package media

import (
	"testing"
)

var (
	image = "../../data/ff.jpg"
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
				host:        "api.weixin.qq.com",
				filename:    image,
				accessToken: "7_AHBHC4MOVpx0IEnEm6dLREs-lloaS0Qs4mEZwKRq2yiCO-d8JbI5YrQgjbhU0U-bUsjNaBTdOwUiUDRLD76r19ydSpK6qi_EFShIeFI3Dksz4b3RKOAe3AM6XHTUVGFWIsqwvLjKOE42pP4nJUHfACAXIV",
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
				accessToken: `7_AHBHC4MOVpx0IEnEm6dLREs-lloaS0Qs4mEZwKRq2yiCO-d8JbI5YrQgjbhU0U-bUsjNaBTdOwUiUDRLD76r19ydSpK6qi_EFShIeFI3Dksz4b3RKOAe3AM6XHTUVGFWIsqwvLjKOE42pP4nJUHfACAXIV`,
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
