// Package main used for weixin (mp.weixin.qq.com)
package weixin

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	wxEncryptType          = "aes"
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
	// EncodingAESKey 旧的消息加密密钥
	OldEncodingAESKey string

	// 保存access_token
	accessToken string
	// access_token的有效期
	expires int
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
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("appid %s read response: %s", err)
	}
	defer res.Body.Close()

	var t Token
	if err = json.Unmarshal(b, &t); err != nil {
		return fmt.Errorf("appid %s unmarshal response json: %s",
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

// Sign 生成签名，ciphertext是空字符串时，只使用token, timestamp, nonce
func Sign(token, timestamp, nonce, ciphertext string) string {
	list := []string{token, timestamp, nonce}
	if ciphertext != "" {
		list = append(list, ciphertext)
	}
	// 排序token timestamp, nonce和ciphertext
	sort.Strings(list)
	h := sha1.New()
	for i := 0; i < len(list); i++ {
		fmt.Fprint(h, list[i])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// VerfiyWxToken 使用*Weixin.Token与timestamp、nonce的sha1值，与signature校验
func (wx *WeiXin) VerfiyWxToken(timestamp, nonce, signature, ciphertext string) bool {
	hashcode := Sign(wx.Token, timestamp, nonce, ciphertext)
	if hashcode == signature {
		return true
	}
	return false
}

func (wx *WeiXin) HandleEncryptEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("encrypt event --- %s\n", r)
	r.ParseForm()
	msgSignature := r.FormValue("msg_signature")
	encryptType := r.FormValue("encrypt_type")
	if msgSignature == "" || encryptType != wxEncryptType {
		wx.HandleEvent(w, r)
		return
	}
	// 校验消息是否来自微信服务器
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")

	if !wx.VerfiyWxToken(timestamp, nonce, signature, "") {
		// 只打印
		fmt.Println("verify event from weixin failed")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("handle encrypt message read body %s\n", err)
		fmt.Fprint(w, "")
		return
	}
	defer r.Body.Close()
	fmt.Printf("emsg:\n%s\n", body)
	var emsg EncryptMessage
	if err := xml.Unmarshal(body, &emsg); err != nil {
		fmt.Printf("handle encrypt message parse xml content first %s\n", err)
		fmt.Fprint(w, "")
		return
	}
	if !wx.VerfiyWxToken(timestamp, nonce, msgSignature, string(emsg.Encrypt)) {
		fmt.Println("handle weixin message verfiy encrypt message failed.\n")
		fmt.Fprint(w, "")
		return
	}

	key := wx.EncodingAESKey

	b, err := Decrypt(key, string(emsg.Encrypt))
	if err != nil {
		fmt.Printf("handle messsage: decrypt by current key %s\n", err)
		if wx.OldEncodingAESKey != "" {
			bo, err := Decrypt(wx.OldEncodingAESKey, string(emsg.Encrypt))
			if err != nil {
				fmt.Printf("handle messsage: decrypt by old key %s\n", err)
				fmt.Fprint(w, "")
			}
			key = wx.OldEncodingAESKey
			b = bo
		}
		return
	}

	fmt.Printf("decrypt emsg:\n%s\n", b)

	plaintext, appid, err := ParseEncryptMessage(b, wx.AppID)
	if err != nil {
		fmt.Printf("parse encrypt message %s\n", err)
		fmt.Fprint(w, "")
		return
	}
	fmt.Printf("decrypt emsg:\n%s\n%s\n", plaintext, appid)

	var msg Message
	if err := xml.Unmarshal(plaintext, &msg); err != nil {
		fmt.Printf("handle message: unmarshal decrypt plaintext %s\n", err)
		fmt.Fprint(w, "")
		return
	}

	s, err := newExampleMsg(string(msg.ToUserName), string(msg.FromUserName), "加密消息应答")
	if err != nil {
		fmt.Printf("handle message: make response to reply %s\n", err)
		fmt.Fprint(w, "")
		return
	}

	ciphertext, err := Encrypt(key, s, appid)
	if err != nil {
		fmt.Printf("handle message: encrypt response %s\n", err)
		fmt.Fprint(w, "")
		return
	}
	eres := NewEncryptResponse(wx.AppID, wx.Token, nonce, ciphertext)
	resp, err := xml.Marshal(eres)
	if err != nil {
		fmt.Printf("handle message marshal xml response %s\n", err)
		fmt.Fprint(w, "")
		return
	}
	fmt.Printf("%s\n", resp)
	w.Header().Set("Content-Type", "application/xml; encoding=utf-8")
	fmt.Fprintf(w, "%s", resp)
}

// HandleEvent 处理微信服务器验证token请求
//	TODO:
//	  1. 消息接收与转发队列
//	  2. 应答队列
func (wx *WeiXin) HandleEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("nomal --- %s\n", r)
	r.ParseForm()
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")

	if !wx.VerfiyWxToken(timestamp, nonce, signature, "") {
		// 只打印
		fmt.Println("verify weixin event failed.")
		return
	}
	switch r.Method {
	// GET方法用于微信服务器配置验证
	case "GET":
		echostr := r.FormValue("echostr")
		fmt.Fprintf(w, "%s", echostr)
	// POST方法接收明文事件
	case "POST":
		var msg Message
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("handle event read context %s", err)
			fmt.Fprint(w, "")
			return
		}
		defer r.Body.Close()
		if err := xml.Unmarshal(body, &msg); err != nil {
			fmt.Printf("handle event parse xml %s\n", err)
			fmt.Fprint(w, "")
			return
		}
		fmt.Printf("message: -----\n%#v\n", msg)
		fmt.Println("------")
		s, err := newExampleMsg(string(msg.ToUserName), string(msg.FromUserName), "回复应答消息")
		if err != nil {
			fmt.Printf("handle event new message for reply %s\n", err)
			fmt.Fprint(w, "")
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
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("get callback ip read body: %s", err)
	}
	if err = xml.Unmarshal(b, &ips); err != nil {
		return nil, fmt.Errorf("get callback ip address of weixin: %s", err)
	}
	return &ips, nil
}
