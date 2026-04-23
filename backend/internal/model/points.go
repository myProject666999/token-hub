package model

import (
	"time"

	"gorm.io/gorm"
)

type PointsRecord struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Type        string         `json:"type" gorm:"size:20;not null"`
	Amount      float64        `json:"amount" gorm:"type:decimal(18,2);not null"`
	Balance     float64        `json:"balance" gorm:"type:decimal(18,2);not null"`
	Description string        `json:"description" gorm:"size:255"`
	RelatedID   string         `json:"related_id" gorm:"size:100"`
	RelatedType string        `json:"related_type" gorm:"size:50"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

const (
	PointsTypeRecharge = "recharge"
	PointsTypeConsume  = "consume"
	PointsTypeRefund   = "refund"
	PointsTypeGift     = "gift"
)

type PointsConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Key         string         `json:"key" gorm:"uniqueIndex;size:50;not null"`
	Value       string         `json:"value" gorm:"size:255"`
	Description string        `json:"description" gorm:"size:255"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

const (
	PointsKeyRate = "points_rate"
)

type PointsStatistics struct {
	UserID          uint    `json:"user_id"`
	TotalPoints     float64 `json:"total_points"`
	UsedPoints      float64 `json:"used_points"`
	Balance         float64 `json:"balance"`
	RechargeCount   int64   `json:"recharge_count"`
	CallCount       int64   `json:"call_count"`
}
