// // internal/service/post.go

package service

// import (
// 	"campus-wall/internal/model"
// 	"campus-wall/pkg/errors"
// 	"campus-wall/pkg/storage"
// 	"mime/multipart"

// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// )

// type PostService struct {
//     db      *gorm.DB
//     storage storage.Storage
// 	cache   *redis.Client
// }

// func NewPostService(db *gorm.DB, storage storage.Storage) *PostService {
//     return &PostService{
//         db:      db,
//         storage: storage,
//     }
// }

// // UploadPostImages 上传帖子图片
// func (s *PostService) UploadPostImages(files []*multipart.FileHeader) ([]string, error) {
//     var urls []string

//     for _, file := range files {
//         // 验证图片
//         if err := validateImageFile(file); err != nil {
//             return nil, err
//         }

//         // 上传图片
//         url, err := s.storage.Upload(file, "posts")
//         if err != nil {
//             return nil, errors.NewInternalServerError("failed to upload image", err)
//         }

//         urls = append(urls, url)
//     }

//     return urls, nil
// }

// // CreatePost 创建帖子
// func (s *PostService) CreatePost(userID uint, req *model.CreatePostRequest) (*model.Post, error) {
//     // 开启事务
//     tx := s.db.Begin()
//     defer func() {
//         if r := recover(); r != nil {
//             tx.Rollback()
//         }
//     }()

//     // 1. 创建帖子
//     post := &model.Post{
//         UserID:  userID,
//         Content: req.Content,
//         Status:  1, // 正常状态
//     }

//     if err := tx.Create(post).Error; err != nil {
//         tx.Rollback()
//         return nil, errors.NewInternalServerError("failed to create post", err)
//     }

//     // 2. 处理图片
//     if len(req.Images) > 0 {
//         var postImages []model.PostImage
//         for i, url := range req.Images {
//             postImages = append(postImages, model.PostImage{
//                 PostID: post.ID,
//                 URL:    url,
//                 Order:  i + 1,
//             })
//         }
//         if err := tx.Create(&postImages).Error; err != nil {
//             tx.Rollback()
//             return nil, errors.NewInternalServerError("failed to create post images", err)
//         }
//     }

//     // 3. 处理标签
//     if len(req.Tags) > 0 {
//         for _, tagName := range req.Tags {
//             var tag model.Tag
//             // 查找或创建标签
//             err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, model.Tag{Name: tagName}).Error
//             if err != nil {
//                 tx.Rollback()
//                 return nil, errors.NewInternalServerError("failed to process tags", err)
//             }

//             // 更新标签使用次数
//             if err := tx.Model(&tag).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error; err != nil {
//                 tx.Rollback()
//                 return nil, errors.NewInternalServerError("failed to update tag count", err)
//             }

//             // 关联帖子和标签
//             if err := tx.Model(post).Association("Tags").Append(&tag); err != nil {
//                 tx.Rollback()
//                 return nil, errors.NewInternalServerError("failed to associate tags", err)
//             }
//         }
//     }

//     // 提交事务
//     if err := tx.Commit().Error; err != nil {
//         return nil, errors.NewInternalServerError("failed to commit transaction", err)
//     }

//     // 4. 返回完整的帖子信息
//     var completePost model.Post
//     if err := s.db.Preload("User").Preload("Images").Preload("Tags").First(&completePost, post.ID).Error; err != nil {
//         return nil, errors.NewInternalServerError("failed to get complete post", err)
//     }

//     return &completePost, nil
// }

// // GetPostList 获取帖子列表
// func (s *PostService) GetPostList(query *model.PostListQuery) (*model.PostListResponse, error) {
//     db := s.db.Model(&model.Post{}).
//         Preload("User").
//         Preload("Images").
//         Preload("Tags").
//         Where("status = ?", 1) // 只查询正常状态的帖子

//     // 1. 关键词搜索
//     if query.Keyword != "" {
//         db = db.Where("content LIKE ?", "%"+query.Keyword+"%")
//     }

//     // 2. 话题筛选
//     if query.TopicID > 0 {
//         db = db.Joins("JOIN post_topics ON posts.id = post_topics.post_id").
//             Where("post_topics.topic_id = ?", query.TopicID)
//     }

