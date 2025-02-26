// internal/model/like.go

package model

import (
	"time"
)

// PostLike 帖子点赞模型
type PostLike struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    PostID    uint      `json:"post_id"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}

// 定义唯一索引
func (PostLike) TableName() string {
    return "post_likes"
}

func (PostLike) Indexes() []string {
    return []string{
        "CREATE UNIQUE INDEX idx_post_likes_user_post ON post_likes(user_id, post_id)",
    }
}