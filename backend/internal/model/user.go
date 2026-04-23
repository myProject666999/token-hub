package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password    string         `json:"-" gorm:"size:255;not null"`
	Email       string         `json:"email" gorm:"size:100"`
	Phone       string         `json:"phone" gorm:"size:20"`
	Nickname    string         `json:"nickname" gorm:"size:50"`
	Avatar      string         `json:"avatar" gorm:"size:255"`
	Role        string         `json:"role" gorm:"size:20;default:'user'"`
	Status      int            `json:"status" gorm:"default:1"`
	Points      float64        `json:"points" gorm:"type:decimal(18,2);default:0"`
	TotalPoints float64        `json:"total_points" gorm:"type:decimal(18,2);default:0"`
	UsedPoints  float64        `json:"used_points" gorm:"type:decimal(18,2);default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type UserUpdateRequest struct {
	Email    string `json:"email" binding:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type UserPointsResponse struct {
	Points      float64 `json:"points"`
	TotalPoints float64 `json:"total_points"`
	UsedPoints  float64 `json:"used_points"`
}
