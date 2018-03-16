package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MaterialUpdater 用来生成json更新永久图文素材
type MaterialUpdater struct {
	MediaID  string     `json:"media_id,omitempty"`
	Index    int        `json:"indext,omitempty"`
	Articles []*Article `json:"articles,omitempty"`
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
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("update material news: %s: %s", materialUpdater.MediaID, res.Status)
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("update material news: %s: %s", materialUpdater.MediaID, err)
	}
	defer res.Body.Close()

	var resp Response
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("update material news: %s: %s", materialUpdater.MediaID, err)

	}
	return &resp, nil
}
