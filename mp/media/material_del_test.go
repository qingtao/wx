package media

import (
	"testing"
)

func TestDeleteMaterial(t *testing.T) {
	type args struct {
		host        string
		path        string
		mediaID     string
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
			name: "DeleteMaterial",
			args: args{
				host:        wxhost,
				path:        WxMaterialDel,
				mediaID:     "ivsqEz6azLj5EqfuahV6mX3AagdoCeHqHfzYxZuLVIE",
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteMaterial(tt.args.host, tt.args.path, tt.args.mediaID, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMaterial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("DeleteMaterial() = %v, want %v", got, tt.want)
		})
	}
}
