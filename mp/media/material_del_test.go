package media

import (
	"testing"
)

func TestDeleteMaterial(t *testing.T) {
	type args struct {
		host        string
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
				mediaID:     "nEemRfdR4u4U9wU3ttFTAtsfJDgB80cE_3GKWf6z-3iVAIUJm6UEWyTwLjxjd_ar",
				accessToken: accessToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteMaterial(tt.args.host, tt.args.accessToken, tt.args.mediaID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMaterial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("DeleteMaterial() = %v, want %v", got, tt.want)
		})
	}
}
