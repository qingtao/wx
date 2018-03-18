package cs

const (
	// WxKfPath 客服帐号管理接口path
	WxKfPath = "customservice/kfaccount"
	// WxKfAdd 添加
	WxKfAdd = "add"
	// WxKfUpdate 修改
	WxKfUpdate = "update"
	// WxKfDel 删除
	WxKfDel = "del"
	// WxKfHeadImg 头像图片
	WxKfHeadImg = "uploadheadimg"
	// WxKfGetKfList 获取客服帐号和发送信息路径
	WxKfGetKfList = "cgi-bin/customservice/getkflist"
	// WxKfSend 发送消息
	WxKfSend = "cgi-bin/message/custom/send"
	// WxKftyping 发送输入状态接口
	WxKftyping = "cgi-bin/message/cutom/typing"
)

// Account 帐号管理
type Account struct {
	KfAccount string `json:"kf_account"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
}

// Response 微信返回错误码和信息
type Response struct {
	// Errcode 错误代码
	Errcode int `json:"errcode"`
	// Errmsg 错误消息
	Errmsg string `json:"errmsg"`
}

// List getkflist返回的客户帐号信息
type List struct {
	KfAccount    string `json:"kf_account,omitempty"`
	KfNick       string `json:"kf_nick,omitempty"`
	KfID         string `json:"kf_id,omitempty"`
	KfHeadimgurl string `json:"kf_headimgurl,omitempty"`
}

// Lists getkflist返回的客户帐号信息, 错误时返回错误码和错误信息
type Lists struct {
	KfList  []*List `json:"kf_lsit,omitempty"`
	Errcode int     `json:"errcode,omitempty"`
	Errmsg  string  `json:"errmsg,omitempty"`
}

// Message 客服消息
type Message struct {
	// ToUser 消息接收方
	ToUser string `json:"touser"`
	// MsgType 消息类型
	MsgType string `json:"msgtype"`
	// Text 文本消息，MsgType为text
	Text *Text `json:"text,omitempty"`
	// Image 图片信息， MsgType为image
	Image *Media `json:"image,omitempty"`
	// Video 视频消息，MsgType为video
	Video *Media `json:"video,omitempty"`
	// Voice 语音消息，MsgType为voice
	Voice *Media `json:"voice,omitempty"`
	// Music 音乐消息，MsgType为music
	Music *Music `json:"music,omitempty"`
	// MpNews 图文消息，MsgType为mpnews
	MpNews *Media `json:"mpnews,omitempty"`
	// News 图文消息，点击跳转到外链
	News *News `json:"news,omitempty"`
	// WxCard 微信卡卷
	WxCard *WxCard `json:"wxcard,omitempty"`
	// MiniProgramPage 小程序页面
	MiniProgramPage *MiniProgramPage `json:"miniprogrampage,omitempty"`
	// CustomService 需要以某个客服帐号来发消息（在微信6.0.2及以上版本中显示自定义头像），则需在JSON数据包的后半部分加入customservice参数
	CustomService *CustomService `json:"customservice,omitempty"`
}

// Text 文本消息
type Text struct {
	Content string `json:"content,omitempty"`
}

// Media 多媒体消息
type Media struct {
	MediaID      string `json:"media_id,omitempty"`
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
}

// NewMedia 新建多媒体消息
func NewMedia(mediaid, title, desc, thumb string) *Media {
	return &Media{mediaid, thumb, title, desc}
}

// Music 音乐消息
type Music struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	MusicURL     string `json:"musicurl,omitempty"`
	HqMusicURL   string `json:"hqmusicurl"`
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
}

// NewMusic 新建音乐消息
func NewMusic(title, desc, musicurl, hqmusicurl, thumb string) *Music {
	return &Music{title, desc, musicurl, hqmusicurl, thumb}
}

// News 图文消息
type News struct {
	Articles []*Article `json:"articles,omitempty"`
}

// Article 图文文章
type Article struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	PicURL      string `json:"picurl,omitempty"`
}

// NewArticle 新建图文文章
func NewArticle(title, desc, URL, picurl string) *Article {
	return &Article{title, desc, URL, picurl}
}

// WxCard 微信卡卷
type WxCard struct {
	CardID string `json:"card_id,omitempty"`
}

// NewWxCard 新建客服消息微信卡卷消息
func NewWxCard(id string) *WxCard {
	return &WxCard{id}
}

// MiniProgramPage 微信小程序页面
type MiniProgramPage struct {
	Title        string `json:"title,omitempty"`
	AppID        string `json:"appid,omitempty"`
	PagePath     string `json:"pagepath,omitempty"`
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
}

// NewMiniProgramPage 新建小程序页面消息
func NewMiniProgramPage(title, appid, pagepath, thumb string) *MiniProgramPage {
	return &MiniProgramPage{title, appid, pagepath, thumb}
}

// CustomService 发送客服消息使用的客服帐号名称
type CustomService struct {
	KfAccount string `json:"kf_account,omitempty"`
}

// SetCustomService 设置使用哪个客服帐号发送消息
func (msg *Message) SetCustomService(account string) {
	if account != "" {
		msg.CustomService = &CustomService{account}
	}
}

// NewTextMessage 创建文本客服消息
func NewTextMessage(touser, fromuser, content string) *Message {
	msg := &Message{
		ToUser:  touser,
		MsgType: "text",
		Text: &Text{
			Content: content,
		},
	}
	msg.SetCustomService(fromuser)
	return msg
}

// newMediaMessage 新建多媒体消息
func newMediaMessage(touser, fromuser, msgtype string, media *Media) *Message {
	msg := &Message{
		ToUser:  touser,
		MsgType: msgtype,
	}
	switch msgtype {
	case "image":
		msg.Image = media
	case "voice":
		msg.Voice = media
	case "video":
		msg.Video = media
	case "mpnews":
		msg.MpNews = media
	}
	msg.SetCustomService(fromuser)
	return msg
}

// NewImageMessage 创建图片客服消息
func NewImageMessage(touser, fromuser string, media *Media) *Message {
	return newMediaMessage(touser, fromuser, "image", media)
}

// NewVoiceMessage 创建语音客服消息
func NewVoiceMessage(touser, fromuser string, media *Media) *Message {
	return newMediaMessage(touser, fromuser, "voice", media)
}

// NewVideoMessage 创建视频客服消息
func NewVideoMessage(touser, fromuser string, media *Media) *Message {
	return newMediaMessage(touser, fromuser, "video", media)
}

// NewMusicMessage 创建音乐客服消息
func NewMusicMessage(touser, fromuser string, music *Music) *Message {
	msg := &Message{
		ToUser:  touser,
		MsgType: "music",
		Music:   music,
	}
	msg.SetCustomService(fromuser)
	return msg
}

// NewNewsMessage 创建图文（外链）客服消息
func NewNewsMessage(touser, fromuser string, articles []*Article) *Message {
	msg := &Message{
		ToUser:  touser,
		MsgType: "news",
		News:    &News{articles},
	}
	msg.SetCustomService(fromuser)
	return msg
}

// NewMpNewsMessage 创建图文客服消息
func NewMpNewsMessage(touser, fromuser string, mpnews *Media) *Message {
	return newMediaMessage(touser, fromuser, "mpnews", mpnews)
}

// NewWxCardMessage 创建微信卡卷客服消息
func NewWxCardMessage(touser, fromuser string, wxcard *WxCard) *Message {
	msg := &Message{
		ToUser:  touser,
		MsgType: "wxcard",
		WxCard:  wxcard,
	}
	msg.SetCustomService(fromuser)
	return msg
}

// NewMiniProgramPageMessage 创建微信小程序页面客服消息
func NewMiniProgramPageMessage(touser, fromuser string, miniprogrampage *MiniProgramPage) *Message {
	msg := &Message{
		ToUser:          touser,
		MsgType:         "miniprogrampage",
		MiniProgramPage: miniprogrampage,
	}
	msg.SetCustomService(fromuser)
	return msg
}
