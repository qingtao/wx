package users

// 2018/03/18
// TODO: 添加完成用户管理后，开始实际测试各API调用情况
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const JSONContentType = "application/json; charset=utf-8"

// WxTagsCreate 创建用户标签API
const WxTagsCreate = `cgi-bin/tags/create`

type TagResponse struct {
	Tag     *Tag   `json:"tag,omitempty"`
	ErrCode uint32 `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

type Tag struct {
	ID   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TagsResponse struct {
	Tags    []*Tag `json:"tags,omitempty"`
	ErrCode uint32 `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// CreateTag 添加标签
func CreateTag(host, accessToken, tagname string) (*TagResponse, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxTagsCreate, accessToken)
	JSON := fmt.Sprintf(`{"tag":{"name":"%s"}}`, tagname)
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tres TagResponse
	if err = json.Unmarshal(b, &tres); err != nil {
		return nil, err
	}
	return &tres, nil
}

type Response struct {
	ErrCode uint32 `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}
