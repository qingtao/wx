package mp

import (
	"encoding/xml"
	"time"
)

// Message 微信公众平台消息结构
type Message struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	// ToUserName 开发者微信号
	ToUserName CDATA
	// FromUserName 发送方帐号（一个OpenID）
	FromUserName CDATA
	// CreateTime 消息创建时间 （整型）
	CreateTime int64
	// MsgId 消息id，64位整型
	MsgID int64 `xml:"MsgId,omitempty" json:"MsgId,omitempty"`
	// MsgType text
	MsgType CDATA
	// Content 文本消息内容
	Content CDATA `xml:",omitempty"`
	// PicURL 图片链接（由系统生成）
	PicURL CDATA `xml:"PicUrl,omitempty" json:"PicUrl,omitempty"`
	// MediaId 图片消息媒体id，可以调用多媒体文件下载接口拉取数据
	MediaID CDATA `xml:"MediaId,omitempty" json:"MediaId,omitempty"`
	// Foramt 语音格式，如amr，speex等
	Format CDATA `xml:",omitempty"`
	// Recognition 语音识别结果，UTF8编码
	Recognition CDATA `xml:",omitempty"`
	// ThumbMediaId 视频消息缩略图的媒体id，
	// 可以调用多媒体文件下载接口拉取数据。
	ThumbMediaID CDATA `xml:"ThumbMediaId,omitempty" json:"ThumbMediaId,omitempty"`
	// LocationX 地理位置维度
	LocationX float64 `xml:"Location_X,omitempty" json:"Location_X,omitempty"`
	// LocationY 地理位置经度
	LocationY float64 `xml:"Location_Y,omitempty" json:"Location_Y,omitempty"`
	// Scale 地图缩放大小
	Scale int64 `xml:",omitempty"`
	// Label 地理位置信息
	Label CDATA `xml:",omitempty"`
	// Title 消息标题
	Title CDATA `xml:",omitempty"`
	// Description 消息描述
	Description CDATA `xml:",omitempty"`
	// URL 消息链接
	URL CDATA `xml:"Url,omitempty"`
	// Event 事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Event CDATA `xml:",omitempty"`
	// EventKey 事件KEY值，qrscene_为前缀，后面为二维码的参数值
	EventKey CDATA `xml:",omitempty"`
	// Ticket 二维码的ticket，可用来换取二维码图片
	Ticket CDATA `xml:",omitempty"`
	// Latitude 地理位置纬度
	Latitude float64 `xml:",omitempty"`
	// Longitude 地理位置经度
	Longitude float64 `xml:",omitempty"`
	// Precision 地理位置精度
	Precision float64 `xml:",omitempty"`
	// Scene 场景值，固定为1
	Scene string `xml:",omitempty"`

	// 自定义菜单事件

	// MenuID 点击菜单跳转链接时的时间推送, Event: VIEW
	MenuID string `xml:"MenuId,omitempty" json:"MenuId,omitempty"`
	// ScanCodeInfo: 扫描事件推送
	//   1. Event是scancode_push：
	//	   扫码推事件的事件推送
	//   2. Event是scancode_waitmsg：
	//     扫码推事件且弹出“消息接收中”提示框的事件推送
	ScanCodeInfo *ScanCodeInfo `xml:",omitempty"`
	// SendPiscInfo 拍照发图的事件推送
	//   1. Event是pic_sysphoto:
	//     弹出系统拍照发图的事件推送
	//   2. Event是pic_photo_or_album:
	//     弹出拍照或者相册发图的事件推送
	//   3. Event是pic_weixin:
	//     弹出微信相册发图器的事件推送
	SendPicsInfo *SendPicsInfo `xml:",omitempty"`
	// 	SendLocationInfo 发送的位置信息
	SendLocationInfo *SendLocationInfo `xml:",omitempty"`
}

// CDATA xml <![CDATA[...]]]格式
type CDATA string

// MarshalXML 实现xml.Marshaler接口
func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var v = struct {
		string `xml:",cdata"`
	}{string(c)}
	return e.EncodeElement(v, start)
}

// ScanCodeInfo 扫描事件推送
type ScanCodeInfo struct {
	// ScanType
	ScanType CDATA `xml:",omitempty"`
	// ScanResult
	ScanResult CDATA `xml:",omitempty"`
}

// Item in PicList
type Item struct {
	// PicMd5Sum 图片的MD5值，开发者若需要，可用于验证接收到图片
	PicMd5Sum CDATA `xml:",omitempty"`
}

// PicList 图片列表
type PicList struct {
	// PicMd5Sum 图片的MD5值，开发者若需要，可用于验证接收到图片
	Item *Item `xml:",omitempty"`
}

// SendPicsInfo 发送图片事件
type SendPicsInfo struct {
	// 发送的图片数量
	Count int `xml:",omitempty"`
	// PicList 图片列表
	PicList []*PicList `xml:",omitempty"`
}

