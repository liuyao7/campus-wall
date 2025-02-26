// internal/service/user.go

package service

import (
	"campus-wall/internal/model"
	"campus-wall/pkg/auth"
	"campus-wall/pkg/errors"
	"campus-wall/pkg/storage"
	"path"
	"strings"

	"gorm.io/gorm"

	"mime/multipart"
)

type UserService struct {
    db         *gorm.DB
    tokenMaker *auth.JWTMaker
    storage storage.Storage
}

func NewUserService(db *gorm.DB, tokenMaker *auth.JWTMaker, storage storage.Storage) *UserService {
    return &UserService{
        db:         db,
        tokenMaker: tokenMaker,
        storage: storage,
    }
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=32"`
    Password string `json:"password" binding:"required,min=6,max=32"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    User        model.User `json:"user"`
    AccessToken string     `json:"access_token"`
}

func (s *UserService) Register(req *RegisterRequest) (*AuthResponse, error) {
    // 检查用户名是否已存在
    var existingUser model.User
    if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
        return nil, errors.NewBadRequestError("username already exists")
    }

    // 加密密码
    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to hash password", err)
    }

    // 创建用户
    user := model.User{
        Username: req.Username,
        Password: hashedPassword,
    }

    if err := s.db.Create(&user).Error; err != nil {
        return nil, errors.NewInternalServerError("failed to create user", err)
    }

    // 生成token
    token, err := s.tokenMaker.CreateToken(user.ID, user.OpenID)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to create token", err)
    }

    return &AuthResponse{
        User:        user,
        AccessToken: token,
    }, nil
}

func (s *UserService) Login(req *LoginRequest) (*AuthResponse, error) {
    // 查找用户
    var user model.User
    if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
        return nil, errors.NewUnauthorizedError("invalid username or password")
    }

    // 验证密码
    if err := auth.CheckPassword(req.Password, user.Password); err != nil {
        return nil, errors.NewUnauthorizedError("invalid username or password")
    }

    // 生成token
    token, err := s.tokenMaker.CreateToken(user.ID, user.OpenID)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to create token", err)
    }

    return &AuthResponse{
        User:        user,
        AccessToken: token,
    }, nil
}

func (s *UserService) UpdateUser(userID uint, req *model.UserUpdateRequest) (*model.User, error) {
    // 1. 查找用户
    var user model.User
    if err := s.db.First(&user, userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.NewNotFoundError("user not found")
        }
        return nil, errors.NewInternalServerError("database error", err)
    }

    // 2. 构建更新数据
    updates := make(map[string]interface{})
    
    if req.Nickname != nil {
        updates["nickname"] = *req.Nickname
    }
    if req.Gender != nil {
        updates["gender"] = *req.Gender
    }
    if req.Birthday != nil {
        updates["birthday"] = *req.Birthday
    }
    if req.Introduction != nil {
        updates["introduction"] = *req.Introduction
    }
    if req.Location != nil {
        updates["location"] = *req.Location
    }
    if req.Email != nil {
        // 检查邮箱是否已被使用
        var count int64
        s.db.Model(&model.User{}).Where("email = ? AND id != ?", *req.Email, userID).Count(&count)
        if count > 0 {
            return nil, errors.NewBadRequestError("email already exists")
        }
        updates["email"] = *req.Email
    }
    if req.Phone != nil {
        // 检查手机号是否已被使用
        var count int64
        s.db.Model(&model.User{}).Where("phone = ? AND id != ?", *req.Phone, userID).Count(&count)
        if count > 0 {
            return nil, errors.NewBadRequestError("phone number already exists")
        }
        updates["phone"] = *req.Phone
    }

    // 3. 更新用户信息
    if err := s.db.Model(&user).Updates(updates).Error; err != nil {
        return nil, errors.NewInternalServerError("failed to update user", err)
    }

    // 4. 重新获取用户信息
    if err := s.db.First(&user, userID).Error; err != nil {
        return nil, errors.NewInternalServerError("failed to get updated user", err)
    }

    return &user, nil
}

// UploadAvatar 上传头像
func (s *UserService) UploadAvatar(userID uint, file *multipart.FileHeader) (*model.User, error) {
    // 1. 验证文件类型
    if err := validateImageFile(file); err != nil {
        return nil, err
    }

    // 2. 查找用户
    var user model.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return nil, errors.NewNotFoundError("user not found")
    }

    // 3. 如果用户已有头像，删除旧头像
    if user.Avatar != "" {
        _ = s.storage.Delete(user.Avatar)
    }

    // 4. 上传新头像
    avatarURL, err := s.storage.Upload(file, "avatars")
    if err != nil {
        return nil, errors.NewInternalServerError("failed to upload avatar", err)
    }

    // 5. 更新用户头像
    if err := s.db.Model(&user).Update("avatar", avatarURL).Error; err != nil {
        return nil, errors.NewInternalServerError("failed to update avatar", err)
    }

    return &user, nil
}

// validateImageFile 验证图片文件
func validateImageFile(file *multipart.FileHeader) error {
    // 检查文件大小（例如最大2MB）
    if file.Size > 2*1024*1024 {
        return errors.NewBadRequestError("file size exceeds maximum limit")
    }

    // 检查文件类型
    ext := strings.ToLower(path.Ext(file.Filename))
    allowedTypes := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
        ".gif":  true,
    }

    if !allowedTypes[ext] {
        return errors.NewBadRequestError("invalid file type")
    }

    return nil
}