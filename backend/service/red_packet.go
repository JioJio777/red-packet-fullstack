package service

import (
	"errors"
	"math/rand"
	"time"

	"red-packet/database"
	"red-packet/model"
	"red-packet/repository"

	"gorm.io/gorm"
)

type SendRedPacketParams struct {
	SenderID    uint64
	Type        int8
	TotalAmount uint64
	TotalCount  uint32
}

type RedPacketDetail struct {
	*model.RedPacket
	SenderName   string
	ClaimedCount int64
	MyClaim      *MyClaim
}

type MyClaim struct {
	Claimed   bool      `json:"claimed"`
	Amount    uint64    `json:"amount,omitempty"`
	ClaimedAt time.Time `json:"claimed_at,omitempty"`
}

type RecordItem struct {
	ReceiverID   uint64    `json:"receiver_id"`
	ReceiverName string    `json:"receiver_name"`
	Amount       uint64    `json:"amount"`
	ClaimedAt    time.Time `json:"claimed_at"`
}

func SendRedPacket(params SendRedPacketParams) (*model.RedPacket, error) {
	if params.TotalAmount < uint64(params.TotalCount) {
		return nil, errors.New("total amount must be >= total count (min 1 fen per person)")
	}

	var redPacket *model.RedPacket

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 扣减发送者余额
		if err := repository.DeductUserBalance(tx, params.SenderID, params.TotalAmount); err != nil {
			return err
		}

		// 验证余额是否足够（受影响行数为 0 说明余额不足）
		senderBalance, err := repository.GetUserBalanceInTx(tx, params.SenderID)
		if err != nil {
			return err
		}
		_ = senderBalance

		// 写流水：支出
		balanceAfter, err := repository.GetUserBalanceInTx(tx, params.SenderID)
		if err != nil {
			return err
		}
		txRecord := &model.Transaction{
			UserID:       params.SenderID,
			Type:         model.TransactionTypeSend,
			Direction:    model.TransactionDirectionOut,
			Amount:       params.TotalAmount,
			BalanceAfter: balanceAfter,
			Remark:       "发红包",
		}

		rp := &model.RedPacket{
			SenderID:        params.SenderID,
			Type:            params.Type,
			TotalAmount:     params.TotalAmount,
			TotalCount:      params.TotalCount,
			RemainingAmount: params.TotalAmount,
			RemainingCount:  params.TotalCount,
			Status:          model.RedPacketStatusActive,
			ExpiredAt:       time.Now().Add(24 * time.Hour),
		}
		if err := repository.CreateRedPacket(rp); err != nil {
			return err
		}

		txRecord.RelatedID = &rp.ID
		if err := repository.CreateTransaction(tx, txRecord); err != nil {
			return err
		}

		redPacket = rp
		return nil
	})

	return redPacket, err
}

func ClaimRedPacket(redPacketID, receiverID uint64) (uint64, error) {
	var claimedAmount uint64

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 加行锁，防止并发超发
		rp, err := repository.GetRedPacketForUpdate(tx, redPacketID)
		if err != nil {
			return errors.New("red packet not found")
		}

		// 状态校验
		if rp.Status != model.RedPacketStatusActive {
			if rp.Status == model.RedPacketStatusEmpty {
				return errors.New("red packet is empty")
			}
			return errors.New("red packet is expired")
		}
		if time.Now().After(rp.ExpiredAt) {
			return errors.New("red packet is expired")
		}

		// 检查是否已领取
		_, err = repository.GetRedPacketRecord(redPacketID, receiverID)
		if err == nil {
			return errors.New("already claimed")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 计算本次领取金额
		amount := calcClaimAmount(rp)
		claimedAmount = amount

		// 更新红包剩余
		rp.RemainingAmount -= amount
		rp.RemainingCount--
		if rp.RemainingCount == 0 {
			rp.Status = model.RedPacketStatusEmpty
		}
		if err := repository.UpdateRedPacket(tx, rp); err != nil {
			return err
		}

		// 写领取记录
		record := &model.RedPacketRecord{
			RedPacketID: redPacketID,
			ReceiverID:  receiverID,
			Amount:      amount,
		}
		if err := repository.CreateRedPacketRecord(tx, record); err != nil {
			return err
		}

		// 增加领取者余额
		if err := repository.AddUserBalance(tx, receiverID, amount); err != nil {
			return err
		}

		// 写流水：收入
		balanceAfter, err := repository.GetUserBalanceInTx(tx, receiverID)
		if err != nil {
			return err
		}
		txRecord := &model.Transaction{
			UserID:       receiverID,
			Type:         model.TransactionTypeReceive,
			Direction:    model.TransactionDirectionIn,
			Amount:       amount,
			BalanceAfter: balanceAfter,
			RelatedID:    &redPacketID,
			Remark:       "领红包",
		}
		return repository.CreateTransaction(tx, txRecord)
	})

	return claimedAmount, err
}

