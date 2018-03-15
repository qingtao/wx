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
	WxImageJPEG = "jpeg"
	// WxImageJPG 微信的图片格式: jpg
	WxImageJPG = "jpg"
	// WxImagePNG 微信的图片格式: png
	WxImagePNG = "png"
	// WxImageGIF 微信的图片格式: gif
	WxImageGIF = "gif"
	// WxVoiceMP3 微信的音频：mp3
	WxVoiceMP3 = "mp3"
	// WxVoiceAMR 微信的音频：amr
	WxVoiceAMR = "amr"
	// WxVideoMP4 微信的视频：mp4
	WxVideoMP4 = "mp4"
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

// Response 微信公众平台媒体资源API响应结构，放置在一起方便解析响应的消息和错误
type Response struct {
	Type      string `json:"type,omitempty"`
	MediaID   string `json:"media_id,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
	ErrMsg    string `json:"errmsg,omitempty"`
}

// uploadMedia 上传素材到公众平台，host 正常是通过微信公众平台的域名，accessToken 是调用接口凭证
func uploadMedia(host, typ, filename, accessToken string) (*Response, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	maxsize := WxImageMaxSize
	switch typ {
	case "image":
		switch ext {
		case WxImageGIF, WxImageJPEG, WxImageJPG, WxImagePNG:
			return nil, fmt.Errorf("image must be one of jpg|jpeg|gif|png")
		}
	case "voice":
		if ext != WxVoiceAMR && ext != WxVoiceMP3 {
			return nil, fmt.Errorf("voice must be one of amr|mp3")
		}
		maxsize = WxVoiceMaxSize
	case "video":
		if ext != WxVideoMP4 {
			return nil, fmt.Errorf("video must be mp4")
		}
		maxsize = WxVideoMaxSize
	case "thumb":
		if ext != WxImageJPG {
			return nil, fmt.Errorf("thumb must be jpg")
		}
		maxsize = WxThumbMaxSize
	default:
		return nil, fmt.Errorf("media type not supported")
	}
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("upload %s %s", typ, err)
	}
	if stat.Size() > int64(maxsize) {
		return nil, fmt.Errorf("%s file too large than %d", typ, maxsize)
	}
	//fmt.Println(stat.Size())
	if stat.Size() <= 0 {
		return nil, fmt.Errorf("%s file size is zero", typ)
	}
	var buf = new(bytes.Buffer)
	multiWriter := multipart.NewWriter(buf)
	w, err := multiWriter.CreateFormFile("media", stat.Name())
	if err != nil {
		return nil, err
	}
	fr, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	if _, err = io.Copy(w, fr); err != nil {
		return nil, err
	}
	multiWriter.Close()
	contentType := multiWriter.FormDataContentType()

	fmt.Println(contentType)
	uri := fmt.Sprintf("https://%s/%s?access_token=%s&type=%s",
		host, WxMediaUpload, accessToken, typ)
	res, err := http.Post(uri, contentType, buf)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadImage 上传图片
func UploadImage(host, filename, accessToken string) (*Response, error) {
	return uploadMedia(host, "image", filename, accessToken)
}

// UploadVoice 上传音频
func UploadVoice(host, filename, accessToken string) (*Response, error) {
	return uploadMedia(host, "voice", filename, accessToken)
}

// UploadVideo 上传视频
func UploadVideo(host, filename, accessToken string) (*Response, error) {
	return uploadMedia(filename, "video", filename, accessToken)
}

// UploadThumb 上传缩略图
func UploadThumb(host, filename, accessToken string) (*Response, error) {
	return uploadMedia(filename, "thumb", filename, accessToken)
}

// DownloadResponse 下载素材时返回错误信息用
type DownloadResponse struct {
	// VideoURL 如果是请求的视频，服务器响应会设置此字段
	VideoURL string `json:"video_url"`
	// ErrCode 错误代码
	ErrCode int `json:"errcode"`
	// ErrMsg  错误消息
	ErrMsg string `json:"errmsg"`
}

// String 只返回错误代码和错误消息
func (dr DownloadResponse) String() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", dr.ErrCode, dr.ErrMsg)
}

// GetMedia 下载素材, 如果error为nil，返回的字符串是文件保存的绝对路径
func GetMedia(host, mediaID, accessToken, dir string) (string, error) {
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
		fmt.Printf("%s\n", b)
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
	// 写入文件file
	if err = ioutil.WriteFile(file, b, 0640); err != nil {
		return "", fmt.Errorf("write file %s", err)
	}
	return file, nil
}

// Material 媒体永久图文素材
type Material struct {
	Articles *Article `json:"articles"`
}

// Article 图文，永久的
type Article struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	ShowCoverPic     int    `json:"show_cover_pic"`
	Content          string `json:"content "`
	ContentSourceURL string `json:"content_source_url"`
}

// WxMaterialImageMaxSize 图文消息内的图片，只支持jpg/png，且大小不能大于1MB
const (
	// WxMaterialImageMaxSize 上传图文中的图片最大值
	WxMaterialImageMaxSize = 1024 * 1024
	// WxMediaMaterialAdd 永久图文素材上传路径,此处上传的图片不受素材库数量的限制
	WxMediaMaterialAdd = "cgi-bin/media/material/add_news"
	// WxMediaMaterialGet 获取永久图文素材路径
	WxMediaMaterialGet = "cgi-bin/media/material/get_material"
	// WxMediaMateriaImagelAdd 图文中图片上传的路径
	WxMediaMateriaImagelAdd = "cgi-bin/media/uploadimg"
	// WxMediaMateriaImagelDel 删除永久图文的路径
	WxMediaMateriaImagelDel = "cgi-bin/media/material/del_material"
)
