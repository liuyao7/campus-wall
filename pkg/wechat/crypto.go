// pkg/wechat/crypto.go

package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// WXBizDataCrypt 微信数据解密
type WXBizDataCrypt struct {
    AppID      string
    SessionKey string
}

// UserInfo 解密后的用户信息
type UserInfo struct {
    OpenID    string `json:"openId"`
    UnionID   string `json:"unionId"`
    NickName  string `json:"nickName"`
    Gender    int    `json:"gender"`
    City      string `json:"city"`
    Province  string `json:"province"`
    Country   string `json:"country"`
    AvatarURL string `json:"avatarUrl"`
    Language  string `json:"language"`
    Watermark struct {
        Timestamp int64  `json:"timestamp"`
        AppID     string `json:"appid"`
    } `json:"watermark"`
}

// PhoneInfo 解密后的手机号信息
type PhoneInfo struct {
    PhoneNumber     string `json:"phoneNumber"`
    PurePhoneNumber string `json:"purePhoneNumber"`
    CountryCode     string `json:"countryCode"`
    Watermark      struct {
        Timestamp int64  `json:"timestamp"`
        AppID     string `json:"appid"`
    } `json:"watermark"`
}

func NewWXBizDataCrypt(appID, sessionKey string) *WXBizDataCrypt {
    return &WXBizDataCrypt{
        AppID:      appID,
        SessionKey: sessionKey,
    }
}

// Decrypt 解密数据
func (w *WXBizDataCrypt) Decrypt(encryptedData, iv string) ([]byte, error) {
    sessionKey, err := base64.StdEncoding.DecodeString(w.SessionKey)
    if err != nil {
        return nil, err
    }

    encData, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return nil, err
    }

    ivBytes, err := base64.StdEncoding.DecodeString(iv)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(sessionKey)
    if err != nil {
        return nil, err
    }

    mode := cipher.NewCBCDecrypter(block, ivBytes)
    mode.CryptBlocks(encData, encData)

    // PKCS#7 unpadding
    padding := int(encData[len(encData)-1])
    if padding < 1 || padding > 32 {
        return nil, errors.New("invalid padding size")
    }
    return encData[:len(encData)-padding], nil
}

// DecryptUserInfo 解密用户信息
func (w *WXBizDataCrypt) DecryptUserInfo(encryptedData, iv string) (*UserInfo, error) {
    data, err := w.Decrypt(encryptedData, iv)
    if err != nil {
        return nil, err
    }

    var userInfo UserInfo
    err = json.Unmarshal(data, &userInfo)
    if err != nil {
        return nil, err
    }

    if userInfo.Watermark.AppID != w.AppID {
        return nil, errors.New("invalid app id in watermark")
    }

    return &userInfo, nil
}

// DecryptPhoneNumber 解密手机号
func (w *WXBizDataCrypt) DecryptPhoneNumber(encryptedData, iv string) (*PhoneInfo, error) {
    data, err := w.Decrypt(encryptedData, iv)
    if err != nil {
        return nil, err
    }

    var phoneInfo PhoneInfo
    err = json.Unmarshal(data, &phoneInfo)
    if err != nil {
        return nil, err
    }

    if phoneInfo.Watermark.AppID != w.AppID {
        return nil, errors.New("invalid app id in watermark")
    }

    return &phoneInfo, nil
}