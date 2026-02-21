package repository

import (
	"errors"
	"red-packet/database"
	"red-packet/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateRedPacket(rp *model.RedPacket) error {
	return database.DB.Create(rp).Error
}

func GetRedPacketByID(id uint64) (*model.RedPacket, error) {
	var rp model.RedPacket
	err := database.DB.First(&rp, id).Error
	if err != nil {
		return nil, err
	}
	return &rp, nil
}

// GetRedPacketForUpdate 加行锁查询，用于领红包的并发控制
func GetRedPacketForUpdate(tx *gorm.DB, id uint64) (*model.RedPacket, error) {
	var rp model.RedPacket
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&rp, id).Error
	if err != nil {
		return nil, err
	}
	return &rp, nil
}

func UpdateRedPacket(tx *gorm.DB, rp *model.RedPacket) error {
	return tx.Save(rp).Error
}

func CreateRedPacketRecord(tx *gorm.DB, record *model.RedPacketRecord) error {
	return tx.Create(record).Error
}

func GetRedPacketRecord(redPacketID, receiverID uint64) (*model.RedPacketRecord, error) {
	var record model.RedPacketRecord
	err := database.DB.Where("red_packet_id = ? AND receiver_id = ?", redPacketID, receiverID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func GetRedPacketRecords(redPacketID uint64, offset, limit int) ([]model.RedPacketRecord, int64, error) {
	var records []model.RedPacketRecord
	var total int64
	database.DB.Model(&model.RedPacketRecord{}).Where("red_packet_id = ?", redPacketID).Count(&total)
	err := database.DB.Where("red_packet_id = ?", redPacketID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&records).Error
	return records, total, err
}

func CountRedPacketClaimed(redPacketID uint64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.RedPacketRecord{}).Where("red_packet_id = ?", redPacketID).Count(&count).Error
	return count, err
}

func GetSentRedPackets(senderID uint64, offset, limit int) ([]model.RedPacket, int64, error) {
	var list []model.RedPacket
	var total int64
	database.DB.Model(&model.RedPacket{}).Where("sender_id = ?", senderID).Count(&total)
	err := database.DB.Where("sender_id = ?", senderID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, total, err
}

func GetReceivedRedPackets(receiverID uint64, offset, limit int) ([]model.RedPacketRecord, int64, error) {
	var list []model.RedPacketRecord
	var total int64
	database.DB.Model(&model.RedPacketRecord{}).Where("receiver_id = ?", receiverID).Count(&total)
	err := database.DB.Where("receiver_id = ?", receiverID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, total, err
}

func CreateTransaction(tx *gorm.DB, t *model.Transaction) error {
	return tx.Create(t).Error
}

func DeductUserBalance(tx *gorm.DB, userID uint64, amount uint64) error {
	result := tx.Model(&model.User{}).Where("id = ? AND balance >= ?", userID, amount).
		UpdateColumn("balance", gorm.Expr("balance - ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient balance")
	}
	return nil
}

func AddUserBalance(tx *gorm.DB, userID uint64, amount uint64) error {
	return tx.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

func GetUserBalanceInTx(tx *gorm.DB, userID uint64) (uint64, error) {
	var user model.User
	err := tx.Select("balance").First(&user, userID).Error
	return user.Balance, err
}
