package handler

import (
	"net/http"
	"strconv"

	"red-packet/pkg/response"
	"red-packet/service"

	"github.com/gin-gonic/gin"
)

type SendRedPacketRequest struct {
	Type        int8   `json:"type" binding:"required,oneof=1 2"`
	TotalAmount uint64 `json:"total_amount" binding:"required,min=1"`
	TotalCount  uint32 `json:"total_count" binding:"required,min=1,max=100"`
}

func SendRedPacket(c *gin.Context) {
	var req SendRedPacketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	senderID, _ := c.Get("user_id")
	rp, err := service.SendRedPacket(service.SendRedPacketParams{
		SenderID:    senderID.(uint64),
		Type:        req.Type,
		TotalAmount: req.TotalAmount,
		TotalCount:  req.TotalCount,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":           rp.ID,
		"type":         rp.Type,
		"total_amount": rp.TotalAmount,
		"total_count":  rp.TotalCount,
		"expired_at":   rp.ExpiredAt,
	})
}

func ClaimRedPacket(c *gin.Context) {
	redPacketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "invalid id")
		return
	}

	receiverID, _ := c.Get("user_id")
	amount, err := service.ClaimRedPacket(redPacketID, receiverID.(uint64))
	if err != nil {
		code := 400
		switch err.Error() {
		case "red packet is empty":
			code = 1002
		case "red packet is expired":
			code = 1003
		case "already claimed":
			code = 1004
		}
		response.Fail(c, http.StatusBadRequest, code, err.Error())
		return
	}

	response.Success(c, gin.H{"amount": amount})
}

func GetRedPacketDetail(c *gin.Context) {
	redPacketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "invalid id")
		return
	}

	currentUserID, _ := c.Get("user_id")
	detail, err := service.GetRedPacketDetail(redPacketID, currentUserID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":               detail.ID,
		"sender_id":        detail.SenderID,
		"sender_name":      detail.SenderName,
		"type":             detail.Type,
		"total_amount":     detail.TotalAmount,
		"total_count":      detail.TotalCount,
		"remaining_amount": detail.RemainingAmount,
		"remaining_count":  detail.RemainingCount,
		"claimed_count":    detail.ClaimedCount,
		"status":           detail.Status,
		"expired_at":       detail.ExpiredAt,
		"created_at":       detail.CreatedAt,
		"my_claim":         detail.MyClaim,
	})
}

func GetRedPacketRecords(c *gin.Context) {
	redPacketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "invalid id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 50 {
		pageSize = 50
	}

	records, total, err := service.GetRedPacketRecords(redPacketID, page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, "internal error")
		return
	}

	response.Success(c, gin.H{"total": total, "list": records})
}

func GetSentRedPackets(c *gin.Context) {
	senderID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := service.GetSentRedPackets(senderID.(uint64), page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, "internal error")
		return
	}

	response.Success(c, gin.H{"total": total, "list": list})
}

func GetReceivedRedPackets(c *gin.Context) {
	receiverID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := service.GetReceivedRedPackets(receiverID.(uint64), page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, "internal error")
		return
	}

	response.Success(c, gin.H{"total": total, "list": list})
}
