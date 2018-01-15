package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// menu paths
const (
	// WxMenuPath 自定义菜单接口的Path
	WxMenuPath = "cgi-bin/menu"
	// WxMenuCreate 创建菜单的Path
	WxMenuCreate = "create"
	// WxMenuDelete 删除菜单的Path
	WxMenuDelete = "delete"
	// WxMenuGet 查询菜单的Path
	WxMenuGet = "get"
	// WxMenuAddConditional 查询个性化菜单的Path
	WxMenuAddConditional = "addconditional"
	// WxMenuDelConditional 删除个性化菜单的Path
	WxMenuDelConditional = "delconditional"
	// WxTryMatch 测试个性化菜单匹配结果
	WxMenuTryMatch = "trymatch"
	//WxGetCurrentSelfMenu 获取自定义菜单配置接口
	WxGetCurrentSelfMenu = "cgi-bin/get_current_selfmenu_info"
)

// 自定义菜单创建接口的按钮类型，根据微信公众平台文档的描述
//   https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141013
const (
	// WxMenuClickType 点击推事件用户点击click类型按钮后，微信服务器会通过消息接口推送消息类型为event的结构给开发者（参考消息接口指南），并且带上按钮中开发者填写的key值，开发者可以通过自定义的key值与用户进行交互
	WxMenuClickType = "click"
	// WxMenuViewType 跳转URL用户点击view类型按钮后，微信客户端将会打开开发者在按钮中填写的网页URL，可与网页授权获取用户基本信息接口结合，获得用户基本信息
	WxMenuViewType = "view"
	// WxMenuScanCodePush 扫码推事件用户点击按钮后，微信客户端将调起扫一扫工具，完成扫码操作后显示扫描结果（如果是URL，将进入URL），且会将扫码的结果传给开发者，开发者可以下发消息
	WxMenuScanCodePush = "scancode_push"
	// WxMenuScanCodeWaitMsg 扫码推事件且弹出“消息接收中”提示框用户点击按钮后，微信客户端将调起扫一扫工具，完成扫码操作后，将扫码的结果传给开发者，同时收起扫一扫工具，然后弹出“消息接收中”提示框，随后可能会收到开发者下发的消息
	WxMenuScanCodeWaitMsg = "scancode_waitmsg"
	// WxMenuPicSysPhoto 弹出系统拍照发图用户点击按钮后，微信客户端将调起系统相机，完成拍照操作后，会将拍摄的相片发送给开发者，并推送事件给开发者，同时收起系统相机，随后可能会收到开发者下发的消息
	WxMenuPicSysPhoto = "pic_sysphoto"
	// WxMenuPicPhotoOrAlbum 弹出拍照或者相册发图用户点击按钮后，微信客户端将弹出选择器供用户选择“拍照”或者“从手机相册选择”。用户选择后即走其他两种流程
	WxMenuPicPhotoOrAlbum = "pic_photo_or_album"
	// WxMenuPicWeixin 弹出微信相册发图器用户点击按钮后，微信客户端将调起微信相册，完成选择操作后，将选择的相片发送给开发者的服务器，并推送事件给开发者，同时收起相册，随后可能会收到开发者下发的消息
	WxMenuPicWeixin = "pic_weixin"
	// WxMenuLocationSelect 出地理位置选择器用户点击按钮后，微信客户端将调起地理位置选择工具，完成选择操作后，将选择的地理位置发送给开发者的服务器，同时收起位置选择工具，随后可能会收到开发者下发的消息
	WxMenuLocationSelect = "location_select"
	// WxMenuMediaID 下发消息（除文本消息）用户点击media_id类型按钮后，微信服务器会将开发者填写的永久素材id对应的素材下发给用户，永久素材类型可以是图片、音频、视频、图文消息。请注意：永久素材id必须是在“素材管理/新增永久素材”接口上传后获得的合法id
	WxMenuMediaID = "media_id"
	// WxMenuViewLimited 跳转图文消息URL用户点击view_limited类型按钮后，微信客户端将打开开发者在按钮中填写的永久素材id对应的图文消息URL，永久素材类型只支持图文消息。请注意：永久素材id必须是在“素材管理/新增永久素材”接口上传后获得的合法id
	WxMenuViewLimited = "view_limited"
	// WxMenuMiniProgram 可能是小程序
	WxMenuMiniProgram = "miniprogram"
)

