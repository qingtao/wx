package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	// WxImageJPEG 微信的图片格式: jpeg
	WxImageJPEG = ".jpeg"
	// WxImageJPG 微信的图片格式: jpg
	WxImageJPG = ".jpg"
	// WxImagePNG 微信的图片格式: png
	WxImagePNG = ".png"
	// WxImageGIF 微信的图片格式: gif
	WxImageGIF = ".gif"
	// WxVoiceMP3 微信的音频：mp3
	WxVoiceMP3 = ".mp3"
	// WxVoiceAMR 微信的音频：amr
	WxVoiceAMR = ".amr"
	// WxVideoMP4 微信的视频：mp4
	WxVideoMP4 = ".mp4"
	// WxMediaUpload 微信临时媒体素材上传的路径
	WxMediaUpload = "cgi-bin/media/upload"
	// WxMediaGet 获取微信临时媒体素材的路径
	WxMediaGet = "cgi-bin/media/get"
)

var (
	// WxImageMaxSize 图片文件最大1MB
	WxImageMaxSize = 2 * 1024 * 1024
	// WxVoiceMaxSize 音频文件最大1MB
	WxVoiceMaxSize = 2 * 1024 * 1024
	// WxVideoMaxSize 视频最大10MB
	WxVideoMaxSize = 10 * 1024 * 1024
	// WxThumbMaxSize 缩略图最大64KB
	WxThumbMaxSize = 64 * 1024
)

