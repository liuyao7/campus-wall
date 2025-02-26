// // internal/service/stats.go
package service

// import (
// 	"campus-wall/internal/model"
// 	"context"
// 	"fmt"

// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// )

// type StatsService struct {
//     db    *gorm.DB
//     cache *redis.Client
// }

// // UpdateCommentStats 更新评论统计
// func (s *StatsService) UpdateCommentStats(postID uint) error {
//     // 统计评论数
//     var stats struct {
//         TotalComments  int64
//         TotalReplies   int64
//         ActiveUsers    int64
//         AvgLikes       float64
//     }

//     err := s.db.Transaction(func(tx *gorm.DB) error {
//         // 统计一级评论数
//         if err := tx.Model(&model.Comment{}).
//             Where("post_id = ? AND parent_id IS NULL", postID).
//             Count(&stats.TotalComments).Error; err != nil {
//             return err
//         }

//         // 统计回复数
//         if err := tx.Model(&model.Comment{}).
//             Where("post_id = ? AND parent_id IS NOT NULL", postID).
//             Count(&stats.TotalReplies).Error; err != nil {
//             return err
//         }

//         // 统计参与评论的用户数
//         if err := tx.Model(&model.Comment{}).
//             Where("post_id = ?", postID).
//             Distinct("user_id").
//             Count(&stats.ActiveUsers).Error; err != nil {
//             return err
//         }

//         // 计算平均点赞数
//         if err := tx.Model(&model.Comment{}).
//             Where("post_id = ?", postID).
//             Select("COALESCE(AVG(likes), 0)").
//             Scan(&stats.AvgLikes).Error; err != nil {
//             return err
//         }

//         return nil
//     })

//     if err != nil {
//         return err
//     }

//     // 更新缓存
//     key := fmt.Sprintf("post:%d:comment_stats", postID)
//     return s.cache.HMSet(context.Background(), key, map[string]interface{}{
//         "total_comments": stats.TotalComments,
//         "total_replies":  stats.TotalReplies,
//         "active_users":   stats.ActiveUsers,
//         "avg_likes":      stats.AvgLikes,
//     }).Err()
// }