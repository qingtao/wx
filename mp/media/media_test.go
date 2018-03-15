package media

import (
	"testing"
)

var mediaID = ""

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
				host:        "api.weixin.qq.com",
				filename:    "ff.jpg",
				accessToken: "7_7AxEenvpgl3LjOQbWOcsX3E8mo2rTEv5-Vc8n9D0nw70faB6Cg4rsXP8S1sj19UJDKHvJY-GJj_y7E7hHgpt4BfF7Ru5Z4nNNHBfoWYKWN1Rv7rtd0wPd7DWDRP3_KSKY0Q1WPCjEVSvI4tzECAhAFAJHF",
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
			mediaID = got.MediaID
			t.Logf("UploadImage() = %v, want %v", got, tt.want)
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
				host:        "api.weixin.qq.com",
				accessToken: `7_7AxEenvpgl3LjOQbWOcsX3E8mo2rTEv5-Vc8n9D0nw70faB6Cg4rsXP8S1sj19UJDKHvJY-GJj_y7E7hHgpt4BfF7Ru5Z4nNNHBfoWYKWN1Rv7rtd0wPd7DWDRP3_KSKY0Q1WPCjEVSvI4tzECAhAFAJHF`,
				dir:         "./",
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
