package mp

import (
	"net/http"
	"testing"
)

const key = "../../data/key.xml"

func TestWx(t *testing.T) {
	wx, err := New(key)
	if err != nil {
		t.Fatalf("---- %s\n", err)
	}
	t.Logf("%#v\n", wx)

	t.Run("Token", func(t *testing.T) {
		if err = wx.GetAccessToken(); err != nil {
			t.Fatal(err)
		}
		t.Logf("access_token: %s\nexpires_in: %d\n", wx.accessToken, wx.expires)
		m, err := wx.GetMenu(wx.accessToken)
		if err != nil {
			t.Fatalf("%#v\n", err)
		}
		t.Logf("%#v\n", m)
		//
		self, err := wx.GetCurrentSelfMenu()
		if err != nil {
			t.Fatalf("%s\n", err)
		}
		t.Logf("%#v\n", self)

		ip, err := wx.GetCallBackIP()
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range ip.IPList {
			t.Logf("%3d: %#v\n", k, v)
		}
	})
	t.Run("Handle", func(t *testing.T) {
		//监听加密事件处理
		http.HandleFunc("/wx", wx.HandleEncryptEvent)
	})
	//运行微信平台测试服务器接口是在腾讯云上进行的，故前端还存在Nginx
	http.ListenAndServe(":30012", nil)
}
