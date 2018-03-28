package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"qingtao/weixin/mp/internal/tools"
)

/*
以下引用自微信公众平台的文档：
获取用户基本信息(UnionID机制)

在关注者与公众号产生消息交互后，公众号可获得关注者的OpenID（加密后的微信号，每个用户对每个公众号的OpenID是唯一的。对于不同公众号，同一用户的openid不同）。公众号可通过本接口来根据OpenID获取用户基本信息，包括昵称、头像、性别、所在城市、语言和关注时间。

请注意，如果开发者有在多个公众号，或在公众号、移动应用之间统一用户帐号的需求，需要前往微信开放平台（open.weixin.qq.com）绑定公众号后，才可利用UnionID机制来满足上述需求。

UnionID机制说明：

开发者可通过OpenID来获取用户基本信息。特别需要注意的是，如果开发者拥有多个移动应用、网站应用和公众帐号，可通过获取用户基本信息中的unionid来区分用户的唯一性，因为只要是同一个微信开放平台帐号下的移动应用、网站应用和公众帐号，用户的unionid是唯一的。换句话说，同一用户，对同一个微信开放平台下的不同应用，unionid是相同的。
*/

// User 公众平台的粉丝用户, 具体文档参考：https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140839
type User struct {
	Subscribe uint32 `json:"subscribe,omitempty"`
	OpenID    string `json:"openid,omitempty"`
	NickName  string `json:"nickname,omitempty"`
	Sex       uint32 `json:"sex,omitempty"`
	Language  string `json:"language,omitempty"`
	City      string `json:"city,omitempty"`
	Province  string `json:"province,omitempty"`
	Country   string `json:"country,omitempty"`
	// HeadImgURL 用户头像，最后一个数值代表正方形头像大小(有0、46、64、96、132数值可选，0代表640*640正方形头像),用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImgURL     string `json:"headimgurl,omitempty"`
	SubscribeTime  uint32 `json:"subscribe_time,omitempty"`
	UnionID        string `json:"unionid,omitempty"`
	Remark         string `json:"remark,omitempty"`
	GroupID        uint32 `json:"groupid,omitempty"`
	TagIDList      []int  `json:"tagid_list,omitempty"`
	SubscribeScene string `json:"subscribe_scene,omitempty"`
	QrScene        uint32 `json:"qr_scene,omitempty"`
	QrSceneStr     string `json:"qr_scene_str,omitempty"`

	// ErrCode 错误代码
	ErrCode uint32 `json:"errcode,omitempty"`
	// Errmsg 错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// WxUserInfoPath 微信公众平台获取用户基本信息的API路径
const WxUserInfoPath = "cgi-bin/user/info"

// GetUserInfo 根据API接口通过GET方法获取用户基本信息
func GetUserInfo(host, accessToken, openid, lang string) (*User, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s&openid=%s", host, WxUserInfoPath, accessToken, openid)
	if lang != "" {
		URL = URL + "&lang=" + lang
	}
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	var user User
	if err = tools.UnmarshalJSON(res, &user); err != nil {
		return nil, fmt.Errorf("when get user info %s: error: %s", openid, err)
	}

	/*
		b, err := tools.ReadResponse(res)
		if err != nil {
			return nil, fmt.Errorf("When get user info %s", err)
		}
		if err = json.Unmarshal(b, &user);
		err != nil {
			return nil, fmt.Errorf("when get user info %s: error: %s", openid, err)
		}
	*/

	return &user, nil
}

// Users 批量获取用户基本信息时API返回的结构体
type Users struct {
	UserInfoList []*User `json:"user_info_list,omitempty"`
	ErrCode      uint32  `json:"errcode,omitempty"`
	ErrMsg       string  `json:"errmsg,omitempty"`
}

// WxUsersInfoPath 微信公众平台获取用户基本信息的API路径
const WxUsersInfoPath = "cgi-bin/user/info/batchget"

// Item 用于批量获取用户信息
type Item struct {
	OpenID string `json:"openid,omitempty"`
	Lang   string `json:"lang,omitempty"`
}

// GetUsersInfo 批量获取用户基本信息
func GetUsersInfo(host, accessToken string, userlist []*Item) (*Users, error) {
	URL := fmt.Sprintf("https://%s/%s?access_token=%s", host, WxUsersInfoPath, accessToken)
	var userList = struct {
		UserList []*Item `json:"user_lsit"`
	}{
		UserList: userlist,
	}

	b, err := json.Marshal(userList)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(URL, JSONContentType, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var users Users
	if err = tools.UnmarshalJSON(res, &users); err != nil {
		return nil, err
	}
	return &users, nil
}