// Button 自定义菜单的button
type Button struct {
	// Name 菜单标题，不超过16个字节，子菜单不超过60个字节
	Name string `json:"name"`
	// Type 菜单的响应类型，view表示网页类型，click表示点击类型，miniprogram表示小程序
	Type string `json:"type,omitempty"`
	// Key click等点击类型必须，菜单KEY值，用于消息接口推送，不超过128字节
	Key string `json:"key,omitempty"`
	// URL view/miniprogram等类型必须，网页链接，用户点击惨淡可以打开链接，不超过1024字节，type为miniprogram时，不支持小程序的老版本客户端打开本URL
	URL string `json:"url,omitempty"`
	// MediaId media_id类型和view_limited类型必须，调用新增永久素材接口返回的合法media_id
	MediaId string `json:"media_id,omitempty"`
	// AppId miniprogram类型必须，小程序的appid
	AppId string `json:"appid,omitempty"`
	// miniprogram必须，小程序的页面路径
	PagePath string `json:"pagepath,omitempty"`
	// SubButton 子菜单
	SubButton []*Button `json:"sub_button,omitempty"`
	// NewsInfo 图文消息的信息
	NewsInfo *NewsInfo `json:"news_info,omitempty"`
}

// MenuOfConditional 查询自定义菜单接口返回数据
type MenuOfConditional struct {
	// Menu 菜单
	Menu *Menu `json:"menu,omitempty"`
	// ConditionalMenu 个性化菜单
	ConditionalMenu *ConditionalMenu `json:"conditionalmenu,omitempty"`
	// ErrCode 自定义菜单错误码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 自定义菜单错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// ConditionalMenu 个性化菜单
type ConditionalMenu struct {
	// Button 按钮
	Button []*Button `json:"button,omitempty"`
	// MatchRule
	MatchRule *MatchRule `json:"matchrule,omitempty"`
	// ErrCode 自定义菜单错误码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 自定义菜单错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// NewsInfo 在公众平台官网通过网站功能发布菜单包含此字段
type NewsInfo struct {
	// List 菜单列表
	List []*NewsInfoList `json:"list,omitempty"`
}

// NewsInfoList 包含在NewsInfo图文消息的信息中
type NewsInfoList struct {
	// Title 	图文消息的标题
	Title string `json:"title,omitempty"`
	// Author 作者
	Author string `json:"author,omitempty"`
	// Digest 摘要
	Digest string `json:"digest,omitempty"`
	// ShowCover是否显示封面，0为不显示，1为显示
	ShowCover int64 `json:"show_conver,omitempty"`
	// CoverURL 封面图片的URL
	CoverURL string `json:"cover_url,omitempty"`
	// ContentURL 正文的URL
	ContentURL string `json:"content_url,omitempty"`
	// SourceURL 原文的URL，若置空则无查看原文入口
	SourceURL string `json:"source_url,omitempty"`
}

// Menu 自定义菜单
type Menu struct {
	// Button 一级菜单组
	Button []*Button `json:"button,omitempty"`
	// Menu 个性化菜单的ID，只有返回的是个性化菜单是才不为空
	MenuId string `json:"menuid,omitempty"`
	// ErrCode 自定义菜单错误码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 自定义菜单错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// Response 微信返回的响应信息, 用在只返回状态的请求
type Response struct {
	// ErrCode 自定义菜单错误码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 自定义菜单错误信息
	ErrMsg string `json:"errmsg,omitempty"`
	MenuId string `json:"menuid,omitempty"`
}

// post 提交自定义菜单操作
func (wx *WeiXin) post(action string, menu interface{}) (*Response, error) {
	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%s",
		wx.Host, WxMenuPath, action, wx.accessToken)
	b, err := json.Marshal(menu)
	if err != nil {
		return nil, fmt.Errorf("appid %s json marshal: %s\n", wx.AppId, err)
	}
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("appid %s post create menu: %s\n",
			wx.AppId, err)
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("appid %s menu response read body: %s", wx.AppId, err)
	}
	defer res.Body.Close()
	var wxinfo Response
	if err = json.Unmarshal(b, &wxinfo); err != nil {
		return nil, fmt.Errorf("appid %s json unmarshal menu: %s", wx.AppId, err)
	}
	return &wxinfo, nil
}

// CreateMenu 创建自定义菜单
func (wx *WeiXin) CreateMenu(menu *Menu) (*Response, error) {
	return wx.post(WxMenuCreate, menu)
}

// GetMenu 查询自定义菜单
func (wx *WeiXin) GetMenu(accessToken string) (*MenuOfConditional, error) {
	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%s",
		wx.Host, WxMenuPath, WxMenuGet, wx.accessToken)
	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("appid %s get menu: %s", wx.AppId, err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("appid %s menu response read body: %s", wx.AppId, err)
	}
	defer res.Body.Close()

	var menu MenuOfConditional
	if err = json.Unmarshal(b, &menu); err != nil {
		return nil, fmt.Errorf("appid %s get menu unmarshal response: %s", err)
	}

	return &menu, nil
}

