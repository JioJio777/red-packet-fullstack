package model

import "time"

// 流水类型
const (
	TransactionTypeRecharge = "recharge" // 充值
	TransactionTypeSend     = "send"     // 发红包（扣款）
	TransactionTypeReceive  = "receive"  // 领红包（到账）
	TransactionTypeRefund   = "refund"   // 红包过期退款
)

// 资金方向
const (
	TransactionDirectionIn  = 1 // 收入
	TransactionDirectionOut = 2 // 支出
)

type Transaction struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	UserID       uint64    `gorm:"not null;index:idx_user_created,priority:1"`
	Type         string    `gorm:"type:varchar(20);not null"`
	Direction    int8      `gorm:"not null"`
	Amount       uint64    `gorm:"not null"`
	BalanceAfter uint64    `gorm:"not null"`
	RelatedID    *uint64   `gorm:"index:idx_related_id"`
	Remark       string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time `gorm:"not null;index:idx_user_created,priority:2"`
}
