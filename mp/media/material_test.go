package media

import (
	"reflect"
	"testing"
)

var (
	wxhost      = "api.weixin.qq.com"
	image       = "../../data/ff.jpg"
	accessToken = "7_obDjOczo257C3dIzm6dLREs-lloaS0Qs4mEZwLS2skqX1Dd_erJf1gJGLtfs1r8B6-N_hx91po2A9BT-pjvMz9JLuEak3RAhlQA9J8uzKK6izwNG21znQgLBT-MMHDZXikEk6PouvZeg9DGXLRJfAEAWCU"
)

func TestUploadImage(t *testing.T) {
	type args struct {
		host        string
		filename    string
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *Response
		wantErr bool
	}{

		{
			name: "api",
			args: args{
				host:        wxhost,
				filename:    image,
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadImage(tt.args.host, tt.args.filename, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("UploadImage() = %v, want %v", got, tt.want)
		})
	}
}

func TestUploadVoice(t *testing.T) {
	type args struct {
		host        string
		filename    string
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "UploadVoice",
			args: args{
				host:        wxhost,
				filename:    "../../data/ff.mp3",
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadVoice(tt.args.host, tt.args.filename, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadVoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("UploadVoice() = %v, want %v", got, tt.want)
		})
	}
}

func TestUploadVideo(t *testing.T) {
	type args struct {
		host        string
		filename    string
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "UploadVideo",
			args: args{
				host:        wxhost,
				filename:    "../../data/ff.mp4",
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadVideo(tt.args.host, tt.args.filename, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("UploadVideo() = %v, want %v", got, tt.want)
		})
	}
}

func TestUploadThumb(t *testing.T) {
	type args struct {
		host        string
		filename    string
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "UploadThumb",
			args: args{
				host:        wxhost,
				filename:    "../../data/ff.jpg",
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadThumb(tt.args.host, tt.args.filename, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadThumb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadThumb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMedia(t *testing.T) {
	type args struct {
		host        string
		mediaID     string
		accessToken string
		dir         string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "api",
			args: args{
				mediaID:     `LMAQxwm98BY1LGorrpi5vHa9NbF6wBvQlNoxZliCeHnYWTpXrBu5ZjTVSLXqRd_w`,
				host:        wxhost,
				accessToken: accessToken,
				dir:         "../../data/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMedia(tt.args.host, tt.args.mediaID, tt.args.accessToken, tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetMedia() = %v, want %v", got, tt.want)
		})
	}
}

func TestMaterialArticle_Upload(t *testing.T) {
	type args struct {
		host        string
		accessToken string
	}
	tests := []struct {
		name    string
		m       *MaterialArticle
		args    args
		want    *MaterialResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "MaterialArticle_Upload",
			m: &MaterialArticle{
				[]*Article{
					{
						Title:            "MaterialArticleTest",
						ThumbMediaID:     "ivsqEz6azLj5EqfuahV6mZeqG7uT7r0Mawcovnh4Fdc",
						Author:           "tom",
						Digest:           "没什么内容",
						ShowCoverPic:     1,
						Content:          "正文也没什么内容的",
						ContentSourceURL: "http://mmbiz.qpic.cn/mmbiz_jpg/4ZsURkMgSkhYibLE2pxF1y0Pib41rdoFMU3WcLe3NoAtVPd8yWzF95gChheUa2IVY6ibFUAclMvwYbAIH37Usnhgg/0",
					},
				},
			},
			args: args{
				host:        wxhost,
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upload(tt.args.host, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaterialArticle.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("MaterialArticle.Upload() = %v, want %v", got, tt.want)
		})
	}
}

func TestMaterialImage_Upload(t *testing.T) {
	type args struct {
		host        string
		accessToken string
	}

	tests := []struct {
		name    string
		m       *MaterialImage
		args    args
		want    *MaterialResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "MaterialImage_Upload_1",
			m: &MaterialImage{
				InMaterial: true,
				FileName:   image,
			},
			args:    args{wxhost, accessToken},
			wantErr: false,
		},
		{
			name: "MaterialImage_Upload_2",
			m: &MaterialImage{
				InMaterial: false,
				FileName:   image,
			},
			args:    args{wxhost, accessToken},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upload(tt.args.host, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaterialImage.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("MaterialImage.Upload() = %v, want %v", got, tt.want)
		})
	}
}

func TestMaterialVideo_Upload(t *testing.T) {
	type args struct {
		host        string
		accessToken string
	}
	tests := []struct {
		name    string
		m       *MaterialVideo
		args    args
		want    *MaterialResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "MaterialVideo_UpLoad",
			m: &MaterialVideo{
				FileName:     "../../data/ff.mp4",
				Title:        "person",
				Introduction: "ok",
			},
			args: args{
				wxhost, accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upload(tt.args.host, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaterialVideo.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("MaterialVideo.Upload() = %v, want %v", got, tt.want)
		})
	}
}

func TestMaterialVoice_Upload(t *testing.T) {
	type args struct {
		host        string
		accessToken string
	}
	tests := []struct {
		name    string
		m       *MaterialVoice
		args    args
		want    *MaterialResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "TestMaterialVoice_Upload",
			wantErr: false,
			m: &MaterialVoice{
				FileName: "../../data/ff.mp3",
			},
			args: args{wxhost, accessToken},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upload(tt.args.host, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaterialVoice.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("MaterialVoice.Upload() = %v, want %v", got, tt.want)
		})
	}
}

func TestMaterialthumb_Upload(t *testing.T) {
	type args struct {
		host        string
		accessToken string
	}
	tests := []struct {
		name    string
		m       *Materialthumb
		args    args
		want    *MaterialResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "TestMaterialthumb_Upload",
			wantErr: false,
			args:    args{wxhost, accessToken},
			m: &Materialthumb{
				FileName: "../../data/ff.jpg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Upload(tt.args.host, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Materialthumb.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Materialthumb.Upload() = %v, want %v", got, tt.want)
		})
	}
}
