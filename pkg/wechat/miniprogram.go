// pkg/wechat/miniprogram.go

package wechat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
    codeToSessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type MiniProgram struct {
    AppID     string
    AppSecret string
}

type Code2SessionResponse struct {
    OpenID     string `json:"openid"`
    SessionKey string `json:"session_key"`
    UnionID    string `json:"unionid"`
    ErrCode    int    `json:"errcode"`
    ErrMsg     string `json:"errmsg"`
}

func NewMiniProgram(appID, appSecret string) *MiniProgram {
    return &MiniProgram{
        AppID:     appID,
        AppSecret: appSecret,
    }
}

func (mp *MiniProgram) Code2Session(code string) (*Code2SessionResponse, error) {
    // 打印一下mp
    log.Printf("MiniProgram: %+v", mp)

    url := fmt.Sprintf(codeToSessionURL, mp.AppID, mp.AppSecret, code)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    log.Printf("Received response with status code: %d", resp.StatusCode)


    var result Code2SessionResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    log.Printf("Decoded response: %+v", result)

    if result.ErrCode != 0 {
        return nil, fmt.Errorf("WeChat error: %d - %s", result.ErrCode, result.ErrMsg)
    }


    return &result, nil
}