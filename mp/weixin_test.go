package mp

import (
	"net/http"
	"testing"
)

func TestWx(t *testing.T) {
	wx, err := New("key.xml")
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
		http.HandleFunc("/wx", wx.HandleEncryptEvent)
	})
	http.ListenAndServe(":8080", nil)
}
