package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type WXLoginResponse struct {
    OpenID     string `json:"openid"`
    SessionKey string `json:"session_key"`
    UnionID    string `json:"unionid"`
    ErrCode    int    `json:"errcode"`
    ErrMsg     string `json:"errmsg"`
}

const (
    WXAppID     = "wxbec01cd32419625b"     // 替换为你的小程序AppID
    WXAppSecret = "3d2aaa51b7befff63ffaea64da2a5c49" // 替换为你的小程序AppSecret
)

func WXLogin(code string) (*WXLoginResponse, error) {
    url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        WXAppID, WXAppSecret, code)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var wxResp WXLoginResponse
    if err := json.NewDecoder(resp.Body).Decode(&wxResp); err != nil {
        return nil, err
    }

    if wxResp.ErrCode != 0 {
        return nil, errors.New(wxResp.ErrMsg)
    }

    return &wxResp, nil
}