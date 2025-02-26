package main

import (
	"campus-wall/pkg/auth"
	"campus-wall/pkg/config"
	"campus-wall/pkg/database"
	"campus-wall/pkg/logger"
	"campus-wall/pkg/wechat"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "campus-wall/docs" // 导入 swagger docs

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// @title Campus Wall API
// @version 1.0
// @description Campus Wall Backend API Documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

func main() {
    // 加载 .env 文件
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // 初始化日志
    logger.InitLogger(os.Getenv("GIN_MODE"))
    defer logger.Log.Sync()

    // 加载配置
    cfg, err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
    if err != nil {
        logger.Fatal("Failed to load config", zap.Error(err))
    }

    // 初始化数据库
    db, err := database.InitDB(&cfg.Database)
    if err != nil {
        logger.Fatal("Failed to initialize database", zap.Error(err))
    }

    // 初始化Redis客户端
    // redisClient, err := redis.NewRedisClient(&cfg.Redis)
    // if err != nil {
    //     logger.Fatal("Failed to initialize Redis", zap.Error(err))
    // }
    // defer redisClient.Close()

    // 初始化存储服务
    // var fileStorage storage.Storage
    // if cfg.Storage.Type == "local" {
    //     fileStorage = storage.NewLocalStorage(
    //         cfg.Storage.Local.Path,
    //         cfg.Storage.Local.BaseURL,
    //     )
    // } else {
    //     fileStorage, err = storage.NewOSSStorage(
    //         cfg.Storage.OSS.Endpoint,
    //         cfg.Storage.OSS.AccessKeyID,
    //         cfg.Storage.OSS.AccessKeySecret,
    //         cfg.Storage.OSS.BucketName,
    //     )
    //     if err != nil {
    //         logger.Fatal("Failed to initialize OSS storage", zap.Error(err))
    //     }
    // }

    // 打印一下配置信息
    logger.Info("Config", zap.Any("config", cfg))

    // 初始化微信小程序客户端
    miniProgram := wechat.NewMiniProgram(
        cfg.WeChat.MiniProgram.AppID,
        cfg.WeChat.MiniProgram.AppSecret,
    )

    // 初始化JWT maker
    tokenMaker := auth.NewJWTMaker(cfg.JWT.SecretKey)

    // 设置gin模式
    gin.SetMode(os.Getenv("GIN_MODE"))

    // 配置路由
    routerConfig := &router.Config{
        TokenMaker:     tokenMaker,
        WechatHandler:  wechatHandler,
        UserHandler:    userHandler,
        PostHandler:    postHandler,
        LikeHandler:    likeHandler,
        CommentHandler: commentHandler,
        RateLimit:      rate.Limit(cfg.Server.RateLimit.Requests),
        RateBurst:      cfg.Server.RateLimit.Burst,
    }

    // 创建路由
    router := router.NewRouter(routerConfig)
    r := router.Engine()

    /* // 设置路由
    // r := router.SetupRouter(routerConfig)

    // 创建路由
    // r := gin.New()

    // 中间件
    // r.Use(
    //     gin.Recovery(),
    //     middleware.Logger(),
    //     middleware.CORS(),
    //     middleware.RateLimit(rate.Limit(cfg.Server.RateLimit.Requests), cfg.Server.RateLimit.Burst),
    // )

    // 初始化Redis工具
    // redisUtil := redis.NewRedisUtil(redisClient)

    // 初始化services
    // userService := service.NewUserService(db, tokenMaker, fileStorage)
    // wechatService := service.NewWeChatService(db, tokenMaker, miniProgram)
    // postService := service.NewPostService(db, fileStorage)
    // likeService := service.NewLikeService(db, redisClient)
    // commentService := service.NewCommentService(db, redisClient)
    // likeService := service.NewLikeService(db, redisUtil)
    // commentService := service.NewCommentService(db, redisUtil)
    
    // 初始化handlers
    // userHandler := handler.NewUserHandler(userService)
    // wechatHandler := handler.NewWeChatHandler(wechatService)
    // postHandler := handler.NewPostHandler(postService)
    // likeHandler := handler.NewLikeHandler(likeService)
    // commentHandler := handler.NewCommentHandler(commentService)

    // API路由组
    // api := r.Group("/api")
    {
        // 公开接口
        // public := api.Group("/")
        // {
            // public.POST("/register", userHandler.Register)
            // public.POST("/login", userHandler.Login)
        //     public.GET("/ping", func(c *gin.Context) {
        //         c.JSON(http.StatusOK, gin.H{"message": "pong"})
        //     })
        // }

        // 微信相关接口
        // wx := api.Group("/wx")
        // {
        //     wx.POST("/login", wechatHandler.MiniProgramLogin)
        //     wx.POST("/decrypt/userinfo", wechatHandler.DecryptUserInfo)
        //     wx.POST("/decrypt/phone", wechatHandler.DecryptPhoneNumber)
        // }

        // 需要认证的接口
        // auth := api.Group("/").Use(middleware.AuthMiddleware(tokenMaker))
        // {
            // 用户相关
            // auth.PUT("/users/me", userHandler.UpdateUser)
            // auth.POST("/users/avatar", userHandler.UploadAvatar)

            // 帖子相关
            // auth.POST("/posts/images", postHandler.UploadImages)
            // auth.POST("/posts", postHandler.CreatePost)
            // auth.GET("/posts", postHandler.GetPostList)
            // auth.GET("/posts/users/:user_id", postHandler.GetUserPosts)
            // auth.GET("/posts/topics/:topic_id", postHandler.GetTopicPosts)

            // 点赞相关
            // auth.POST("/posts/:id/like", likeHandler.LikePost)
            // auth.DELETE("/posts/:id/like", likeHandler.UnlikePost)
            // auth.GET("/posts/:id/likes", likeHandler.GetPostLikeUsers)
            // auth.GET("/likes/posts", likeHandler.GetUserLikedPosts)

            // 评论相关
            // auth.POST("/posts/:post_id/comments", commentHandler.CreateComment)
            // auth.DELETE("/comments/:id", commentHandler.DeleteComment)
            // auth.GET("/posts/:post_id/comments", commentHandler.GetCommentList)
            // auth.POST("/comments/:id/like", commentHandler.LikeComment)
            // auth.DELETE("/comments/:id/like", commentHandler.UnlikeComment)
        // }
    // } */

    // 静态文件服务
    if cfg.Storage.Type == "local" {
        r.Static("/uploads", cfg.Storage.Local.Path)
    }

    // 启动服务器
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      r,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    // 优雅关闭
    go func() {
        logger.Info("Server starting", zap.String("port", port))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed to start", zap.Error(err))
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down server...")

    // 设置关闭超时
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown", zap.Error(err))
    }

    logger.Info("Server exiting")
}