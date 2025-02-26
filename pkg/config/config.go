// pkg/config/config.go

package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig     `mapstructure:"jwt"`
    Redis    RedisConfig   `mapstructure:"redis"`
    Logger   LoggerConfig  `mapstructure:"logger"`
    WeChat   WeChatConfig `mapstructure:"wechat"`
    Storage  StorageConfig `mapstructure:"storage"`
    Swagger struct {
        Enabled bool   `mapstructure:"enabled"`
        Host    string `mapstructure:"host"`
        BasePath string `mapstructure:"base_path"`
    } `mapstructure:"swagger"`
}

type StorageConfig struct {
    Type  string             `mapstructure:"type"`
    Local LocalStorageConfig `mapstructure:"local"`
    OSS     OSSStorageConfig   `mapstructure:"oss"`
}

type OSSStorageConfig struct {
    Endpoint string `mapstructure:"endpoint"`
    AccessKeyID      string `mapstructure:"access_key_id"`
    AccessKeySecret  string `mapstructure:"access_key_secret"`
    BucketName        string `mapstructure:"bucket_name"`
}

type LocalStorageConfig struct {
    Path string `mapstructure:"root_path"`
    BaseURL string `mapstructure:"base_url"`
}

type WeChatConfig struct {
    MiniProgram MiniProgramConfig `mapstructure:"miniprogram"`
}

type MiniProgramConfig struct {
    AppID     string `mapstructure:"app_id"`
    AppSecret string `mapstructure:"app_secret"`
}

type RateLimitConfig struct {
    Requests int           `mapstructure:"max_requests"`
    Burst    int           `mapstructure:"burst"`
}

type ServerConfig struct {
    Port         string        `mapstructure:"port"`
    RateLimit    RateLimitConfig           `mapstructure:"rate_limit"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
    ParseTime bool   `mapstructure:"parse_time"`
}

type JWTConfig struct {
    SecretKey     string        `mapstructure:"secret_key"`
    TokenDuration time.Duration `mapstructure:"token_duration"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize     int          `mapstructure:"pool_size"`
    MinIdleConns int          `mapstructure:"min_idle_conns"`
    MaxRetries   int          `mapstructure:"max_retries"`
    DialTimeout  time.Duration `mapstructure:"dial_timeout"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type LoggerConfig struct {
    Level      string `mapstructure:"level"`
    Filename   string `mapstructure:"filename"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}


func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    config := &Config{}
    if err := viper.Unmarshal(config); err != nil {
        return nil, err
    }

    return config, nil
}

// GetDSN 返回数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host,
        c.Port,
        c.Username,
        c.Password,
        c.DBName,
        c.SSLMode,
    )
}