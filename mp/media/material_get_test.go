package media

import (
	"testing"
)

func TestGetMaterial(t *testing.T) {
	type args struct {
		host        string
		path        string
		typ         string
		mediaID     string
		accessToken string
		dir         string
	}
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
				path:        WxMaterailGet,
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
				path:        WxMaterailGet,
				typ:         "voice",
				mediaID:     "ivsqEz6azLj5EqfuahV6mX3AagdoCeHqHfzYxZuLVIE",
				accessToken: accessToken,
				dir:         "../../data/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, resp, err := GetMaterial(tt.args.host,
				tt.args.path,
				tt.args.typ,
				tt.args.mediaID,
				tt.args.accessToken,
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