// UploadResponse 微信公众平台媒体资源API响应结构，放置在一起方便解析响应的消息和错误
type UploadResponse struct {
	Type      string `json:"type,omitempty"`
	MediaID   string `json:"media_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
	ErrMsg    string `json:"errmsg,omitempty"`
}

// ParseFile 读取文件并检查文件的大小、扩展名，返回mutltipart的Content-Type，io.Reader, 如果任何错误，则err非空
func ParseFile(typ, filename string, maxsize int, desc []byte) (contentType string, r io.Reader, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch typ {
	case "image":
		// 检查图片扩展名
		switch ext {
		case WxImageGIF, WxImageJPEG, WxImageJPG, WxImagePNG:
		default:
			return "", nil, fmt.Errorf("image must be one of jpg|jpeg|gif|png")
		}
		if maxsize == 0 {
			maxsize = WxImageMaxSize
		}
	case "voice":
		if ext != WxVoiceAMR && ext != WxVoiceMP3 {
			return "", nil, fmt.Errorf("voice must be one of amr|mp3")
		}
		maxsize = WxVoiceMaxSize
	case "video":
		if ext != WxVideoMP4 {
			return "", nil, fmt.Errorf("video must be mp4")
		}
		maxsize = WxVideoMaxSize
	case "thumb":
		if ext != WxImageJPG {
			return "", nil, fmt.Errorf("thumb must be jpg")
		}
		maxsize = WxThumbMaxSize
		// 不符合扩展名的，返回错误
	default:
		return "", nil, fmt.Errorf("media type not supported")
	}
	// 取文件状态，得到文件名称和大小
	stat, err := os.Stat(filename)
	if err != nil {
		return "", nil, fmt.Errorf("upload %s %s", typ, err)
	}
	// 大于最大值，返回提示错误
	if stat.Size() > int64(maxsize) {
		return "", nil, fmt.Errorf("%s file too large than %d", typ, maxsize)
	}
	// 保证文件大小大于0
	if stat.Size() <= 0 {
		return "", nil, fmt.Errorf("size of file %s is zero", filename)
	}
	var buf = new(bytes.Buffer)
	multiWriter := multipart.NewWriter(buf)
	// multipart的文件，名称“media”是微信要求的参数
	w, err := multiWriter.CreateFormFile("media", stat.Name())
	if err != nil {
		return "", nil, err
	}
	// 打开文件
	fr, err := os.Open(filename)
	if err != nil {
		return "", nil, err
	}
	defer fr.Close()
	// 文件内容写入到w->buf
	if _, err = io.Copy(w, fr); err != nil {
		return "", nil, err
	}
	//写入 MaterialVideo 的 description
	if desc != nil {
		w, err = multiWriter.CreateFormField("description")
		if err != nil {
			return "", nil, err
		}
		if _, err = bytes.NewBuffer(desc).WriteTo(w); err != nil {
			return "", nil, err
		}
	}
	multiWriter.Close()
	// 获取http头部的Content-Type
	r = buf
	contentType = multiWriter.FormDataContentType()
	return

}

// uploadMedia 上传素材到公众平台，host 正常是通过微信公众平台的域名，accessToken 是调用接口凭证
func uploadMedia(host, accessToken, typ, filename string) (*UploadResponse, error) {
	contentType, r, err := ParseFile(typ, filename, 0, nil)
	if err != nil {
		return nil, err
	}
	// 使用host,WxMediaUpload,accessToken和typ连接成url
	uri := fmt.Sprintf("https://%s/%s?access_token=%s&type=%s",
		host, WxMediaUpload, accessToken, typ)
	res, err := http.Post(uri, contentType, r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("when post %s %s: %s", typ, filename, err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("when post file %s, message received: %s", filename, err)
	}
	defer res.Body.Close()
	var resp UploadResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadImage 上传图片
func UploadImage(host, accessToken, filename string) (*UploadResponse, error) {
	return uploadMedia(host, "image", filename, accessToken)
}

// UploadVoice 上传音频
func UploadVoice(host, accessToken, filename string) (*UploadResponse, error) {
	return uploadMedia(host, "voice", filename, accessToken)
}

// UploadVideo 上传视频
func UploadVideo(host, accessToken, filename string) (*UploadResponse, error) {
	return uploadMedia(host, "video", filename, accessToken)
}

// UploadThumb 上传缩略图
func UploadThumb(host, accessToken, filename string) (*UploadResponse, error) {
	return uploadMedia(host, "thumb", filename, accessToken)
}

// DownloadResponse 下载素材时返回错误信息用
type DownloadResponse struct {
	// VideoURL 如果是请求的视频，服务器响应会设置此字段
	VideoURL string `json:"video_url,omitempty"`
	// ErrCode 错误代码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg  错误消息
	ErrMsg string `json:"errmsg,omitempty"`
}

// String 只返回错误代码和错误消息
func (dr DownloadResponse) String() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", dr.ErrCode, dr.ErrMsg)
}

// GetMedia 下载素材, 如果error为nil，返回的字符串是文件保存的绝对路径
func GetMedia(host, accessToken, mediaID, dir string) (string, error) {
	uri := fmt.Sprintf("http://%s/%s?access_token=%smedia_id=%s",
		host, WxMediaUpload, accessToken, mediaID)
	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("%s", res.Status)
	}
	mediaType := res.Header.Get("Content-Type")
	filename := ""
	switch {
	// 如果响应中的Content-Type 是jpg|png|gif， 提取文件名
	case strings.Contains(mediaType, WxImageGIF) ||
		strings.Contains(mediaType, WxImageJPEG) ||
		strings.Contains(mediaType, WxImagePNG):
		_, params, err := mime.ParseMediaType(
			res.Header.Get("Content-disposition"))
		if err != nil {
			return "", fmt.Errorf("get filename %s", err)
		}
		// 赋值filename
		filename = params["filename"]
		// 如果响应的Content-Type是json或者text，尝试读取body并解析json
	case strings.Contains(mediaType, "json") ||
		strings.Contains(mediaType, "text"):
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("read body %s", err)
		}
		//fmt.Printf("%s\n", b)
		defer res.Body.Close()
		var downResponse DownloadResponse
		if err = json.Unmarshal(b, &downResponse); err != nil {
			return "", fmt.Errorf("read body ok, %s", err)
		}
		if downResponse.VideoURL == "" {
			return "", fmt.Errorf("%s", downResponse)
		}
		uri, err := url.Parse(downResponse.VideoURL)
		if err != nil {
			return "", fmt.Errorf(`video's url is invalid %s`, err)
		}
		//赋值filename
		filename = uri.Path

		res, err = http.Get(downResponse.VideoURL)
		if err != nil {
			return "", fmt.Errorf("get %s %s", downResponse.VideoURL, err)
		}
	}
	if filename == "" {
		return "", fmt.Errorf("get media: find filename error")
	}
	file := filepath.Join(dir, filename)
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response body %s", err)
	}
	defer res.Body.Close()
	// 写入文件file
	if err = ioutil.WriteFile(file, b, 0640); err != nil {
		return "", fmt.Errorf("write file %s", err)
	}
	return file, nil
}