//     // 3. 用户筛选
//     if query.UserID > 0 {
//         db = db.Where("user_id = ?", query.UserID)
//     }

//     // 4. 标签筛选
//     if query.Tag != "" {
//         db = db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
//             Joins("JOIN tags ON post_tags.tag_id = tags.id").
//             Where("tags.name = ?", query.Tag)
//     }

//     // 5. 时间范围筛选
//     if !query.StartTime.IsZero() {
//         db = db.Where("created_at >= ?", query.StartTime)
//     }
//     if !query.EndTime.IsZero() {
//         db = db.Where("created_at <= ?", query.EndTime)
//     }

//     // 6. 排序
//     switch query.Sort {
//     case "hot":
//         db = db.Order("likes DESC")
//     case "comment":
//         db = db.Order("comments DESC")
//     default:
//         db = db.Order("created_at DESC")
//     }

//     // 7. 获取总数
//     var total int64
//     if err := db.Count(&total).Error; err != nil {
//         return nil, errors.NewInternalServerError("failed to count posts", err)
//     }

//     // 8. 分页查询
//     offset := (query.Page - 1) * query.Size
//     var posts []model.Post
//     if err := db.Offset(offset).Limit(query.Size).Find(&posts).Error; err != nil {
//         return nil, errors.NewInternalServerError("failed to get posts", err)
//     }

//     // 9. 处理帖子数据（例如：添加是否点赞标记）
//     if err := s.processPostList(&posts); err != nil {
//         return nil, err
//     }

//     return &model.PostListResponse{
//         Posts: posts,
//         Total: total,
//         Page:  query.Page,
//         Size:  query.Size,
//     }, nil
// }

// // processPostList 处理帖子列表数据
// func (s *PostService) processPostList(posts *[]model.Post) error {
//     // 这里可以添加一些额外的处理逻辑
//     // 例如：添加是否点赞标记、处理图片URL等
//     return nil
// }

// // GetUserPosts 获取用户的帖子列表
// func (s *PostService) GetUserPosts(userID uint, query *model.PostListQuery) (*model.PostListResponse, error) {
//     query.UserID = userID
//     return s.GetPostList(query)
// }

// // GetTopicPosts 获取话题下的帖子列表
// func (s *PostService) GetTopicPosts(topicID uint, query *model.PostListQuery) (*model.PostListResponse, error) {
//     query.TopicID = topicID
//     return s.GetPostList(query)
// }

// // GetPostList with cache
// // func (s *PostService) GetPostList(query *model.PostListQuery) (*model.PostListResponse, error) {
// //     // 1. 尝试从缓存获取
// //     cacheKey := s.generateCacheKey(query)
// //     if response, err := s.getFromCache(cacheKey); err == nil {
// //         return response, nil
// //     }

// //     // 2. 从数据库获取
// //     response, err := s.getPostListFromDB(query)
// //     if err != nil {
// //         return nil, err
// //     }

// //     // 3. 存入缓存
// //     go s.setCache(cacheKey, response)

// //     return response, nil
// // }

// // // 生成缓存key
// // func (s *PostService) generateCacheKey(query *model.PostListQuery) string {
// //     return fmt.Sprintf("posts:%s:%d:%d:%s:%d:%d",
// //         query.Keyword,
// //         query.TopicID,
// //         query.UserID,
// //         query.Sort,
// //         query.Page,
// //         query.Size,
// //     )
// // }

// // // 从缓存获取
// // func (s *PostService) getFromCache(key string) (*model.PostListResponse, error) {
// //     data, err := s.cache.Get(context.Background(), key).Bytes()
// //     if err != nil {
// //         return nil, err
// //     }

// //     var response model.PostListResponse
// //     if err := json.Unmarshal(data, &response); err != nil {
// //         return nil, err
// //     }

// //     return &response, nil
// // }

// // // 设置缓存
// // func (s *PostService) setCache(key string, response *model.PostListResponse) {
// //     data, err := json.Marshal(response)
// //     if err != nil {
// //         return
// //     }

// //     s.cache.Set(context.Background(), key, data, time.Minute*5)
// // }