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
	"os"
	"path/filepath"
	"strings"
)

const (
	// WxImageJPEG 微信的图片格式: jpeg
	WxImageJPEG = "jpeg"
	// WxImageJPG 微信的图片格式: jpeg
	WxImageJPG = "jpg"
	// WxImagePNG 微信的图片格式: png
	WxImagePNG = "png"
	// WxImageGIF 微信的图片格式: gif
	WxImageGIF = "gif"
	// WxVoiceMP3 微信的音频：mp3
	WxVoiceMP3 = "mp3"
	// WxVoiceAMR 微信的音频：amr
	WxVoiceAMR = "amr"
	// WxMediaUpload 微信媒体素材上传的路径
	WxMediaUpload = "cgi-bin/media/upload"
	// WxMediaGet 获取微信媒体素材的路径
	WxMediaGet = "cgi-bin/media/get"
)

var (
	// WxImageMaxSize 图片文件最大1MB
	WxImageMaxSize = 1024 * 1024
	// WxVoiceMaxSize 音频文件最大1MB
	WxVoiceMaxSize = 1024 * 1024
	// WxVideoMaxSize 视频最大10MB
	WxVideoMaxSize = 10 * 1024 * 1024
	// WxThumbMaxSize 缩略图最大64KB
	WxThumbMaxSize = 64 * 1024
)

// Response 微信公众平台媒体资源API响应结构，放置在一起方便解析响应的消息和错误
type Response struct {
	Type      string `json:"type,omitempty"`
	MediaId   string `json:"media_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
	ErrMsg    string `json:"errmsg,omitempty"`
}

// VerifyImageFormat 素材格式是否正确
func VerifyImageFormat(filename, imageType string) bool {
	ext := filepath.Ext(filename)
	// 扩展名转化为小写并与typ比较，一致返回true
	if strings.ToLower(ext) == imageType {
		return true
	}
	return false
}

// uploadMedia 上传素材到公众平台，host 正常是通过微信公众平台的域名，accessToken 是调用接口凭证
func uploadMedia(host, typ, filename, accessToken string) (*Response, error) {
	ct, r, err := ParseFile("media", filename)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("https://%s/%s?access_token=%s&type=%s",
		host, WxMediaUpload, accessToken, typ)
	res, err := http.Post(uri, ct, r)
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

func UploadImage(host, filename, access_token string) (*Response, error) {
	return uploadMedia(host, "image", filename, access_token)
}

func UploadVoice(host, filename, access_token string) (*Response, error) {
	return uploadMedia(host, "voice", filename, access_token)
}

func UploadVideo(host, filename, access_token string) (*Response, error) {
	return uploadMedia(filename, "video", access_token)
}

func Upload(host, filename, access_token string) (*Response, error) {
	return uploadMedia(filename, "thumb", access_token)
}

// ParseFile 解析文件，验证文件的大小，返回 contentType io.Reader err, 异常时返回 err 非空
func ParseFile(typ, filename string) (contentType string, r io.Reader, err error) {
	maxsize := WxImageMaxSize
	switch typ {
	case "image":
	case "voice":
		maxsize = WxVoiceMaxSize
	case "video":
		maxsize = WxVideoMaxSize
	case "thumb":
		maxsize = WxThumbMaxSize
	}
	stat, err := os.Stat(filename)
	if err != nil {
		return
	}
	if stat.Size() > int64(maxsize) {
		return "", nil, fmt.Errorf("%s file too large than %d", typ, maxsize)
	}
	var buf bytes.Buffer
	multiWriter := multipart.NewWriter(&buf)
	defer multiWriter.Close()
	w, err := multiWriter.CreateFormFile("media", stat.Name())
	if err != nil {
		return
	}
	fr, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fr.Close()
	io.Copy(w, fr)
	contentType = multiWriter.FormDataContentType()
	return
}

// ErrResponse 下载素材时返回错误信息用
type DownloadResponse struct {
	VideoURL string `json:"video_url"`
	ErrCode  string `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
}

func (dr DownloadResponse) String() string {
	return "errcode: " + dr.ErrCode + "errmsg: " + dr.ErrMsg
}

// GetMedia 下载素材
func GetMedia(host, mediaID, accessToken, dir string) (string, error) {
	uri := fmt.Sprintf("http://%s?access_token=%smedia_id=%s",
		host, accessToken, mediaID)
	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	ct := res.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", err
	}
	switch {
	case strings.HasSuffix(mediaType, "json"):
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("read body %s", err)
		}
		defer res.Body.Close()
		var downResponse DownloadResponse
		if err = json.Unmarshal(b, &downResponse); err != nil {
			return "", fmt.Errorf("read body ok, %s", err)
		}
		if downResponse.VideoURL == "" {
			return "", fmt.Errorf("%s", downResponse)
		}
		return downResponse.VideoURL, nil
	case strings.HasSuffix(mediaType, WxImageGIF):
		fallthrough
	case strings.HasSuffix(mediaType, WxImageJPEG):
		fallthrough
	case strings.HasSuffix(mediaType, WxImagePNG):
		fallthrough
	case strings.HasSuffix(mediaType, WxImageJPG):
		
		return "", nil
	}
	return "", nil
}