// MaterialArticle 媒体永久图文素材
type MaterialArticle struct {
	Articles []*Article `json:"articles"`
}

// Parse 实现 MaterialMedia接口
func (m *MaterialArticle) Parse() (string, io.Reader, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", nil, err
	}
	return "application/json; charset=utf-8", bytes.NewReader(b), nil
}

// Upload 使用参数host, access_token，上传图文素材*MaterialArticle到公众平台永久素材库
func (m *MaterialArticle) Upload(host, accessToken string) (*MaterialResponse, error) {
	return UploadMaterial(host, accessToken, WxMaterailAdd, "", m)
}

// Article 图文，永久的
type Article struct {
	// Title 标题
	Title string `json:"title"`
	// ThumbMediaID 图文消息的封面素材ID
	ThumbMediaID string `json:"thumb_media_id"`
	// Author 作者
	Author string `json:"author"`
	// Digest 图文消息的摘要，仅有单图消息才有摘要，多图文此处为空; 若该值为空，默认选取前64个字
	Digest string `json:"digest"`
	// ShowCoverPic 是否显示封面，0：false，1：true
	ShowCoverPic int `json:"show_cover_pic"`
	// Content 图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS,涉及图片url必须来源 "上传图文消息内的图片获取URL"接口获取。外部图片url将被过滤。
	Content string `json:"content"`
	// ContentSourceURL 图文消息的原文地址，即点击“阅读原文”后的URL
	ContentSourceURL string `json:"content_source_url"`
	// NeedOpenComment 可以留言或者评论
	NeedOpenComment uint32 `json:"need_open_comment,omitempty"`
	// OnlyFansCanComment 只有公众号粉丝评价
	OnlyFansCanComment uint32 `json:"only_fans_can_comment,omitempty"`
}

// WxMaterialImageMaxSize 图文消息内的图片，只支持jpg/png，且大小不能大于1MB
const (
	// WxMaterialImageMaxSize 上传图文中的图片最大值
	WxMaterialImageMaxSize = 1024 * 1024
	// WxMaterailAdd 永久图文素材上传路径,此处上传的图片不受素材库数量的限制
	WxMaterailAdd = "cgi-bin/material/add_news"
	// WxMaterailAddOther 除图文和图文中图片以外的其他类型上传路径
	WxMaterailAddOther = "cgi-bin/material/add_material"
	// WxMediaUploadImg 图文中图片上传的路径
	WxMediaUploadImg = "cgi-bin/media/uploadimg"
)

// MaterialResponse 上传图文素材的响应结构
type MaterialResponse struct {
	// Title 标题, 视频素材的返回值中使用
	Title string `json:"title,omitempty"`
	// MediaID, 素材的mediaId, 视频和图文素材中使用
	MediaID string `json:"media_id,omitempty"`
	// URL 素材链接
	URL string `json:"url,omitempty"`
	// ErrCode 错误代码
	ErrCode int `json:"errcode,omitempty"`
	// ErrMsg 错误消息
	ErrMsg string `json:"errmsg,omitempty"`
}

