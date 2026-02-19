package model

import "time"

// 红包类型
const (
	RedPacketTypeNormal = 1 // 普通红包（等额）
	RedPacketTypeLucky  = 2 // 拼手气红包（随机）
)

// 红包状态
const (
	RedPacketStatusActive  = 1 // 可领取
	RedPacketStatusEmpty   = 2 // 已抢完
	RedPacketStatusExpired = 3 // 已过期
)

type RedPacket struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	SenderID        uint64    `gorm:"not null;index:idx_sender_id"`
	Type            int8      `gorm:"not null"`
	TotalAmount     uint64    `gorm:"not null"`
	TotalCount      uint32    `gorm:"not null"`
	RemainingAmount uint64    `gorm:"not null"`
	RemainingCount  uint32    `gorm:"not null"`
	Status          int8      `gorm:"not null;default:1;index:idx_status_expired"`
	ExpiredAt       time.Time `gorm:"not null;index:idx_status_expired"`
	CreatedAt       time.Time `gorm:"not null"`
}
