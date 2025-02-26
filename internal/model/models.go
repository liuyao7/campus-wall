package model

import (
	"time"
)

// Post 帖子模型
// type Post struct {
//     ID            uint           `gorm:"primarykey" json:"id"`
//     UserID        uint           `json:"user_id"`
//     Content       string         `gorm:"type:text" json:"content"`
//     Images        string         `gorm:"type:text" json:"images"`
//     Type          int           `gorm:"type:tinyint;default:0" json:"type"` // 0:普通 1:匿名
//     Location      string         `gorm:"type:varchar(255)" json:"location"`
//     Likes         int           `gorm:"default:0" json:"likes"`
//     Comments      int           `gorm:"default:0" json:"comments"`
//     Views         int           `gorm:"default:0" json:"views"`
//     Status        int           `gorm:"type:tinyint;default:1" json:"status"` // 0:待审核 1:已发布 2:已删除
//     CreatedAt     time.Time      `json:"created_at"`
//     UpdatedAt     time.Time      `json:"updated_at"`
//     DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
//     User          User           `gorm:"foreignKey:UserID" json:"user"`
// }

// Comment 评论模型
// type Comment struct {
//     ID            uint           `gorm:"primarykey" json:"id"`
//     PostID        uint           `json:"post_id"`
//     UserID        uint           `json:"user_id"`
//     ParentID      *uint          `json:"parent_id"` // 父评论ID，用于回复功能
//     Content       string         `gorm:"type:text" json:"content"`
//     Likes         int           `gorm:"default:0" json:"likes"`
//     Status        int           `gorm:"type:tinyint;default:1" json:"status"`
//     CreatedAt     time.Time      `json:"created_at"`
//     UpdatedAt     time.Time      `json:"updated_at"`
//     DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
//     User          User           `gorm:"foreignKey:UserID" json:"user"`
// }

// Like 点赞模型
type Like struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    uint      `json:"user_id"`
    TargetID  uint      `json:"target_id"`   // 点赞目标ID（可能是帖子或评论）
    Type      int       `json:"type"`        // 1:帖子点赞 2:评论点赞
    CreatedAt time.Time `json:"created_at"`
}

// Report 举报模型
type Report struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    uint      `json:"user_id"`
    TargetID  uint      `json:"target_id"`   // 举报目标ID
    Type      int       `json:"type"`        // 1:帖子举报 2:评论举报
    Reason    string    `gorm:"type:text" json:"reason"`
    Status    int       `gorm:"type:tinyint;default:0" json:"status"` // 0:待处理 1:已处理
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}