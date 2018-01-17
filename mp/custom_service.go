package mp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"qingtao/weixin/mp/kf"
)

// postAcount 微信客服消息接口管理客服帐号，action: add/update/del
func (wx *WeiXin) postAcount(action string, acc *kf.Account) (*kf.Response, error) {
	b, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("https://%s/%s/%s/?access_token=%s",
		wx.Host, kf.WxKfPath, action, wx.accessToken)
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
	var status kf.Response
	if err = json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("%s kfacount %s", action, err)
	}
	return &status, nil
}

// AddAccount 新增客服帐号
func (wx *WeiXin) AddAccount(acc *kf.Account) (*kf.Response, error) {
	return wx.postAcount(kf.WxKfAdd, acc)
}

// UpdateAccount 修改客服帐号
func (wx *WeiXin) UpdateAccount(acc *kf.Account) (*kf.Response, error) {
	return wx.postAcount(kf.WxKfUpdate, acc)
}

// DelAccount 删除客服帐号
func (wx *WeiXin) DelAccount(acc *kf.Account) (*kf.Response, error) {
	return wx.postAcount(kf.WxKfDel, acc)
}

// UploadHeadImage 上传客服头像，未实现
func (wx *WeiXin) UploadHeadImage(account string, r io.Reader) (*kf.Response, error) {
	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%skf_account=%s",
		wx.Host, kf.WxKfPath, kf.WxKfHeadImg, wx.accessToken, account)
	res, err := http.Post(uri, "image/jpeg", r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var status kf.Response
	if err := json.Unmarshal(b, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// GetList 获取所有客户帐号
func (wx *WeiXin) GetList() (*kf.Lists, error) {
	uri := fmt.Sprintf("https://%s/%s?access_token=%s",
		wx.Host, kf.WxKfGetKfList, wx.accessToken)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("getkflist %s", err)
	}
	defer res.Body.Close()
	var list kf.Lists
	if err = json.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("getkflist %s", err)
	}
	return &list, nil
}

// SendMessage 发送客服消息
func (wx *WeiXin) SendMessage(msg *kf.Message) (*kf.Response, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("send custom message %s", err)
	}
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", wx.Host,
		kf.WxKfSend, wx.accessToken)
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("send custom message %s", err)
	}
	var status kf.Response
	if err = json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("send custom message, read response %s", err)
	}
	return &status, nil
}

// SendTyping 发送输入状态
func (wx *WeiXin) SendTyping(touser string) (*kf.Response, error) {
	typing := `{"touser":"` + touser + `", "command":"typing"}`
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", wx.Host,
		kf.WxKftyping, wx.accessToken)
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader([]byte(typing)))

	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("send typing %s", err)
	}
	var status kf.Response
	if err = json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("send typing read response %s", err)
	}
	return &status, nil
}
