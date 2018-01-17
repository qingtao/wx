package mp

import (
	"log"
	"net/http"
	"testing"
)

func TestWx(t *testing.T) {
	wx, err := New("key.xml")
	if err != nil {
		log.Fatalf("---- %s\n", err)
	}
	log.Printf("%#v\n", wx)

	t.Run("Token", func(t *testing.T) {
		if err = wx.GetAccessToken(); err != nil {
			log.Fatalln(err)
		}
		log.Printf("access_token: %s\nexpires_in: %d\n", wx.accessToken, wx.expires)
		m, err := wx.GetMenu(wx.accessToken)
		if err != nil {
			log.Fatalf("%#v\n", err)
		}
		log.Printf("%#v\n", m)
		//
		self, err := wx.GetCurrentSelfMenu()
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		log.Printf("%#v\n", self)

		ip, err := wx.GetCallBackIP()
		if err != nil {
			log.Fatalln(err)
		}
		for k, v := range ip.IPList {
			log.Printf("%3d: %#v\n", k, v)
		}
	})
	t.Run("Handle", func(t *testing.T) {
		http.HandleFunc("/wx", wx.HandleEncryptEvent)
	})
	http.ListenAndServe(":80", nil)
}