// calcClaimAmount 计算本次领取金额
// 普通红包：等额分配；拼手气红包：随机分配（保证最后一人也有得拿）
func calcClaimAmount(rp *model.RedPacket) uint64 {
	if rp.RemainingCount == 1 {
		return rp.RemainingAmount
	}
	if rp.Type == model.RedPacketTypeNormal {
		return rp.TotalAmount / uint64(rp.TotalCount)
	}
	// 拼手气：随机金额，范围 [1, 剩余均值*2-1]
	maxAmount := rp.RemainingAmount/uint64(rp.RemainingCount)*2 - 1
	if maxAmount < 1 {
		maxAmount = 1
	}
	return uint64(rand.Int63n(int64(maxAmount))) + 1
}

func GetRedPacketDetail(redPacketID, currentUserID uint64) (*RedPacketDetail, error) {
	rp, err := repository.GetRedPacketByID(redPacketID)
	if err != nil {
		return nil, errors.New("red packet not found")
	}

	sender, err := repository.GetUserByID(rp.SenderID)
	if err != nil {
		return nil, err
	}

	claimedCount, _ := repository.CountRedPacketClaimed(redPacketID)

	detail := &RedPacketDetail{
		RedPacket:    rp,
		SenderName:   sender.Username,
		ClaimedCount: claimedCount,
	}

	// 查询当前用户的领取情况
	record, err := repository.GetRedPacketRecord(redPacketID, currentUserID)
	if err == nil {
		detail.MyClaim = &MyClaim{
			Claimed:   true,
			Amount:    record.Amount,
			ClaimedAt: record.CreatedAt,
		}
	} else {
		detail.MyClaim = &MyClaim{Claimed: false}
	}

	return detail, nil
}

func GetRedPacketRecords(redPacketID uint64, page, pageSize int) ([]RecordItem, int64, error) {
	offset := (page - 1) * pageSize
	records, total, err := repository.GetRedPacketRecords(redPacketID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]RecordItem, 0, len(records))
	for _, r := range records {
		user, _ := repository.GetUserByID(r.ReceiverID)
		name := ""
		if user != nil {
			name = user.Username
		}
		items = append(items, RecordItem{
			ReceiverID:   r.ReceiverID,
			ReceiverName: name,
			Amount:       r.Amount,
			ClaimedAt:    r.CreatedAt,
		})
	}
	return items, total, nil
}

func GetSentRedPackets(senderID uint64, page, pageSize int) ([]model.RedPacket, int64, error) {
	offset := (page - 1) * pageSize
	return repository.GetSentRedPackets(senderID, offset, pageSize)
}

func GetReceivedRedPackets(receiverID uint64, page, pageSize int) ([]map[string]interface{}, int64, error) {
	offset := (page - 1) * pageSize
	records, total, err := repository.GetReceivedRedPackets(receiverID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]map[string]interface{}, 0, len(records))
	for _, r := range records {
		rp, _ := repository.GetRedPacketByID(r.RedPacketID)
		sender, _ := repository.GetUserByID(rp.SenderID)
		senderName := ""
		if sender != nil {
			senderName = sender.Username
		}
		result = append(result, map[string]interface{}{
			"red_packet_id": r.RedPacketID,
			"sender_name":   senderName,
			"amount":        r.Amount,
			"claimed_at":    r.CreatedAt,
		})
	}
	return result, total, nil
}
