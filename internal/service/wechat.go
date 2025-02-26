// internal/service/wechat.go

package service

import (
	"campus-wall/internal/model"
	"campus-wall/pkg/auth"
	"campus-wall/pkg/errors"
	"campus-wall/pkg/wechat"
	"fmt"

	"gorm.io/gorm"
)

type WeChatService struct {
    db          *gorm.DB
    tokenMaker  *auth.JWTMaker
    miniProgram *wechat.MiniProgram
}

type WXDecryptRequest struct {
    EncryptedData string `json:"encryptedData" binding:"required"`
    IV            string `json:"iv" binding:"required"`
    SessionKey    string `json:"sessionKey" binding:"required"`
}

func NewWeChatService(db *gorm.DB, tokenMaker *auth.JWTMaker, miniProgram *wechat.MiniProgram) *WeChatService {
    return &WeChatService{
        db:          db,
        tokenMaker:  tokenMaker,
        miniProgram: miniProgram,
    }
}

type WXLoginRequest struct {
    Code     string `json:"code" binding:"required"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
    Gender   int    `json:"gender"`
    Introduction string `json:"introduction"`
}

func (s *WeChatService) MiniProgramLogin(req *WXLoginRequest) (*AuthResponse, error) {
    // 通过code获取openid和session_key
    session, err := s.miniProgram.Code2Session(req.Code)
    fmt.Printf("session: %+v\n", session)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to get session", err)
    }

    var user model.User
    err = s.db.Where("open_id = ?", session.OpenID).First(&user).Error

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            // 新用户注册
            user = model.User{
                OpenID:       session.OpenID,
                UnionID:      session.UnionID,
                SessionKey:   session.SessionKey,
                Nickname:     req.Nickname,
                Avatar:       req.Avatar,
                Gender:      req.Gender,
            }
            if err := s.db.Create(&user).Error; err != nil {
                return nil, errors.NewInternalServerError("failed to create user", err)
            }
        } else {
            return nil, errors.NewInternalServerError("database error", err)
        }
    } else {
        // 更新用户信息
        updates := map[string]interface{}{
            "session_key": session.SessionKey,
            "nickname":    req.Nickname,
            "avatar":      req.Avatar,
            "gender":      req.Gender,
            "introduction":   req.Introduction,
        }
        if err := s.db.Model(&user).Updates(updates).Error; err != nil {
            return nil, errors.NewInternalServerError("failed to update user", err)
        }
    }

    // 生成token
    token, err := s.tokenMaker.CreateToken(user.ID, user.OpenID)
    fmt.Printf("token: %+v\n", token)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to create token", err)
    }

    return &AuthResponse{
        User:        user,
        AccessToken: token,
    }, nil
}

func (s *WeChatService) DecryptUserInfo(req *WXDecryptRequest) (*wechat.UserInfo, error) {
    crypto := wechat.NewWXBizDataCrypt(s.miniProgram.AppID, req.SessionKey)
    userInfo, err := crypto.DecryptUserInfo(req.EncryptedData, req.IV)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to decrypt user info", err)
    }
    return userInfo, nil
}

func (s *WeChatService) DecryptPhoneNumber(req *WXDecryptRequest) (*wechat.PhoneInfo, error) {
    crypto := wechat.NewWXBizDataCrypt(s.miniProgram.AppID, req.SessionKey)
    phoneInfo, err := crypto.DecryptPhoneNumber(req.EncryptedData, req.IV)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to decrypt phone number", err)
    }
    return phoneInfo, nil
}