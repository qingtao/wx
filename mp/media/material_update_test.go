package media

import (
	"testing"
)

func TestGetMaterialCount(t *testing.T) {
	type args struct {
		host        string
		accessToken string
	}
	accessToken =
		"7__wSCo-gaCevTBNcjhMrL57HSgRC7UcR_S3pWnS1m0mQJsJkg7_auLDLr86VaS16ZTjg5md4hL3ZQ04QXV7CW7y0TWUnASifgkv62Z0JnH8y2xlNYOiK0qyiMsc5MMo0tLsZtDZObiWPfWBBRHYQjAFARIA"
	tests := []struct {
		name    string
		args    args
		want    *Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "GetMaterialCount",
			args: args{
				host:        wxhost,
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMaterialCount(tt.args.host, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaterialCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetMaterialCount() = %v, want %v", got, tt.want)
		})
	}
}

func TestGetMaterialList(t *testing.T) {
	type args struct {
		host        string
		accessToken string
		req         *MaterialListRequest
	}
	accessToken =
		"7__wSCo-gaCevTBNcjhMrL57HSgRC7UcR_S3pWnS1m0mQJsJkg7_auLDLr86VaS16ZTjg5md4hL3ZQ04QXV7CW7y0TWUnASifgkv62Z0JnH8y2xlNYOiK0qyiMsc5MMo0tLsZtDZObiWPfWBBRHYQjAFARIA"
	tests := []struct {
		name    string
		args    args
		want    *MaterialList
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "GetMaterialList",
			args: args{
				host:        wxhost,
				accessToken: accessToken,
				req: &MaterialListRequest{
					Type:   "image",
					Offset: 0,
					Count:  10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMaterialList(tt.args.host, tt.args.accessToken, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaterialList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetMaterialList() = %#v, want %v", got, tt.want)
		})
	}
}
