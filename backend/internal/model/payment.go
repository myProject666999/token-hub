package model

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Description string        `json:"description" gorm:"size:255"`
	Logo        string         `json:"logo" gorm:"size:255"`
	Config      string         `json:"-" gorm:"type:text"`
	Status      int            `json:"status" gorm:"default:1"`
	Sort        int            `json:"sort" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type RechargeRecord struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrderNo        string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID         uint           `json:"user_id" gorm:"not null;index"`
	User           User           `json:"user" gorm:"foreignKey:UserID"`
	PaymentMethodID uint          `json:"payment_method_id" gorm:"not null"`
	PaymentMethod  PaymentMethod `json:"payment_method" gorm:"foreignKey:PaymentMethodID"`
	Amount         float64        `json:"amount" gorm:"type:decimal(18,2);not null"`
	Points         float64        `json:"points" gorm:"type:decimal(18,2);not null"`
	PointsRate     float64        `json:"points_rate" gorm:"type:decimal(10,2);not null"`
	Status         string         `json:"status" gorm:"size:20;default:'pending'"`
	TradeNo        string         `json:"trade_no" gorm:"size:100"`
	PaidAt         *time.Time     `json:"paid_at"`
	ExpiredAt      *time.Time     `json:"expired_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusExpired = "expired"
	PaymentStatusRefund  = "refund"
)

type RechargeRequest struct {
	PaymentMethodID uint    `json:"payment_method_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,min=1"`
}

type PaymentNotifyRequest struct {
	OrderNo string  `json:"order_no" binding:"required"`
	TradeNo string  `json:"trade_no"`
	Amount  float64 `json:"amount" binding:"required"`
	Status  string  `json:"status" binding:"required"`
}

type AlipayConfig struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
	Sandbox    bool   `json:"sandbox"`
}

type WechatConfig struct {
	AppID     string `json:"app_id"`
	MchID     string `json:"mch_id"`
	APIKey    string `json:"api_key"`
	NotifyURL string `json:"notify_url"`
	CertPath  string `json:"cert_path"`
	KeyPath   string `json:"key_path"`
}
