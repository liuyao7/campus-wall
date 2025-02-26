// internal/model/comment.go

package model

import (
	"time"
)

// Comment 评论模型
type Comment struct {
    ID         uint      `json:"id" gorm:"primarykey"`
    PostID     uint      `json:"post_id"`
    UserID     uint      `json:"user_id"`
    ParentID   *uint     `json:"parent_id"`                    // 父评论ID，为空表示一级评论
    Content    string    `json:"content" gorm:"type:text"`     // 评论内容
    Likes      int       `json:"likes" gorm:"default:0"`       // 点赞数
    Status     int       `json:"status" gorm:"default:1"`      // 状态：1-正常 2-删除
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    // 关联
    User       User      `json:"user" gorm:"foreignKey:UserID"`
    Parent     *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Replies    []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// CommentRequest 评论请求
type CommentRequest struct {
    Content   string `json:"content" binding:"required,min=1,max=500"`
    ParentID  *uint  `json:"parent_id"`
}

// CommentResponse 评论响应
type CommentResponse struct {
    Comment
    IsLiked bool `json:"is_liked"` // 当前用户是否点赞
}

// CommentListQuery 评论列表查询参数
type CommentListQuery struct {
    PostID uint   `form:"post_id" binding:"required"`
    Sort   string `form:"sort" binding:"omitempty,oneof=new hot"` // new-最新 hot-最热
    Page   int    `form:"page,default=1" binding:"min=1"`
    Size   int    `form:"size,default=20" binding:"min=1,max=100"`
}

// 评论列表响应
type CommentListResponse struct {
    Total   int             `json:"total"`
    Comments []CommentResponse `json:"comments"`
}