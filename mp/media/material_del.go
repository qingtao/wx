package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response 响应错误代码和消息
type Response struct {
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// WxMaterialDel 删除永久图文的路径
const WxMaterialDel = "cgi-bin/material/del_material"

// DeleteMaterial 删除永久素材
func DeleteMaterial(host, path, mediaID, accessToken string) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, path, accessToken)
	JSON := `{"media_id":"` + mediaID + `"}`
	res, err := http.Post(URL, "application/json; charset=utf-8", strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("delete material failed %s", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
