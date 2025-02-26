// // internal/service/notification.go
package service

// import (
// 	"campus-wall/internal/model"

// 	"gorm.io/gorm"
// )

// type NotificationService struct {
//     db *gorm.DB
// }

// func (s *NotificationService) CreateCommentNotification(comment *model.Comment) error {
//     // 创建评论通知
//     notification := &model.Notification{
//         Type:      "comment",
//         UserID:    comment.UserID,
//         TargetID:  comment.PostID,
//         Content:   comment.Content,
//         Status:    1,
//     }

//     // 如果是回复评论，通知被回复的用户
//     if comment.ParentID != nil {
//         var parentComment model.Comment
//         if err := s.db.First(&parentComment, *comment.ParentID).Error; err != nil {
//             return err
//         }
//         notification.ReceiverID = parentComment.UserID
//     } else {
//         // 如果是评论帖子，通知帖子作者
//         var post model.Post
//         if err := s.db.First(&post, comment.PostID).Error; err != nil {
//             return err
//         }
//         notification.ReceiverID = post.UserID
//     }

//     return s.db.Create(notification).Error
// }