// deleteMenu 提交删除自定义菜单请求
func (wx *WeiXin) deleteMenu(menuid string) (*Response, error) {
	action := WxMenuDelete
	if menuid != "" {
		action = WxMenuDelConditional
	}
	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%s",
		wx.Host, WxMenuPath, action, wx.accessToken)

	var res = new(http.Response)
	var err error
	if menuid != "" {
		s := `{"menuid":"` + menuid + `"}`
		res, err = http.Post(uri, "application/json; charset=utf-8",
			bytes.NewReader([]byte(s)))
	} else {
		res, err = http.Get(uri)
	}
	if err != nil {
		return nil, fmt.Errorf("appid %s the request of delete menu %s",
			wx.AppId, err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("appid %s del menu response read body: %s", wx.AppId, err)
	}
	defer res.Body.Close()

	var wxinfo Response
	if err := json.Unmarshal(b, &wxinfo); err != nil {
		return nil, fmt.Errorf("appid %s unmarshal response when delete menu: %s", wx.AppId, err)
	}
	return &wxinfo, nil
}

// DeleteMenu 删除自定义菜单
func (wx *WeiXin) DeleteMenu() (*Response, error) {
	return wx.deleteMenu("")
}

// MatchRule 菜单匹配规则
type MatchRule struct {
	// TagId 用户标签的id，可通过用户标签管理接口获取
	TagId string `json:"tag_id,omitempty"`
	// Sex 性别：男（1）女（2），不填则不做匹配
	Sex string `json:"sex,omitempty"`
	// Country 国家信息，是用户在微信中设置的地区，具体请参考地区信息表
	Country string `json:"country,omitempty"`
	// Province 省份信息，是用户在微信中设置的地区，具体请参考地区信息表
	Province string `json:"province,omitempty"`
	// City 城市信息，是用户在微信中设置的地区，具体请参考地区信息表
	City string `json:"city,omitempty"`
	// ClientPlatformType 客户端版本，当前只具体到系统型号:
	// IOS(1), Android(2),Others(3)，不填则不做匹配
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	// Language 语言信息，是用户在微信中设置的语言，具体请参考语言表:
	// 1、简体中文 "zh_CN" 2、繁体中文TW "zh_TW" 3、繁体中文HK "zh_HK";
	// 4、英文 "en" 5、印尼 "id" 6、马来 "ms" 7、西班牙 "es" 8、韩国 "ko";
	// 9、意大利 "it" 10、日本 "ja" 11、波兰 "pl" 12、葡萄牙 "pt";
	// 13、俄国 "ru" 14、泰文 "th" 15、越南 "vi" 16、阿拉伯语 "ar";
	// 17、北印度 "hi" 18、希伯来 "he" 19、土耳其 "tr" 20、德语 "de";
	// 21、法语 "fr"
	Language string `json:"language,omitempty"`
}

// CreateCustomMenu 创建个性化菜单
func (wx *WeiXin) CreateCustomMenu(menu *ConditionalMenu) (*Response, error) {
	return wx.post(WxMenuAddConditional, menu)
}

// DeleteCustomMenu删除个性化菜单
func (wx *WeiXin) DeleteCustomMenu(menuid string) (*Response, error) {
	return wx.deleteMenu(menuid)
}

// TryCustomMenu 测试个性化菜单匹配结果
func (wx *WeiXin) TryCustomMenu(userid string) (*Menu, error) {
	uri := fmt.Sprintf("https://%s/%s/%s?access_token=%s", wx.Host,
		WxMenuPath, WxMenuTryMatch, wx.accessToken)
	s := `{"user_id":"` + userid + `"}`
	res, err := http.Post(uri, "application/json; charset=utf-8",
		bytes.NewReader([]byte(s)))
	if err != nil {
		return nil, fmt.Errorf("appid %s trymatch custom menu: %s",
			wx.AppId, err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("appid %s trymatch read body: %s", wx.AppId, err)
	}
	defer res.Body.Close()

	var wxinfo Menu
	if err = json.Unmarshal(b, &wxinfo); err != nil {
		return nil, fmt.Errorf("appid %s parse content when trymatch custom menu: %s", wx.AppId, err)
	}
	return &wxinfo, nil
}

// MenuConfig 获取自定义菜单配置接口
type CurrentSelfMenu struct {
	// IsMenuOpen 菜单是否开启，0代表未开启，1代表开启
	IsMenuOpen int `json:"is_menu_open"`
	// SelfMenuInfo 菜单信息
	SelfMenuInfo []*Button `json:"self_menu_info"`
	// ErrCode 自定义菜单错误码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 自定义菜单错误信息
	ErrMsg string `json:"errmsg,omitempty"`
}

// GetCurrentMenu 获取自定义菜单配置接口
func (wx *WeiXin) GetCurrentSelfMenu() (*CurrentSelfMenu, error) {
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", wx.Host,
		WxGetCurrentSelfMenu, wx.accessToken)
	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("appid %s get_current_selfmenu_info: %s",
			wx.AppId, err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("appid %s getcurrentmenu read body: %s", wx.AppId, err)
	}
	defer res.Body.Close()
	var menu CurrentSelfMenu
	if err = json.Unmarshal(b, &menu); err != nil {
		return nil, fmt.Errorf("appid %s get_current_selfmenu_info unmarshal response: %s", wx.AppId, err)
	}
	return &menu, nil
}
