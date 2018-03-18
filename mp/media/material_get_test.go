package media

import (
	"testing"
)

func TestGetMaterial(t *testing.T) {
	type args struct {
		host        string
		typ         string
		mediaID     string
		accessToken string
		dir         string
	}
	accessToken = "7__wSCo-gaCevTBNcjhMrL57HSgRC7UcR_S3pWnS1m0mQJsJkg7_auLDLr86VaS16ZTjg5md4hL3ZQ04QXV7CW7y0TWUnASifgkv62Z0JnH8y2xlNYOiK0qyiMsc5MMo0tLsZtDZObiWPfWBBRHYQjAFARIA"
	tests := []struct {
		name    string
		args    args
		want    *MaterialGetResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "GetMaterial_video",
			args: args{
				host:        wxhost,
				typ:         "video",
				mediaID:     "ivsqEz6azLj5EqfuahV6mSrcr5mNZX15SUN7EtDaxzI",
				accessToken: accessToken,
				dir:         "../../data/",
			},
			wantErr: false,
		},
		{
			name: "GetMaterial_voice",
			args: args{
				host:        wxhost,
				typ:         "voice",
				mediaID:     "nEemRfdR4u4U9wU3ttFTAtsfJDgB80cE_3GKWf6z-3iVAIUJm6UEWyTwLjxjd_ar",
				accessToken: accessToken,
				dir:         "../../data/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, resp, err := GetMaterial(tt.args.host,
				tt.args.accessToken,
				tt.args.typ,
				tt.args.mediaID,
				tt.args.dir,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaterial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetMaterial() filename = %v, reps = %v, want %v", filename, resp, tt.want)
		})
	}
}
