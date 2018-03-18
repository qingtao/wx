package cs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// postAcount 微信客服消息接口管理客服帐号，action: add/update/del
func postAcount(host, action, accessToken string, acc *Account) (*Response, error) {
	b, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("https://%s/%s/%s/?access_token=%s",
		host, WxKfPath, action, accessToken)
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

// AddAccount 新增客服帐号
func AddAccount(host, accessToken string, acc *Account) (*Response, error) {
	return postAcount(host, WxKfAdd, accessToken, acc)
}

// UpdateAccount 修改客服帐号
func UpdateAccount(host, accessToken string, acc *Account) (*Response, error) {
	return postAcount(host, WxKfUpdate, accessToken, acc)
}

// DelAccount 删除客服帐号
func DelAccount(host, accessToken string, acc *Account) (*Response, error) {
	return postAcount(host, WxKfDel, accessToken, acc)
}

// UploadHeadImage 上传客服头像
func UploadHeadImage(host, accessToken, account, filename string) (*Response, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	// 检查图片扩展名
	if ext != ".jpg" {
		return nil, fmt.Errorf("image must be one of .jpg")
	}
	// 取文件状态，得到文件名称和大小
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("upload %s %s", filename, err)
	}

	var buf = new(bytes.Buffer)
	multiWriter := multipart.NewWriter(buf)
	// multipart的文件，名称“media”是微信要求的参数
	w, err := multiWriter.CreateFormFile("media", stat.Name())
	if err != nil {
		return nil, err
	}
	// 打开文件
	fr, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	// 文件内容写入到w->buf
	if _, err = io.Copy(w, fr); err != nil {
		return nil, err
	}
	multiWriter.Close()
	// 获取http头部的Content-Type
	contentType := multiWriter.FormDataContentType()

	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%skf_account=%s",
		host, WxKfPath, WxKfHeadImg, accessToken, account)
	res, err := http.Post(uri, contentType, buf)
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

// GetList 获取所有客户帐号
func GetList(host, accessToken string) (*Lists, error) {
	uri := fmt.Sprintf("https://%s/%s?access_token=%s",
		host, WxKfGetKfList, accessToken)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("getkflist %s", err)
	}
	defer res.Body.Close()
	var list Lists
	if err = json.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("getkflist %s", err)
	}
	return &list, nil
}

// SendMessage 发送客服消息
func SendMessage(host, accessToken string, msg *Message) (*Response, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("send custom message %s", err)
	}
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", host,
		WxKfSend, accessToken)
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

// SendTyping 发送输入状态
func SendTyping(host, accessToken, toUser string) (*Response, error) {
	typing := `{"touser":"` + toUser + `", "command":"typing"}`
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", host,
		WxKftyping, accessToken)
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
