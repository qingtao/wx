package users

// 2018/03/18
// TODO: 添加完成用户管理后，开始实际测试各API调用情况
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response 只包含errcode和errmsg的响应信息
type Response struct {
	ErrCode uint32 `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// JSONContentType HTTP POST中的Content-Type
const JSONContentType = "application/json; charset=utf-8"

// WxTagsCreate 创建用户标签API
const WxTagsCreate = `cgi-bin/tags/create`

// TagResponse 创建标签时返回的响应结构
type TagResponse struct {
	Tag     *Tag   `json:"tag,omitempty"`
	ErrCode uint32 `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// Tag 标签字段
type Tag struct {
	ID   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// changeTag 修改标签
func changeTag(host, accessToken, action, tagname string, id int) (*TagResponse, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, action, accessToken)
	JSON := fmt.Sprintf(`{"tag":{"id":%d,name":"%s"}}`, id, tagname)
	if tagname == "" {
		JSON = fmt.Sprintf(`{"tag":{"id":%d}}`, id)
	} else if id < 3 {
		JSON = fmt.Sprintf(`{"tag":{"name":"%s"}}`, tagname)
	}
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var tres TagResponse
	if err = json.Unmarshal(b, &tres); err != nil {
		return nil, err
	}
	return &tres, nil
}

// CreateTag 添加标签
func CreateTag(host, accessToken, tagname string) (*TagResponse, error) {
	return changeTag(host, accessToken, WxTagsCreate, tagname, 0)
}

// WxTagsGet 获取已创建的标签API
const WxTagsGet = "cgi-bin/tags/get"

// TagsResponse 获取标签列表时的响应
type TagsResponse struct {
	Tags    []*Tag `json:"tags,omitempty"`
	ErrCode uint32 `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// GetTags 获取已有标签
func GetTags(host, accessToken string) (*TagsResponse, error) {
	res, err := http.Get(fmt.Sprintf("https://%s/%s?access_token=%s", host, WxTagsGet, accessToken))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var tsres TagsResponse
	if err = json.Unmarshal(b, &tsres); err != nil {
		return nil, err
	}
	return &tsres, nil
}

// WxTagsUpdate 编辑标签API
const WxTagsUpdate = "cgi-bin/tags/update"

// UpdateTag 修改标签，id不能是0/1/2
func UpdateTag(host, accessToken, tagname string, id int) (*TagResponse, error) {
	// 检查id是否大于等于3，实际是微信公众平台会返回错误代码:45058
	if id < 3 {
		return nil, fmt.Errorf("ID of tag must be greater or equal to three, but it is %d", id)
	}
	return changeTag(host, accessToken, WxTagsUpdate, tagname, id)
}

// WxTagsDelete 删除标签API
const WxTagsDelete = "cgi-bin/tags/delete"

// DeleteTag 删除标签
func DeleteTag(host, accessToken string, id int) (*TagResponse, error) {
	// 检查id是否大于等于3，实际是微信公众平台会返回错误代码:45058
	if id < 3 {
		return nil, fmt.Errorf("ID of tag must be greater or equal to three, but it is %d", id)
	}
	return changeTag(host, accessToken, WxTagsDelete, "", id)
}

// UserOfTag 标签下的粉丝列表
type UserOfTag struct {
	Count      uint32  `json:"count,omitempty"`
	Data       []*Data `json:"data,omitempty"`
	NextOpenID string  `json:"next_openid,omitempty"`
}

// Data in UsersOfTag
type Data struct {
	OpenID []string `json:"openid"`
}

// WxGetTagUsers 获取标签下粉丝列表API
const WxGetTagUsers = "cgi-bin/user/tag/get"

// GetUsersOfTag 获取标签下的用户
func GetUsersOfTag(host, accessToken, next string, id int) (*UserOfTag, error) {
	JSON := fmt.Sprintf(`{"next_openid":"%s","tagid":%d}`, next, id)
	if next == "" {
		JSON = fmt.Sprintf(`{"tagid":%d}`, id)
	}
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxGetTagUsers, accessToken)
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var users UserOfTag
	if err = json.Unmarshal(b, &users); err != nil {
		return nil, err
	}
	return &users, nil
}

// WxBatchTagging 批量打标签API
const WxBatchTagging = "tags/members/batchtagging"

// BatchTag 提交的批量打标签数据
type BatchTag struct {
	TagID      uint32   `json:"tagid,omitempty"`
	OpenIDList []string `json:"openid_list,omitempty"`
}

// batchTagging 批量操作标签
func batchTagging(host, accessToken, action string, btag *BatchTag) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, action, accessToken)
	b, err := json.Marshal(btag)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(URL, JSONContentType, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchTagging 批量打标签
func BatchTagging(host, accessToken string, btag *BatchTag) (*Response, error) {
	return batchTagging(host, accessToken, WxBatchTagging, btag)
}

// WxUnBatchTagging 批量取消标签的API
const WxUnBatchTagging = "cgi-bin/tags/members/batchuntagging"

// UnBatchTagging 批量取消标签
func UnBatchTagging(host, accessToken string, btag *BatchTag) (*Response, error) {
	return batchTagging(host, accessToken, WxUnBatchTagging, btag)
}

// UserTagsList 获取用户所属的标签列表, 一个用户可以最多有20个标签
type UserTagsList struct {
	TagIDList []uint32 `json:"tagi_list,omitempty"`
}

// WxTagsGetIDList 获取用户身上的标签API
const WxTagsGetIDList = "cgi-bin/tags/getidlist"

// GetTagsOfUser 获取用户所属标签
func GetTagsOfUser(host, accessToken, openid string) (*UserTagsList, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxTagsGetIDList, accessToken)
	JSON := `{"openid":"` + openid + `"}`
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp UserTagsList
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
