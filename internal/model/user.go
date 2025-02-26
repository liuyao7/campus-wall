package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
    ID           uint64         `gorm:"primarykey" json:"id"`
    Username     string         `gorm:"uniqueIndex;size:32" json:"username"`
    Password     string         `gorm:"size:128" json:"-"`
    OpenID       string         `gorm:"uniqueIndex;size:64" json:"open_id"`
	UnionID      string         `gorm:"uniqueIndex;size:64" json:"union_id"`
	SessionKey   string         `gorm:"size:255" json:"-"`
    Nickname     string         `gorm:"size:32" json:"nickname"`
    Avatar       string         `gorm:"size:255" json:"avatar"`
    Gender       int            `gorm:"default:0" json:"gender"`
    Introduction string         `gorm:"size:255" json:"introduction"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserUpdateRequest 用户信息更新请求
type UserUpdateRequest struct {
    Nickname    *string `json:"nickname" binding:"omitempty,min=2,max=32"`
    Gender      *int    `json:"gender" binding:"omitempty,oneof=0 1 2"` // 0-未知 1-男 2-女
    Birthday    *string `json:"birthday" binding:"omitempty,datetime=2006-01-02"`
    Introduction *string `json:"introduction" binding:"omitempty,max=200"`
    Location    *string `json:"location" binding:"omitempty,max=100"`
    Email       *string `json:"email" binding:"omitempty,email"`
    Phone       *string `json:"phone" binding:"omitempty,len=11"`
}