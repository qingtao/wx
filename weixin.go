// Package main used for weixin (mp.weixin.qq.com)
package weixin

import (
	"crypto/cipher"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
)

const (
	wxConfigHost           = "Host"
	wxConfigAppID          = "AppID"
	wxConfigAppSecret      = "AppSecret"
	wxConfigToken          = "Token"
	wxConfigEncodingAESKey = "EncodingAESKey"
)

// WeiXin 微信公众号配置参数
type WeiXin struct {
	XMLName xml.Name `xml:"weixin" json:"-"`
	// Host 微信服务器主机名
	Host string
	// 微信开发者ID
	// AppID 应用ID
	AppID string
	// AppSecret 应用密钥
	AppSecret string
	// Token 令牌
	Token string
	// EncodingAESKey 消息加密密钥, 即消息加解密Key，长度固定为43个字符，从a-z,A-Z,0-9共62个字符中选取。由开发者在创建公众号插件时填写，也可申请修改
	EncodingAESKey string

	// 保存access_token
	accessToken string
	// access_token的有效期
	expires int
	// aes加解密
	block    cipher.Block
	blockold cipher.Block
}

// New 读取filename文件, 生成新的*WeiXin, 失败时返回error非空
func New(filename string) (*WeiXin, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read weixin config: %s", err)
	}
	var wx WeiXin
	if err = xml.Unmarshal(b, &wx); err != nil {
		return nil, fmt.Errorf("read weixin config: %s", err)
	}
	return &wx, nil
}

// CreateWeiXinFile 创建weixin参数的xml文件, 按照申请的微信公众平台修改参数
func CreateWeiXinFile(filename string) error {
	wx := &WeiXin{
		Host:           "api.weixin.qq.com",
		AppID:          "appid",
		AppSecret:      "appsecret",
		Token:          "token",
		EncodingAESKey: "-",
	}
	b, err := xml.MarshalIndent(wx, "", "  ")
	if err != nil {
		return err
	}
	// 添加xml.Header到文件第一行
	b = append([]byte(xml.Header), b[0:]...)
	return ioutil.WriteFile(filename, b, 0664)
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

// Unmarshal 使用函数f([]byte, interface{})解析微信的消息结构
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
//	TODO:
//	  * 将获取和更新access_token分离成服务
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
	wx.accessToken = t.AccessToken
	wx.expires = t.ExpiresIn
	return nil
}

// Sign 生成签名，ciphertext是空字符串时，只是用token, timestamp, nonce
func Sign(token, timestamp, nonce, ciphertext string) string {
	list := []string{token, timestamp, nonce}
	if ciphertext != "" {
		list = append(list, ciphertext)
	}
	sort.Strings(list)
	h := sha1.New()
	for i := 0; i < len(list); i++ {
		fmt.Fprint(h, list[i])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// VerfiyWxToken 使用Weixin.Token与timestamp、nonce的sha1值，与signature校验
func (wx *WeiXin) VerfiyWxToken(timestamp, nonce, signature, ciphertext string) bool {
	hashcode := Sign(wx.Token, timestamp, nonce, ciphertext)
	if hashcode == signature {
		return true
	}
	return false
}

// HandleEvent 处理微信服务器验证token请求
//	TODO:
//	  1. 消息接收与转发队列
//	  2. 应答队列
func (wx *WeiXin) HandleEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r)
	r.ParseForm()
	signature := r.FormValue("signature")
	msg_signature := r.FormValue("msg_signature")
	if signature == "" && msg_signature == "" {
		return
	}
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")

	switch r.Method {
	// GET方法用于微信服务器配置验证
	case "GET":
		if !wx.VerfiyWxToken(timestamp, nonce, signature, "") {
			// 只打印
			fmt.Println("handle weixin message verfiy failed.")
			fmt.Fprint(w, "")
			return
		}
		echostr := r.FormValue("echostr")
		fmt.Fprintf(w, "%s", echostr)
		// POST
	case "POST":
		if msg_signature != "" {
			var msg EncryptMessage
			if err := Unmarshal(r.Body, &msg, xml.Unmarshal); err != nil {
				fmt.Printf("HandleMsg: %s\n", err)
				return
			}
			if !wx.VerfiyWxToken(timestamp, nonce, msg_signature, string(msg.Encrypt)) {
				fmt.Println("handle weixin message verfiy failed.")
				fmt.Fprint(w, "")
				return
			}
			block, err := NewCipherBlock(wx.EncodingAESKey)
			if err != nil {
				fmt.Printf("HandleMsg: cipher block %s\n", err)
				fmt.Fprint(w, "success")
				return
			}
			bs, err := Decrypt(block, string(msg.Encrypt))
			if err != nil {
				fmt.Printf("HandleMsg: decrypt %s\n", err)
				fmt.Fprint(w, "success")
				return
			}

			xmltext, appid, err := ParseEncryptMessage(bs, wx.AppID)
			if err != nil {
				fmt.Printf("HandleMsg: decrypt %s\n", err)
				fmt.Fprint(w, "success")
				return
			}

			var msg1 Message
			if err := xml.Unmarshal(xmltext, &msg1); err != nil {
				fmt.Printf("HandleMsg: unmarshal plaintext %s\n", err)
				fmt.Fprint(w, "success")
				return
			}

			s, err := newExampleMsg(string(msg1.ToUserName), string(msg1.FromUserName))
			if err != nil {
				fmt.Printf("HandleMsg: make msg %s\n", err)
				fmt.Fprint(w, "success")
				return
			}

			ctext, err := Encrypt(block, s, appid)
			if err != nil {
				fmt.Printf("HandleMsg: encrypt %s\n", err)
				fmt.Fprint(w, "success")
				return
			}
			eres := NewEncryptResponse(wx.AppID, wx.Token, nonce, ctext)
			resp, err := xml.Marshal(eres)
			if err != nil {
				fmt.Printf("HandleMsg: encrypt %s\n", err)
				fmt.Fprint(w, "success")
				return
			}
			fmt.Printf("%s\n", resp)
			fmt.Fprintf(w, "%s", resp)
			return
		}

		var msg Message
		if err := Unmarshal(r.Body, &msg, xml.Unmarshal); err != nil {
			fmt.Printf("HandleMsg: %s\n", err)
			return
		}
		fmt.Printf("%#v\n", msg)
		fmt.Println("------")
		s, err := newExampleMsg(string(msg.ToUserName), string(msg.FromUserName))
		if err != nil {
			fmt.Printf("HandleMsg: %s\n", err)
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
