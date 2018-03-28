package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// WxUserUpdaterMark 设置用户备注
const WxUserUpdaterMark = "cgi-bin/user/info/updateremark"

// MarkUser 备注用户, openid 是用户标识符，根据微信公众平台的api说明，remark不能超过30个字符
func MarkUser(host, accessToken, openid, remark string) (*Response, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxUserUpdaterMark, accessToken)
	JSON := `{"openid":"` + openid + `","remark":"` + remark + `"}`
	res, err := http.Post(URL, JSONContentType, strings.NewReader(JSON))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
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
