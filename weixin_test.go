package wx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestWx(t *testing.T) {
	wx, err := New("key.txt")
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	if err = wx.GetAccessToken(); err != nil {
		log.Fatalln(err)
	}
	m, err := wx.GetMenu(wx.accessToken)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	fmt.Printf("%#s\n", m)

	/*
		res, err := GetCurrentSelfMenu(api, tk.AccessToken)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
	*/
	http.HandleFunc("/wx", wx.HandleEvent)
	http.ListenAndServe(":80", nil)
}

func TestButton(t *testing.T) {
	var s = `
{
     "button":[
     {    
          "type":"click",
          "name":"今日歌曲",
          "key":"V1001_TODAY_MUSIC"
      },
      {
           "name":"菜单",
           "sub_button":[
           {    
               "type":"view",
               "name":"搜索",
               "url":"http://www.soso.com/"
            },
            {
                 "type":"miniprogram",
                 "name":"wxa",
                 "url":"http://mp.weixin.qq.com",
                 "appid":"wx286b93c14bbf93aa",
                 "pagepath":"pages/lunar/index"
             },
            {
               "type":"click",
               "name":"赞一下我们",
               "key":"V1001_GOOD"
            }]
       }]
}`

	var cb Menu
	if err := json.Unmarshal([]byte(s), &cb); err != nil {
		log.Fatalln(err)
	}
	for i, button := range cb.Button {
		fmt.Printf("%d - name: %s\n", i, button.Name)
		if button.SubButton != nil {
			for j, sub := range button.SubButton {
				fmt.Printf("- sub_button: %d\n", j)
				out("-- name: %s\n", sub.Name)
				out("-- type: %s\n", sub.Type)
				out("-- key: %s\n", sub.Key)
				out("-- url: %s\n", sub.URL)
				out("-- media_id: %s\n", sub.MediaId)
				out("-- appid: %s\n", sub.AppId)
				out("-- pagepath: %s\n", sub.PagePath)
			}
			continue
		}

		out("- type: %s\n", button.Type)
		out("- key: %s\n", button.Key)
		out("- url: %s\n", button.URL)
		out("- media_id: %s\n", button.MediaId)
		out("- appid: %s\n", button.AppId)
		out("- pagepath: %s\n", button.PagePath)
	}
}

func out(f, s string) {
	if s != "" {
		fmt.Printf(f, s)
	}
}

func TestXML(t *testing.T) {
	var e = &Info{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Event:        "event",
		EventKey:     "eventKey",
		MenuID:       "1111",
		ScanCodeInfo: &ScanCodeInfo{
			ScanType:   "1",
			ScanResult: "2",
		},
		SendPicsInfo: &SendPicsInfo{
			Count: 1,
			PicList: []*PicList{
				&PicList{
					&Item{
						PicMd5Sum: "5a75aaca956d97be686719218f275c6b",
					},
				},
			},
		},
	}
	b, err := xml.MarshalIndent(e, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", b)
	fmt.Println("---------------")
	ss := `<xml>
  <ToUserName><![CDATA[toUser]]></ToUserName>
  <FromUserName><![CDATA[fromUser]]></FromUserName>
  <CreateTime>1515467979</CreateTime>
  <Event><![CDATA[event]]></Event>
  <EventKey><![CDATA[eventKey]]></EventKey>
</xml>`
	var ev Info
	if err := xml.Unmarshal([]byte(ss), &ev); err != nil {
		log.Fatalln("%s\n", err)
	}
	fmt.Println("weixin event:")
	fmt.Printf("-- ToUserName: %s\n", ev.ToUserName)
	fmt.Printf("-- ToUserName: %v\n", reflect.TypeOf(ev.ToUserName))
	fmt.Printf("-- FromUserName: %s\n", ev.FromUserName)
	fmt.Printf("-- CreateTime: %d\n", ev.CreateTime)
	fmt.Printf("-- Event: %s\n", ev.Event)
	fmt.Printf("-- EventKey: %s\n", ev.EventKey)
}

func TestInfo(t *testing.T) {
	info := &Info{
		ToUserName:   "gh_763d78092799",
		FromUserName: "okNT7wrJ00zXBowaRS-CAeFcQ7rc",
		CreateTime:   1515567192,
		MsgId:        6509311524975152819,
		MsgType:      "location",
		LocationX:    46.682068,
		LocationY:    217.085495,
		Scale:        16,
		Label:        "a省b市c区d路1号)",
	}
	bj, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", bj)
	b, err := xml.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", b)
	var i Info
	if err = xml.Unmarshal(b, &i); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", i.MsgType)
}
