package model

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:100;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Description string        `json:"description" gorm:"size:500"`
	Logo        string         `json:"logo" gorm:"size:255"`
	Website     string         `json:"website" gorm:"size:255"`
	APIEndpoint string        `json:"api_endpoint" gorm:"size:255"`
	APIKey      string         `json:"-" gorm:"size:500"`
	APISecret   string         `json:"-" gorm:"size:500"`
	Status      int            `json:"status" gorm:"default:1"`
	Sort        int            `json:"sort" gorm:"default:0"`
	Models      []Model        `json:"models" gorm:"foreignKey:ProviderID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Model struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"size:100;not null"`
	Code         string         `json:"code" gorm:"size:50;not null"`
	ProviderID   uint           `json:"provider_id" gorm:"not null"`
	Provider     Provider       `json:"provider" gorm:"foreignKey:ProviderID"`
	Description  string        `json:"description" gorm:"size:500"`
	ModelType    string         `json:"model_type" gorm:"size:50"`
	PointsPer1KInput  float64   `json:"points_per_1k_input" gorm:"type:decimal(10,4);default:0"`
	PointsPer1KOutput float64   `json:"points_per_1k_output" gorm:"type:decimal(10,4);default:0"`
	MaxTokens    int            `json:"max_tokens" gorm:"default:4096"`
	ContextLimit int            `json:"context_limit" gorm:"default:4096"`
	SupportsVision bool         `json:"supports_vision" gorm:"default:false"`
	SupportsFunction bool       `json:"supports_function" gorm:"default:false"`
	Status       int            `json:"status" gorm:"default:1"`
	Sort         int            `json:"sort" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type ProviderCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Website     string `json:"website"`
	APIEndpoint string `json:"api_endpoint" binding:"required"`
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
}

type ModelCreateRequest struct {
	Name              string  `json:"name" binding:"required"`
	Code              string  `json:"code" binding:"required"`
	ProviderID        uint    `json:"provider_id" binding:"required"`
	Description       string  `json:"description"`
	ModelType         string  `json:"model_type"`
	PointsPer1KInput  float64 `json:"points_per_1k_input"`
	PointsPer1KOutput float64 `json:"points_per_1k_output"`
	MaxTokens         int     `json:"max_tokens"`
	ContextLimit      int     `json:"context_limit"`
	SupportsVision    bool    `json:"supports_vision"`
	SupportsFunction  bool    `json:"supports_function"`
	Status            int     `json:"status"`
	Sort              int     `json:"sort"`
}
