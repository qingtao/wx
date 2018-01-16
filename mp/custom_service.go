package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// 客服帐号管理接口path
	wxKfPath    = "customservice/kfaccount"
	wxKfAdd     = "add"
	wxKfUpdate  = "update"
	wxKfDel     = "del"
	wxKfHeadImg = "uploadheadimg"
	// 获取客服帐号和发送信息路径
	wxKfGetAll = "cgi-bin/customservice/getkflist"
	wxKfSend   = "cgi-bin/message/custom/send"
	// 发送输入状态接口
	wxKftyping = "cgi-bin/message/cutom/typing"
)

// Request 帐号管理
type KfAccount struct {
	Kf_account string `json:"kf_account"`
	Nickname   string `json:"nickname"`
	Password   string `json:"password"`
}

// Response 微信返回错误码和信息
type Response struct {
	// Errcode 错误代码
	Errcode int `json:"errcode"`
	// Errmsg 错误消息
	Errmsg string `json:"errmsg"`
}

// postKfAcount 微信客服消息接口管理客服帐号，action: add/update/del
func (wx *WeiXin) postKfAcount(action string, acc *KfAccount) (*Response, error) {
	b, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("https://%s/%s/%s/?access_token=%s",
		wx.Host, wxKfPath, action, wx.accessToken)
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%s kfacount %s", action, err)
	}
	defer res.Body.Close()
	var status Response
	if err = json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("%s kfacount %s", action, err)
	}
	return &status, nil
}

// AddKfAccount 新增客服帐号
func (wx *WeiXin) AddKfAccount(acc *KfAccount) (*Response, error) {
	return wx.postKfAcount(wxKfAdd, acc)
}

// UpdateKfAccount 修改客服帐号
func (wx *WeiXin) UpdateKfAccount(acc *KfAccount) (*Response, error) {
	return wx.postKfAcount(wxKfUpdate, acc)
}

// DelKfAccount 删除客服帐号
func (wx *WeiXin) DelKfAccount(acc *KfAccount) (*Response, error) {
	return wx.postKfAcount(wxKfDel, acc)
}

// UploadKfHeadImage 上传客服头像，未实现
func (wx *WeiXin) UploadKfHeadImage(kfaccount string, r io.Reader) (*Response, error) {
	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%skf_account=%s",
		wx.Host, wxKfPath, wxKfHeadImg, wx.accessToken, kfaccount)
	res, err := http.Post(uri, "image/jpeg", r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var status Response
	if err := json.Unmarshal(b, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

type Kflist struct {
	Kf_account    string `json:"kf_account"`
	Kf_nick       string `json:"kf_nick"`
	Kf_id         string `json:'kf_id"`
	Kf_headimgurl string `json:"kf_headimgurl"`
}

type AllKfaccount struct {
	Kf_list []*Kflist `json:"kf_lsit"`
	Errcode int       `json:"errcode"`
	Errmsg  string    `json:"errmsg"`
}

func (wx *WeiXin) GetKfList() (*AllKfaccount, error) {
	uri := fmt.Sprintf("https://%s/%s?access_token=%s",
		wx.Host, wxKfGetAll, wx.accessToken)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("getkflist %s", err)
	}
	defer res.Body.Close()
	var list AllKfaccount
	if err = json.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("getkflist %s", err)
	}
	return &list, nil
}

type CustomMessage struct {
	ToUser          string                 `json:"touser"`
	MsgType         string                 `json:"msgtype"`
	Text            *Text                  `json:"text,omitempty"`
	Image           *CustomMedia           `json:"image,omitempty"`
	Video           *CustomMedia           `json:"video,omitempty"`
	Voice           *CustomMedia           `json:"voice,omitempty"`
	Music           *CustomMusic           `json:"music,omitempty"`
	MpNews          *CustomMedia           `json:"mpnews,omitempty"`
	News            *CustomNews            `json:"news,omitempty"`
	WxCard          *CustomWxCard          `json:"wxcard,omitempty"`
	MiniProgramPage *CustomMiniProgramPage `json:"miniprogrampage,omitempty"`
	CustomService   *CustomService         `json:"customservice,omitempty"`
}

type Text struct {
	Content string `josn:"content,omitempty"`
}

type CustomMedia struct {
	MediaId      string `json:"media_id"`
	ThumbMediaId string `json:"thumb_media_id,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
}

type CustomMusic struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	MusicUrl     string `json:"musicurl,omitempty"`
	HqMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id,omitempty"`
}

type CustomNews struct {
	Articles []*CustomNew `json:"articles"`
}

type CustomNew struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	PicUrl      string `json:"picurl,omitempty"`
}

type CustomWxCard struct {
	CardId string `json:"card_id,omitempty"`
}

func (wx *WeiXin) SendCustomMessage(v interface{}) (*Response, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("send custom message %s", err)
	}
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", wx.Host,
		wxKfSend, wx.accessToken)
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("send custom message %s", err)
	}
	var status Response
	if err = json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("send custom message, read response %s", err)
	}
	return &status, nil
}

type CustomMiniProgramPage struct {
	Title        string `json:"title,omitempty"`
	AppId        string `json:"appid,omitempty"`
	PagePath     string `json:"pagepath,omitempty"`
	ThumbMediaId string `json:"thumb_media_id,omitempty"`
}

type CustomService struct {
	KfAcount string `json:"kf_account,omitempty"`
}

func (wx *WeiXin) SendCustomTyping(touser string) (*Response, error) {
	typing := `{"touser":"` + touser + `", "command":"typing"}`
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", wx.Host,
		wxKftyping, wx.accessToken)
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader([]byte(typing)))

	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("send typing %s", err)
	}
	var status Response
	if err = json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("send typing read response %s", err)
	}
	return &status, nil
}