// SendLocationInfo 弹出地理位置选择器的事件推送
type SendLocationInfo struct {
	// LocationX X坐标信息
	LocationX CDATA `xml:"Location_X,omitempty"`
	// LocationY Y坐标信息
	LocationY CDATA `xml:"Location_Y,omitempty"`
	// Scale 精度，可理解为精度或者比例尺、越精细的话 scale越高
	Scale CDATA `xml:",omitempty"`
	// Lable 地理位置的字符串信息
	Label CDATA `xml:",omitempty"`
	// Poiname 朋友圈POI的名字，可能为空
	Poiname CDATA `xml:",omitempty"`
}

// ResponseMessage 被动回复微信消息事件的消息结构
type ResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   int64
	MsgType      CDATA
	Content      CDATA     `xml:",omitempty"`
	Image        *Media    `xml:",omitempty"`
	Voice        *Media    `xml:",omitempty"`
	Video        *Media    `xml:",omitempty"`
	Music        *Music    `xml:",omitempty"`
	ArticleCount int       `xml:",omitempty"`
	Articles     *Articles `xml:",omitempty"`
}

// Media 多媒体类型的：image/voice/video，MediaId必须
type Media struct {
	MediaID     CDATA `xml:"MediaId,omitempty"`
	Title       CDATA `xml:",omitempty"`
	Description CDATA `xml:",omitempty"`
}

// Music 回复音乐消息
type Music struct {
	Title        CDATA `xml:",omitempty"`
	Description  CDATA `xml:",omitempty"`
	MusicURL     CDATA `xml:",omitempty"`
	HQMusicURL   CDATA `xml:"HQMusicUrl,omitempty"`
	ThumbMediaID CDATA `xml:"ThumbMediaId,omitempty"`
}

// Articles 图文消息
type Articles struct {
	Item []*Article `xml:"item,omitempty"`
}

// Article 图文消息
type Article struct {
	Title       CDATA `xml:",omitempty"`
	Description CDATA `xml:",omitempty"`
	PicURL      CDATA `xml:"PicUrl,omitempty"`
	URL         CDATA `xml:"Url,omitempty"`
}

// NewArticle 新建图文文章
func NewArticle(Title, Description, PicURL, URL string) *Article {
	return &Article{CDATA(Title), CDATA(Description), CDATA(PicURL), CDATA(URL)}
}

// NewTextMessage 创建被动回复文本消息
func NewTextMessage(ToUserName, FromUserName CDATA, Content string) *ResponseMessage {
	return &ResponseMessage{
		ToUserName:   ToUserName,
		FromUserName: FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      CDATA(Content),
	}
}

// NewMedia 新的多媒体消息，包含：image/voice/video
func NewMedia(MediaID, Title, Description string) *Media {
	return &Media{CDATA(MediaID), CDATA(Title), CDATA(Description)}
}

// NewMusic 新音乐消息结构
func NewMusic(Title, Description, MusicURL, HQMusicURL, ThumbMediaID string) *Music {
	return &Music{CDATA(Title), CDATA(Description), CDATA(MusicURL), CDATA(HQMusicURL), CDATA(ThumbMediaID)}
}

// newMediaMessage 新的多媒体消息
func newMediaMessage(ToUserName, FromUserName CDATA, MsgType string, media *Media) *ResponseMessage {
	msg := &ResponseMessage{
		ToUserName:   ToUserName,
		FromUserName: FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      CDATA(MsgType),
	}
	switch MsgType {
	case "image":
		msg.Image = media
	case "voice":
		msg.Voice = media
	case "video":
		msg.Video = media
	}
	return msg
}

// NewImageMessage 创建被动回复图片消息
func NewImageMessage(ToUserName, FromUserName CDATA, media *Media) *ResponseMessage {
	return newMediaMessage(ToUserName, FromUserName, "image", media)
}

// NewVoiceMessage 创建被动回复语音消息
func NewVoiceMessage(ToUserName, FromUserName CDATA, media *Media) *ResponseMessage {
	return newMediaMessage(ToUserName, FromUserName, "voice", media)
}

// NewVideoMessage 创建被动回复视频消息, meida需要MediaId, title, description
func NewVideoMessage(ToUserName, FromUserName CDATA, media *Media) *ResponseMessage {
	return newMediaMessage(ToUserName, FromUserName, "video", media)
}

// NewMusicMessage 创建被动回复音乐消息
func NewMusicMessage(ToUserName, FromUserName CDATA, music *Music) *ResponseMessage {
	return &ResponseMessage{
		ToUserName:   ToUserName,
		FromUserName: FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "music",
		Music:        music,
	}
}

// NewArticlesMessage 创建被动回复图文消息
func NewArticlesMessage(ToUserName, FromUserName CDATA, articles []*Article) *ResponseMessage {
	return &ResponseMessage{
		ToUserName:   ToUserName,
		FromUserName: FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "news",
		ArticleCount: len(articles),
		Articles:     &Articles{articles},
	}
}
