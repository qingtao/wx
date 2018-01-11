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
	PicURL CDATA `xml:"PicUrl,omitempty"`
	// MediaId 图片消息媒体id，可以调用多媒体文件下载接口拉取数据
	MediaID CDATA `xml:"Mediaid,omitempty"`
	// Foramt 语音格式，如amr，speex等
	Format CDATA `xml:",omitempty"`
	// Recognition 语音识别结果，UTF8编码
	Recognition CDATA `xml:",omitempty"`
	// ThumbMediaId 视频消息缩略图的媒体id，
	// 可以调用多媒体文件下载接口拉取数据。
	ThumbMediaId CDATA `xml:",omitempty"`
	// LocationX 地理位置维度
	LocationX float64 `xml:"Location_X,omitempty"`
	// LocationY 地理位置经度
	LocationY float64 `xml:"Location_Y,omitempty"`
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
	MenuID string `xml:,omitempty"`
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
	ScanType CDATA `xml:"ScanType,omitempty"`
	// ScanResult
	ScanResult CDATA `xml:"ScanResult,omitempty"`
}

// Item in PicList
type Item struct {
	// PicMd5Sum 图片的MD5值，开发者若需要，可用于验证接收到图片
	PicMd5Sum CDATA `xml:omitempty"`
}

// PicList 图片列表
type PicList struct {
	// PicMd5Sum 图片的MD5值，开发者若需要，可用于验证接收到图片
	Item *Item `xml:"omitempty"`
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

func newExampleMsg(from, to string) (string, error) {
	msg := &Message{
		ToUserName:   CDATA(to),
		FromUserName: CDATA(from),
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "欢迎关注！",
	}
	b, err := xml.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshal the message to xml: %s", err)
	}
	return string(b), nil
}
