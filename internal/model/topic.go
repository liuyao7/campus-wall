// internal/model/topic.go

package model

import (
	"time"

	"gorm.io/gorm"
)

// Topic 话题模型
type Topic struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Name        string         `json:"name" gorm:"uniqueIndex;size:50"`
    Description string         `json:"description" gorm:"size:200"`
    Cover       string         `json:"cover"`                    // 话题封面图
    PostCount   int           `json:"post_count" gorm:"default:0"` // 帖子数量
    FollowCount int           `json:"follow_count" gorm:"default:0"` // 关注数量
    Status      int           `json:"status" gorm:"default:1"`   // 1-正常 2-禁用
    IsHot       bool          `json:"is_hot" gorm:"default:false"` // 是否热门
    CreatedBy   uint          `json:"created_by"`               // 创建者ID
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// PostTopic 帖子话题关联表
type PostTopic struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    PostID    uint      `json:"post_id"`
    TopicID   uint      `json:"topic_id"`
    CreatedAt time.Time `json:"created_at"`
}

// TopicFollow 话题关注表
type TopicFollow struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    UserID    uint      `json:"user_id"`
    TopicID   uint      `json:"topic_id"`
    CreatedAt time.Time `json:"created_at"`
}

// 请求和响应结构体
type CreateTopicRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=50"`
    Description string `json:"description" binding:"required,max=200"`
    Cover       string `json:"cover" binding:"omitempty,url"`
}

type TopicListQuery struct {
    Keyword string `form:"keyword"`
    IsHot   *bool  `form:"is_hot"`
    Sort    string `form:"sort" binding:"omitempty,oneof=new hot"` // new-最新 hot-最热
    Page    int    `form:"page,default=1" binding:"min=1"`
    Size    int    `form:"size,default=10" binding:"min=1,max=50"`
}

// 更新Post模型，添加Topics关联
func (Post) TableName() string {
    return "posts"
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
    // 更新话题的帖子数量
    var topicIDs []uint
    if err := tx.Model(&PostTopic{}).Where("post_id = ?", p.ID).Pluck("topic_id", &topicIDs).Error; err != nil {
        return err
    }
    
    if len(topicIDs) > 0 {
        return tx.Model(&Topic{}).Where("id IN ?", topicIDs).
            UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
    }
    return nil
}