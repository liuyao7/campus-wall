package repository

import (
	"campus-wall/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// User相关
func (r *Repository) CreateUser(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *Repository) GetUserByOpenID(openID string) (*model.User, error) {
    var user model.User
    err := r.db.Where("open_id = ?", openID).First(&user).Error
    return &user, err
}

// Post相关
func (r *Repository) CreatePost(post *model.Post) error {
    return r.db.Create(post).Error
}

func (r *Repository) GetPosts(page, pageSize int) ([]model.Post, int64, error) {
    var posts []model.Post
    var total int64

    err := r.db.Model(&model.Post{}).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    err = r.db.Preload("User").
        Where("status = ?", 1).
        Order("created_at DESC").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&posts).Error

    return posts, total, err
}

func (r *Repository) GetPostByID(id uint) (*model.Post, error) {
    var post model.Post
    err := r.db.Preload("User").First(&post, id).Error
    return &post, err
}

// Comment相关
func (r *Repository) CreateComment(comment *model.Comment) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(comment).Error; err != nil {
            return err
        }
        
        // 更新帖子评论数
        return tx.Model(&model.Post{}).
            Where("id = ?", comment.PostID).
            UpdateColumn("comments", gorm.Expr("comments + ?", 1)).Error
    })
}

func (r *Repository) GetCommentsByPostID(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
    var comments []model.Comment
    var total int64

    err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    err = r.db.Preload("User").
        Where("post_id = ? AND status = ?", postID, 1).
        Order("created_at DESC").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&comments).Error

    return comments, total, err
}

// Like相关
func (r *Repository) CreateLike(like *model.Like) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(like).Error; err != nil {
            return err
        }

        // 更新点赞数
        if like.Type == 1 {
            return tx.Model(&model.Post{}).
                Where("id = ?", like.TargetID).
                UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
        } else {
            return tx.Model(&model.Comment{}).
                Where("id = ?", like.TargetID).
                UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
        }
    })
}

func (r *Repository) DeleteLike(userID, targetID uint, likeType int) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("user_id = ? AND target_id = ? AND type = ?", 
            userID, targetID, likeType).Delete(&model.Like{}).Error; err != nil {
            return err
        }

        // 更新点赞数
        if likeType == 1 {
            return tx.Model(&model.Post{}).
                Where("id = ?", targetID).
                UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
        } else {
            return tx.Model(&model.Comment{}).
                Where("id = ?", targetID).
                UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
        }
    })
}

// IncreaseViewCount
func (r *Repository) IncreaseViewCount(postID uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// AddComment
func (r *Repository) AddComment(comment *model.Comment) error {
	return r.db.Create(comment).Error
}
