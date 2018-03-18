package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// WxMaterailGet 获取永久图文素材路径
const WxMaterailGet = "cgi-bin/material/get_material"

// MaterialGetResponse 获取永久图文消息响应
type MaterialGetResponse struct {
	// Title 视频素材用到
	Title string `json:"title,omitempty"`
	// Description 视频素材用到
	Description string `json:"description,omitempty"`
	// DownURL 视频素材用到
	DownURL string `json:"down_url,omitempty"`
	// NewsItem 正常图文
	NewsItem []*Article `json:"news_item,omitempty"`
	ErrCode  int        `json:"errcode,omitempty"`
	ErrMsg   string     `json:"errmsg,omitempty"`
}

// GetMaterial 获取永久图文素材，返回的*MaterialGetResponse需要注意请求的是视频还是图文
func GetMaterial(host, accessToken, typ, mediaID, dir string) (filename string, resp *MaterialGetResponse, err error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, accessToken, WxMaterailGet)
	JSON := `{"media_id":"` + mediaID + `"}`
	res, err := http.Post(URL, "application/json; charset=utf-8", strings.NewReader(JSON))
	if err != nil {
		return "", nil, fmt.Errorf("when get material %s-%s", mediaID, err)
	}
	if res.StatusCode != 200 {
		return "", nil, fmt.Errorf("when get material %s+%s", mediaID, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, fmt.Errorf("when get material %s/%s", mediaID, err)
	}
	defer res.Body.Close()

	if err = json.Unmarshal(b, &resp); err != nil {
		if typ != "video" && typ != "news" {
			filename = filepath.Join(dir, mediaID)
			if err = ioutil.WriteFile(filename, b, 0640); err != nil {
				return "", nil, err
			}
			return filename, nil, nil
		}
	}
	return
}
