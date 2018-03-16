package media

// 注释先记一下大概

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MaterialUpdater 用来生成json更新永久图文素材
type MaterialUpdater struct {
	MediaID  string   `json:"media_id,omitempty"`
	Index    int      `json:"indext,omitempty"`
	Articles *Article `json:"articles,omitempty"`
}

// WxMaterailUpdateNews 素材更新路径
const WxMaterailUpdateNews = "cgi-bin/material/update_news"

// UpdateMaterial 更新永久图文素材
func UpdateMaterial(host, path, accessToken string, materialUpdater *MaterialUpdater) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, path, accessToken)
	b, err := json.Marshal(materialUpdater)
	if err != nil {
		return nil, fmt.Errorf("update material news: %s: %s", materialUpdater.MediaID, err)
	}
	res, err := http.Post(URL, "application/json; charset=utf-8", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	b, err = readResponse(res)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("update material news: %s: %s", materialUpdater.MediaID, err)

	}
	return &resp, nil
}

func readResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return b, nil
}

// MaterialCounter 素材计数器
type MaterialCounter struct {
	VoiceCount int    `json:"voice_count,omitempty"`
	VideoCount int    `json:"video_count,omitempty"`
	ImageCount int    `json:"image_count,omitempty"`
	NewsCount  int    `json:"news_count,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// WxGetMaterialCount 获取永久素材数量的路径
const WxGetMaterialCount = "cgi-bin/material/get_materialcount"

// GetMaterialCount 获取永久素材数量
func GetMaterialCount(host, path, accessToken string) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, path, accessToken)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	b, err := readResponse(res)
	if err != nil {
		return nil, err
	}
	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, errors.New("get material count failed")
	}
	return &resp, nil
}

// MaterialListRequest 获取永久素材的列表
type MaterialListRequest struct {
	// Type 素材的类型，图片（image）、视频（video）、语音 （voice）、图文（news)
	Type string `json:"type"`
	// Offset 从全部素材的该偏移位置开始返回，0表示从第一个素材返回
	Offset int `json:"offset"`
	// Count 返回素材的数量，取值在1到20之间
	Count int `json:"count"`
}

// MaterialList 永久素材列表
type MaterialList struct {
	TotalCount int     `json:"total_count,omitempty"`
	ItemCount  int     `json:"item_count,omitempty"`
	Items      []*Item `json:"item,omitempty"`
	ErrCode    int     `json:"errcode,omitempty"`
	ErrMsg     string  `json:"errmsg,omitempty"`
}

// Item 素材列表项
type Item struct {
	MediaID    string   `json:"media_id,omitempty"`
	Content    *Content `json:"content,omitempty"`
	UpdateTime string   `json:"update_time,omitempty"`
	Name       string   `json:"name,omitempty"`
	URL        string   `json:"url,omitempty"`
}

// Content 图文内容
type Content struct {
	NewsItems []*NewsItem `json:"news_item,omitempty"`
}

// NewsItem 图文列表项
type NewsItem struct {
	// Title 标题
	Title string `json:"title,omitempty"`
	// ThumbMediaID 图文消息的封面素材ID
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
	// Author 作者
	Author string `json:"author,omitempty"`
	// Digest 图文消息的摘要，仅有单图消息才有摘要，多图文此处为空; 若该值为空，默认选取前64个字
	Digest string `json:"digest,omitempty"`
	// ShowCoverPic 是否显示封面，0：false，1：true
	ShowCoverPic int `json:"show_cover_pic,omitempty"`
	// Content 图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS,涉及图片url必须来源 "上传图文消息内的图片获取URL"接口获取。外部图片url将被过滤。
	Content string `json:"content,omitempty"`
	// ContentSourceURL 图文消息的原文地址，即点击“阅读原文”后的URL
	ContentSourceURL string `json:"content_source_url,omitempty"`
	// URL 图文页的URL，或者，当获取的列表是图片素材列表时，该字段是图片的URL
	URL string `json:"url,omitempty"`
	// NeedOpenComment 可以留言或者评论
	NeedOpenComment uint32 `json:"need_open_comment,omitempty"`
	// OnlyFansCanComment 只有公众号粉丝评价
	OnlyFansCanComment uint32 `json:"only_fans_can_comment,omitempty"`
}

// WxMaterailGetList 素材列表路径
const WxMaterailGetList = "cgi-bin/material/batchget_material"

// GetMaterialList 获取永久素材的列表
func GetMaterialList(host, path, accessToken string, req *MaterialListRequest) (*MaterialList, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, path, accessToken)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(URL, "application/json; charset=utf-8", bytes.NewReader(b))
	b, err = readResponse(res)
	if err != nil {
		return nil, err
	}
	var resp MaterialList
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

const (
	// WxOpenComment 打开评论
	WxOpenComment = "cgi-bin/comment/open"
	// wxCloseComment 关闭评论
	wxCloseComment = "cgi-bin/comment/close"
)

// TODO
// 添加评论的操作
