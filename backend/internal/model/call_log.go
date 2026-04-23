package model

import (
	"time"

	"gorm.io/gorm"
)

type CallLog struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"not null;index"`
	User            User           `json:"user" gorm:"foreignKey:UserID"`
	ProviderID      uint           `json:"provider_id" gorm:"not null"`
	Provider        Provider       `json:"provider" gorm:"foreignKey:ProviderID"`
	ModelID         uint           `json:"model_id" gorm:"not null"`
	Model           Model          `json:"model" gorm:"foreignKey:ModelID"`
	RequestID       string         `json:"request_id" gorm:"uniqueIndex;size:100;not null"`
	InputTokens     int            `json:"input_tokens" gorm:"default:0"`
	OutputTokens    int            `json:"output_tokens" gorm:"default:0"`
	TotalTokens     int            `json:"total_tokens" gorm:"default:0"`
	PointsConsumed  float64        `json:"points_consumed" gorm:"type:decimal(18,4);default:0"`
	Status          string         `json:"status" gorm:"size:20;default:'success'"`
	ErrorMessage    string         `json:"error_message" gorm:"type:text"`
	Duration        int64          `json:"duration" gorm:"default:0"`
	IP              string         `json:"ip" gorm:"size:50"`
	UserAgent       string         `json:"user_agent" gorm:"size:500"`
	RequestContent  string         `json:"-" gorm:"type:text"`
	ResponseContent string         `json:"-" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

const (
	CallStatusSuccess = "success"
	CallStatusFailed  = "failed"
	CallStatusPending = "pending"
)

type CallStatistics struct {
	UserID          uint    `json:"user_id"`
	ProviderID      uint    `json:"provider_id"`
	ModelID         uint    `json:"model_id"`
	TotalCalls      int64   `json:"total_calls"`
	SuccessCalls    int64   `json:"success_calls"`
	FailedCalls     int64   `json:"failed_calls"`
	TotalInputTokens int64  `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	TotalTokens     int64   `json:"total_tokens"`
	TotalPoints     float64 `json:"total_points"`
}

type DailyStatistics struct {
	Date            string  `json:"date"`
	UserID          uint    `json:"user_id,omitempty"`
	ProviderID      uint    `json:"provider_id,omitempty"`
	ModelID         uint    `json:"model_id,omitempty"`
	TotalCalls      int64   `json:"total_calls"`
	SuccessCalls    int64   `json:"success_calls"`
	FailedCalls     int64   `json:"failed_calls"`
	TotalInputTokens int64  `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	TotalTokens     int64   `json:"total_tokens"`
	TotalPoints     float64 `json:"total_points"`
}
