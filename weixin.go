// Package main used for weixin (mp.weixin.qq.com)
package wx

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// TODO: replace to config
const (
	wxConfigLines          = 5
	wxConfigHost           = "Host"
	wxConfigAppID          = "AppID"
	wxConfigAppSecret      = "AppSecret"
	wxConfigToken          = "Token"
	wxConfigEncodingAESKey = "EncodingAESKey"
)

// WeiXin 微信公众号配置参数
type WeiXin struct {
	// Host 微信服务器主机名
	Host string
	// 微信开发者ID
	// AppID 应用ID
	AppID string
	// AppSecret 应用密钥
	AppSecret string
	// Token 令牌
	Token string

	// EncodingAESKey 消息加密密钥
	EncodingAESKey string

	accessToken string
	expires     string
}

// New 读取key.txt, 生成新的*WeiXin, 失败时返回error非空
func New(filename string) (*WeiXin, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read weixin config: %s", err)
	}
	s := strings.Replace(string(b), "\r\n", "\n", -1)
	lines := strings.Split(s, "\n")
	fmt.Printf("%s\n--\n%s\n", b, s)
	if len(lines) != wxConfigLines {
		return nil, fmt.Errorf("weixin config lines must be %d", wxConfigLines)
	}
	var wx = new(WeiXin)
	for i := 0; i < len(lines); i++ {
		args := strings.Split(strings.TrimSpace(lines[i]), "=")
		if len(args) != 2 {
			return nil, fmt.Errorf("weixin config invalid, line: %d - %s",
				i, lines[i])
		}
		switch args[0] {
		case wxConfigHost:
			wx.Host = args[1]
		case wxConfigAppID:
			wx.AppID = args[1]
		case wxConfigAppSecret:
			wx.AppSecret = args[1]
		case wxConfigToken:
			wx.Token = args[1]
		case wxConfigEncodingAESKey:
			wx.EncodingAESKey = args[1]
		}
	}
	return wx, nil
}

const (
	// WxTokenPath 获取access_token时使用
	WxTokenPath = "cgi-bin/token"
	// WxGrantType 获取access_token时使用
	WxGrantType = "client_credential"
	// WxGetCallBackIP 获取微信服务器IP地址时使用
	WxGetCallBackIpPath = "cgi-bin/getcallbackip"
)

// Token 从微信公众平台申请Token的响应
type Token struct {
	// AccessToken 是成功时服务器响应中access_token
	AccessToken string `json:"access_token,omitempty"`
	// ExpiresIn 是成功是服务器响应中的过期时间，单位秒
	ExpiresIn int `json:"expires_in,omitempty"`

	//错误代码和信息
	// ErrMsg 是失败时服务器响应中的错误代码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 是失败时服务器响应中的错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

func Unmarshal(rc io.ReadCloser, a interface{}, f func([]byte, interface{}) error) error {
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("read http response: %s", err)
	}
	defer rc.Close()
	// 解析json
	if err := f(b, a); err != nil {
		return fmt.Errorf("unmarshal: %s", err)
	}
	return nil
}

// GetAccessToken 获取微信公众平台的access_token，
// error为nil时，access_token = Token.AccessToken和expires_in = Token.ExpiresIn
func (wx *WeiXin) GetAccessToken() error {
	// 组合URL
	uri := fmt.Sprintf("https://%s/%s?grant_type=%s&appid=%s&secret=%s",
		wx.Host, WxTokenPath, WxGrantType, wx.AppID, wx.AppSecret)
	res, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("appid %s get access_token: %s",
			wx.AppID, err)
	}
	var t Token
	if err = Unmarshal(res.Body, &t, json.Unmarshal); err != nil {
		return fmt.Errorf("appid %s unmarshal response: %s",
			wx.AppID, err)
	}

	// 检查t.AccessToken为空，返回nil，错误代码和错误信息
	if t.AccessToken == "" {
		return fmt.Errorf("appid %s get access_token errcode: %d, errmsg: %s", wx.AppID, t.ErrCode, t.ErrMsg)
	}
	return nil
}

// 计算微信服务器发送的token，timestamp、nonce的sha1散列值，与signature校验
func (wx *WeiXin) VerfiyWxToken(timestamp, nonce, signature string) bool {
	list := []string{wx.Token, timestamp, nonce}
	// 排序字符串
	sort.Strings(list)

	h := sha1.New()
	for i := 0; i < len(list); i++ {
		fmt.Fprint(h, list[i])
	}
	// 16进制字符串
	hashcode := fmt.Sprintf("%x", h.Sum(nil))
	if hashcode == signature {
		return true
	}
	return false
}

// 处理微信服务器验证token请求
func (wx *WeiXin) HandleEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r)
	r.ParseForm()
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")

	ok := true
	if !wx.VerfiyWxToken(timestamp, nonce, signature) {
		ok = false
		// 只打印
		fmt.Println("handle weixin info verfiy failed.")
		return
	}
	switch r.Method {
	// GET方法用于微信服务器配置验证
	case "GET":
		if !ok {
			fmt.Fprint(w, "")
			return
		}
		echostr := r.FormValue("echostr")
		fmt.Fprintf(w, "%s", echostr)
		// POST
	case "POST":
		if !ok {
			return
		}
		var wxinfo Info
		if err := Unmarshal(r.Body, &wxinfo, xml.Unmarshal); err != nil {
			fmt.Printf("HandleWxInfo: %s\n", err)
			return
		}
		fmt.Printf("%#v\n", wxinfo)
		fmt.Println("------")
		s, err := newExampleInfo(string(wxinfo.ToUserName), string(wxinfo.FromUserName))
		if err != nil {
			fmt.Printf("HandleWxInfo: %s\n", err)
			fmt.Fprint(w, "success")
			return
		}
		fmt.Printf("%s\n", s)
		w.Header().Set("Content-Type", "application/xml; encoding=utf-8")
		fmt.Fprintf(w, "%s", s)
	}
}

// CallBackIP微信服务器IP地址
type CallBackIP struct {
	//微信服务器IP地址列表
	IPList []string `json:"ip_list"`
	// ErrMsg 是失败时服务器响应中的错误代码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 是失败时服务器响应中的错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// GetCallBackIP 获取微信服务器IP地址
// 如果公众号基于安全等考虑，需要获知微信服务器的IP地址列表，
// 以便进行相关限制，可以通过该接口获得微信服务器IP地址列表或者IP网段信息。
func (wx *WeiXin) GetCallBackIP() (*CallBackIP, error) {
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", wx.Host,
		WxGetCallBackIpPath, wx.accessToken)
	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("get callback ip address of weixin: %s", err)
	}
	var ips CallBackIP
	if err = Unmarshal(res.Body, &ips, json.Unmarshal); err != nil {
		return nil, fmt.Errorf("get callback ip address of weixin: %s", err)
	}
	return &ips, nil
}
