// internal/model/post.go

package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 帖子模型
type Post struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    UserID    uint           `json:"user_id"`
    Content   string         `json:"content"`
    Images    []PostImage    `json:"images"`
    Tags      []Tag         `json:"tags" gorm:"many2many:post_tags;"`
    Likes     int           `json:"likes" gorm:"default:0"`
    Comments  int           `json:"comments" gorm:"default:0"`
    Status    int           `json:"status" gorm:"default:1"` // 1-正常 2-待审核 3-已删除
    User      User          `json:"user" gorm:"foreignKey:UserID"`
    CreatedAt time.Time     `json:"created_at"`
    UpdatedAt time.Time     `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PostImage 帖子图片模型
type PostImage struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    PostID    uint      `json:"post_id"`
    URL       string    `json:"url"`
    Order     int       `json:"order"` // 图片顺序
    CreatedAt time.Time `json:"created_at"`
}

// Tag 标签模型
type Tag struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    Name      string    `json:"name" gorm:"uniqueIndex"`
    Count     int       `json:"count" gorm:"default:0"` // 使用次数
    CreatedAt time.Time `json:"created_at"`
}

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
    Content string   `json:"content" binding:"required,min=1,max=5000"`
    Tags    []string `json:"tags" binding:"omitempty,dive,min=1,max=20"`
    Images  []string `json:"images" binding:"omitempty,dive,url"`
}

// PostListQuery 帖子列表查询参数
type PostListQuery struct {
    Keyword   string    `form:"keyword"`                         // 搜索关键词
    TopicID   uint      `form:"topic_id"`                       // 话题ID
    UserID    uint      `form:"user_id"`                        // 用户ID
    Tag       string    `form:"tag"`                            // 标签
    StartTime time.Time `form:"start_time" time_format:"2006-01-02"` // 开始时间
    EndTime   time.Time `form:"end_time" time_format:"2006-01-02"`   // 结束时间
    Sort      string    `form:"sort" binding:"omitempty,oneof=new hot comment"` // 排序方式: new-最新 hot-最热 comment-最多评论
    Page      int       `form:"page,default=1" binding:"min=1"`
    Size      int       `form:"size,default=10" binding:"min=1,max=50"`
}

// PostListResponse 帖子列表响应
type PostListResponse struct {
    Posts []Post `json:"posts"`
    Total int64  `json:"total"`
    Page  int    `json:"page"`
    Size  int    `json:"size"`
}