package database

import (
	"campus-wall/internal/model"
	"campus-wall/pkg/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=%t",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
        // cfg.Charset,
        cfg.ParseTime,
        // cfg.Loc,
    )

	// 连接数据库
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

	// 获取底层的sqlDB
    // sqlDB, err := db.DB()
    // if err != nil {
    //     return nil, err
    // }

	// 运行数据库迁移
    // if err := RunMigrations(sqlDB, "campus-wall/migrations"); err != nil {
    //     return nil, fmt.Errorf("failed to run migrations: %v", err)
    // }

    // 设置连接池
    // sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    // sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    // sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

    // 自动迁移
    err = db.AutoMigrate(
        &model.User{},
        // &model.Post{},
        // &model.Comment{},
        // &model.Like{},
        // &model.Report{},
    )
    if err != nil {
        return nil, err
    }

    return db, nil
}