// MaterialMedia 永久的多媒体素材
type MaterialMedia interface {
	// Parse 返回 Content-Type, io.Reader，error
	Parse() (string, io.Reader, error)
}

// UploadMaterial 上传图文素材到公众平台
func UploadMaterial(host, accessToken, path, typ string, m MaterialMedia) (*MaterialResponse, error) {
	//使用接口MaterialMedia可以简化视频、图片和音频等的操作
	contentType, r, err := m.Parse()
	if err != nil {
		return nil, err
	}
	//uri是图文素材上传的URL地址
	uri := fmt.Sprintf("https://%s/%s?access_token=%s", host, path, accessToken)
	// video等类型的文件需要type查询参数
	if typ != "" {
		uri = uri + "&type=" + typ
	}

	res, err := http.Post(uri, contentType, r)
	if err != nil {
		return nil, err
	}
	// TODO: 可能需要检查其他前期的操作判断http相应码
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("when post material, response status is %d", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("when post material, read response failed: %s", err)
	}
	defer res.Body.Close()

	var materialResponse MaterialResponse
	if err = json.Unmarshal(b, &materialResponse); err != nil {
		return nil, fmt.Errorf("when post material, read json message %s", err)
	}
	return &materialResponse, nil
}

// MaterialImage 图片素材，永久的,如果 InMaterial为true，不添加type查询字符串
type MaterialImage struct {
	InMaterial bool
	FileName   string
}

// Parse 实现 MaterialMedia接口
func (m *MaterialImage) Parse() (string, io.Reader, error) {
	size := WxImageMaxSize
	// 判断图片大小限制
	if m.InMaterial {
		size = WxMaterialImageMaxSize
	}
	return ParseFile("image", m.FileName, size, nil)
}

// Upload 上传图片素材
func (m *MaterialImage) Upload(host, accessToken string) (*MaterialResponse, error) {
	typ, path := "image", WxMaterailAddOther
	if m.InMaterial {
		typ = ""
		path = WxMediaUploadImg
	}
	return UploadMaterial(host, accessToken, path, typ, m)
}

// MaterialVideo 视频素材，永久的
type MaterialVideo struct {
	FileName     string `json:"-"`
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
}

// Describe 返回 video 要求的 description
func (m *MaterialVideo) Describe() ([]byte, error) {
	return json.Marshal(m)
}

// Parse 实现 MaterialMedia 接口
func (m *MaterialVideo) Parse() (string, io.Reader, error) {
	desc, err := m.Describe()
	if err != nil {
		return "", nil, err
	}
	return ParseFile("video", m.FileName, WxVideoMaxSize, desc)
}

// Upload 上传图片文件到微信公共平台
func (m *MaterialVideo) Upload(host, accessToken string) (*MaterialResponse, error) {
	return UploadMaterial(host, accessToken, WxMaterailAddOther, "video", m)
}

// MaterialVoice 音频素材，永久的
type MaterialVoice struct {
	FileName string
}

// Parse 实现 MaterialMedia 接口
func (m *MaterialVoice) Parse() (string, io.Reader, error) {
	return ParseFile("voice", m.FileName, WxVoiceMaxSize, nil)
}

// Upload 上传音频文件到微信公共平台
func (m *MaterialVoice) Upload(host, accessToken string) (*MaterialResponse, error) {
	return UploadMaterial(host, accessToken, WxMaterailAddOther, "voice", m)
}

// Materialthumb 素材的缩略图
type Materialthumb struct {
	FileName string
}

// Parse 实现 MaterialMedia 接口
func (m *Materialthumb) Parse() (string, io.Reader, error) {
	return ParseFile("thumb", m.FileName, WxThumbMaxSize, nil)
}

// Upload 上传缩略图文件到微信公共平台
func (m *Materialthumb) Upload(host, accessToken string) (*MaterialResponse, error) {
	return UploadMaterial(host, accessToken, WxMaterailAddOther, "thumb", m)
}
