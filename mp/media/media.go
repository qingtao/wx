package media

const (
	WxImage       = "0xFFD8FF"
	WxMediaUpload = "cgi-bin/media/upload"
	WxMediaGet    = "cgi-bin/media/get"
)

var (
	WxImageMaxSize = 1024 * 1024
	WxVoiceMaxSize = 1024 * 1024
	WxVideoMaxSize = 10 * 1024 * 1024
	WxThumbMaxSize = 64 * 1024
)

type Response struct {
	Type      string `json:"type,omitempty"`
	MediaId   string `json:"media_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
	ErrMsg    string `json:"errmsg,omitempty"`
}

func VerifyImageFormat(filename string) bool {
	b, _ := hex.DecodeString(WxImage)
	fileheader := make([]byte, 8)
	r, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer r.Close()
	r.Read(fileheader)
	if bytes.HasPrefix(fileheader, b) {
		return true
	}
	return false
}

func uploadMedia(host, typ, filename, access_token) (*Response, error) {
	/*
		if !VerifyImageFormat(filename) {
			return nil, fmt.Errorf("image must be jpg")
		}
	*/
	ct, r, err := ParseFile("media", filename)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("https://%s/%s?access_token=%s&type=%s",
		host, WxMediaUpload, access_token, typ)
	res, err := http.Post(uri, ct, r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var status Response
	if err := json.Unmarshal(b, &status); err != nil {
		return nil, err
	}
	return &status, nil
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

func ParseFile(typ, filename string) (contentType string, buf io.Reader, err error) {
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
	if stat.Size() > maxsize {
		return "", nil, fmt.Errorf("%s file too large than %d", typ, maxsize)
	}
	var buf bytes.Buffer
	multiwriter := multipart.NewWriter(&buf)
	defer multiwriter.Close()
	w, err := multiwriter.CreateFormFile(fieldname, stat.Name())
	if err != nil {
		return
	}
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	io.Copy(w, r)
	contentType = multiwriter.FormDataContentType()
	return
}

func GetMedia(host, media_id, access_token, dir string) error {
	uri := fmt.Sprintf("https://%s?access_token=%smedia_id=%s",
		host, access_token, media_id)
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	ct := res.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}
	if bytes.HasSuffix(mediaType, "json") {
	}
}
