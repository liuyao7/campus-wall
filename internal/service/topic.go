// // internal/service/topic.go

package service

// import (
// 	"campus-wall/internal/model"
// 	"campus-wall/pkg/errors"

// 	"gorm.io/gorm"
// )

// type TopicService struct {
//     db *gorm.DB
// }

// func NewTopicService(db *gorm.DB) *TopicService {
//     return &TopicService{db: db}
// }

// // CreateTopic 创建话题
// func (s *TopicService) CreateTopic(userID uint, req *model.CreateTopicRequest) (*model.Topic, error) {
//     // 检查话题名是否已存在
//     var exists bool
//     if err := s.db.Model(&model.Topic{}).
//         Where("name = ?", req.Name).
//         Select("1").
//         Scan(&exists).Error; err != nil {
//         return nil, err
//     }
//     if exists {
//         return nil, errors.NewBadRequestError("topic already exists")
//     }

//     topic := &model.Topic{
//         Name:        req.Name,
//         Description: req.Description,
//         Cover:       req.Cover,
//         CreatedBy:   userID,
//     }

//     if err := s.db.Create(topic).Error; err != nil {
//         return nil, errors.NewInternalServerError("failed to create topic", err)
//     }

//     return topic, nil
// }

// // GetTopicList 获取话题列表
// func (s *TopicService) GetTopicList(query *model.TopicListQuery) ([]model.Topic, int64, error) {
//     db := s.db.Model(&model.Topic{})

//     // 条件查询
//     if query.Keyword != "" {
//         db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
//     }
//     if query.IsHot != nil {
//         db = db.Where("is_hot = ?", *query.IsHot)
//     }

//     // 获取总数
//     var total int64
//     if err := db.Count(&total).Error; err != nil {
//         return nil, 0, err
//     }

//     // 排序
//     switch query.Sort {
//     case "hot":
//         db = db.Order("post_count DESC")
//     default:
//         db = db.Order("created_at DESC")
//     }

//     // 分页
//     offset := (query.Page - 1) * query.Size
//     var topics []model.Topic
//     if err := db.Offset(offset).Limit(query.Size).Find(&topics).Error; err != nil {
//         return nil, 0, err
//     }

//     return topics, total, nil
// }

// // FollowTopic 关注话题
// func (s *TopicService) FollowTopic(userID, topicID uint) error {
//     // 检查话题是否存在
//     var topic model.Topic
//     if err := s.db.First(&topic, topicID).Error; err != nil {
//         return errors.NewNotFoundError("topic not found")
//     }

//     // 检查是否已关注
//     var exists bool
//     if err := s.db.Model(&model.TopicFollow{}).
//         Where("user_id = ? AND topic_id = ?", userID, topicID).
//         Select("1").
//         Scan(&exists).Error; err != nil {
//         return err
//     }
//     if exists {
//         return errors.NewBadRequestError("already followed")
//     }

//     // 开启事务
//     return s.db.Transaction(func(tx *gorm.DB) error {
//         // 创建关注记录
//         follow := &model.TopicFollow{
//             UserID:  userID,
//             TopicID: topicID,
//         }
//         if err := tx.Create(follow).Error; err != nil {
//             return err
//         }

//         // 更新话题关注数
//         if err := tx.Model(&topic).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
//             return err
//         }

//         return nil
//     })
// }

// // UnfollowTopic 取消关注话题
// func (s *TopicService) UnfollowTopic(userID, topicID uint) error {
//     return s.db.Transaction(func(tx *gorm.DB) error {
//         // 删除关注记录
//         result := tx.Where("user_id = ? AND topic_id = ?", userID, topicID).Delete(&model.TopicFollow{})
//         if result.Error != nil {
//             return result.Error
//         }
//         if result.RowsAffected == 0 {
//             return errors.NewBadRequestError("not followed")
//         }

//         // 更新话题关注数
//         if err := tx.Model(&model.Topic{}).Where("id = ?", topicID).
//             UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
//             return err
//         }

//         return nil
//     })
// }

// // GetUserFollowedTopics 获取用户关注的话题
// func (s *TopicService) GetUserFollowedTopics(userID uint) ([]model.Topic, error) {
//     var topics []model.Topic
//     err := s.db.Model(&model.Topic{}).
//         Joins("JOIN topic_follows ON topics.id = topic_follows.topic_id").
//         Where("topic_follows.user_id = ?", userID).
//         Find(&topics).Error
//     return topics, err
// }