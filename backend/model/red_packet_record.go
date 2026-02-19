package model

import "time"

type RedPacketRecord struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	RedPacketID uint64    `gorm:"not null;uniqueIndex:uk_packet_receiver;index:idx_red_packet_id"`
	ReceiverID  uint64    `gorm:"not null;uniqueIndex:uk_packet_receiver;index:idx_receiver_id"`
	Amount      uint64    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
}
