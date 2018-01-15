package weixin

import (
	"encoding/xml"
	"fmt"
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
	MsgId int64 `xml:",omitempty"`
	// MsgType text
	MsgType CDATA
	// Content 文本消息内容
	Content CDATA `xml:",omitempty"`
	// PicURL 图片链接（由系统生成）
	PicUrl CDATA `xml:",omitempty"`
	// MediaId 图片消息媒体id，可以调用多媒体文件下载接口拉取数据
	MediaId CDATA `xml:",omitempty"`
	// Foramt 语音格式，如amr，speex等
	Format CDATA `xml:",omitempty"`
	// Recognition 语音识别结果，UTF8编码
	Recognition CDATA `xml:",omitempty"`
	// ThumbMediaId 视频消息缩略图的媒体id，
	// 可以调用多媒体文件下载接口拉取数据。
	ThumbMediaId CDATA `xml:",omitempty"`
	// LocationX 地理位置维度
	Location_X float64 `xml:",omitempty"`
	// LocationY 地理位置经度
	Location_Y float64 `xml:",omitempty"`
	// Scale 地图缩放大小
	Scale int64 `xml:",omitempty"`
	// Label 地理位置信息
	Label CDATA `xml:",omitempty"`
	// Title 消息标题
	Title CDATA `xml:",omitempty"`
	// Description 消息描述
	Description CDATA `xml:",omitempty"`
	// URL 消息链接
	Url CDATA `xml:",omitempty"`
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
	MenuId string `xml:",omitempty"`
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

// ScanCodeInfo扫描事件推送
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
	Location_X CDATA `xml:",omitempty"`
	// LocationY Y坐标信息
	Location_Y CDATA `xml:",omitempty"`
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
	Image        Media     `xml:",omitempty"`
	Voice        Media     `xml:",omitempty"`
	Video        Media     `xml:",omitempty"`
	Music        Music     `xml:",omitempty"`
	ArticleCount int       `xml:",omitempty"`
	Articles     []Article `xml:",omitempty"`
}

// Media 多媒体类型的：image/voice/video，MediaId必须
type Media struct {
	MediaId     CDATA
	Titile      CDATA `xml:",omitempty"`
	Description CDATA `xml:",omitempty"`
}

// Music 回复音乐消息
type Music struct {
	Titile       CDATA `xml:",omitempty"`
	Description  CDATA `xml:",omitempty"`
	MusicURL     CDATA `xml:",omitempty"`
	HQMusicUrl   CDATA `xml:",omitempty"`
	ThumbMediaId CDATA
}

// Article 图文消息, 需要和ArticleCount一起设置
type Article struct {
	Item ArticleItem `xml:"item,omitempty"`
}

// ArticleItem 图文消息的项目
type ArticleItem struct {
	Title       CDATA
	Description CDATA
	PicUrl      CDATA
	Url         CDATA
}

func NewTextMessage(ToUserName, FromUserName, Content string) ([]byte, error) {
	msg := &ResponseMessage{
		ToUserName:   ToUserName,
		FromUserName: FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      Content,
	}
	return xml.Marshal(msg)
}

func NewImageMessage(ToUserName, FromUserName, MediaId string) ([]byte, error) {
	msg = &ResponseMessage{
		ToUserName:   CDATA(ToUserName),
		FromUserName: CDATA(FromUserName),
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		Image:        Media{MediaId: MediaId},
	}
	return xml.Marshal(msg)
}

func NewVoiceMessage(ToUserName, FromUserName, MediaId string) ([]byte, error) {
	msg = &ResponseMessage{
		ToUserName:   CDATA(ToUserName),
		FromUserName: CDATA(FromUserName),
		CreateTime:   time.Now().Unix(),
		MsgType:      "voice",
		Voice:        Media{MediaId: MediaId},
	}
	return xml.Marshal(msg)
}

func NewVideoMessage(ToUserName, FromUserName, MediaId, Title, Description string) ([]byte, error) {
	msg = &ResponseMessage{
		ToUserName:   CDATA(ToUserName),
		FromUserName: CDATA(FromUserName),
		CreateTime:   time.Now().Unix(),
		MsgType:      "video",
		Video: Media{
			MediaId:     MediaId,
			Title:       Title,
			Description: Description,
		},
	}
	return xml.Marshal(msg)
}

func NewMusicMessage(ToUserName, FromUserName, Title, Description, MusicURL, HQMusicUrl, ThumbMediaId string) ([]byte, error) {
	msg = &ResponseMessage{
		ToUserName:   CDATA(ToUserName),
		FromUserName: CDATA(FromUserName),
		CreateTime:   time.Now().Unix(),
		MsgType:      "music",
		Music: Music{
			Titile:       Title,
			Description:  Description,
			MusicURL:     MusicURL,
			HQMusicUrl:   HQMusicUrl,
			ThumbMediaId: ThumbMediaId,
		},
	}
	return xml.Marshal(msg)
}

func NewArticleMessage(ToUserName, FromUserName, Articles []Article) ([]byte, error) {
	msg = &ResponseMessage{
		ToUserName:   CDATA(ToUserName),
		FromUserName: CDATA(FromUserName),
		CreateTime:   time.Now().Unix(),
		MsgType:      "music",
		ArticleCount: len(Articles),
		Articles:     Articles,
	}
	return xml.marshal(msg)
}

func newExampleMsg(from, to, content string) (string, error) {
	b, err := NewTextMessage(to, from, content)
	if err != nil {
		return "", fmt.Errorf("marshal the message to xml: %s", err)
	}
	return string(b), nil
}
