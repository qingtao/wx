package mp

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"
	"time"
)

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
		t.Fatal(err)
	}
	for i, button := range cb.Button {
		t.Logf("%d - name: %s\n", i, button.Name)
		if button.SubButton != nil {
			for j, sub := range button.SubButton {
				t.Logf("- sub_button: %d\n", j)
				out(t, "-- name: %s\n", sub.Name)
				out(t, "-- type: %s\n", sub.Type)
				out(t, "-- key: %s\n", sub.Key)
				out(t, "-- url: %s\n", sub.URL)
				out(t, "-- media_id: %s\n", sub.MediaId)
				out(t, "-- appid: %s\n", sub.AppId)
				out(t, "-- pagepath: %s\n", sub.PagePath)
			}
			continue
		}

		out(t, "- type: %s\n", button.Type)
		out(t, "- key: %s\n", button.Key)
		out(t, "- url: %s\n", button.URL)
		out(t, "- media_id: %s\n", button.MediaId)
		out(t, "- appid: %s\n", button.AppId)
		out(t, "- pagepath: %s\n", button.PagePath)
	}
}

func out(t *testing.T, f, s string) {
	if s != "" {
		t.Logf(f, s)
	}
}

func TestXML(t *testing.T) {
	var e = &Message{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Event:        "event",
		EventKey:     "eventKey",
		MenuId:       "1111",
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
		t.Fatal(err)
	}
	t.Logf("%s\n", b)
	t.Logf("---------------\n")
	ss := `<xml>
  <ToUserName><![CDATA[toUser]]></ToUserName>
  <FromUserName><![CDATA[fromUser]]></FromUserName>
  <CreateTime>1515467979</CreateTime>
  <Event><![CDATA[event]]></Event>
  <EventKey><![CDATA[eventKey]]></EventKey>
</xml>`
	var msg Message
	if err := xml.Unmarshal([]byte(ss), &msg); err != nil {
		t.Fatalf("%s\n", err)
	}
	t.Logf("weixin event:\n")
	t.Logf("-- ToUserName: %s\n", msg.ToUserName)
	t.Logf("-- ToUserName: %v\n", reflect.TypeOf(msg.ToUserName))
	t.Logf("-- FromUserName: %s\n", msg.FromUserName)
	t.Logf("-- CreateTime: %d\n", msg.CreateTime)
	t.Logf("-- Event: %s\n", msg.Event)
	t.Logf("-- EventKey: %s\n", msg.EventKey)
}

func TestMsg(t *testing.T) {
	info := &Message{
		ToUserName:   "gh_763d78092799",
		FromUserName: "okNT7wrJ00zXBowaRS-CAeFcQ7rc",
		CreateTime:   1515567192,
		MsgId:        6509311524975152819,
		MsgType:      "location",
		Location_X:   46.682068,
		Location_Y:   217.085495,
		Scale:        16,
		Label:        "a省b市c区d路1号)",
	}
	b, err := xml.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", b)
	var msg Message
	if err = xml.Unmarshal(b, &msg); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", msg.MsgType)
}

func TestResMsg(t *testing.T) {
	var msgs = make([]*ResponseMessage, 0)
	t.Run("text", func(t *testing.T) {
		msg := NewTextMessage("to_user", "from_user", "text_msg")
		/*
			b, err := xml.MarshalIndent(msg, "", "  ")
			if err != nil {
				log.Fatalln(err)
			}
		*/
		msgs = append(msgs, msg)
	})
	t.Run("image", func(t *testing.T) {
		image := NewMedia("image_id", "", "")
		msg := NewImageMessage("to_user", "from_user", image)
		msgs = append(msgs, msg)
	})
	t.Run("voice", func(t *testing.T) {
		voice := NewMedia("voice_id", "", "")
		msg := NewVoiceMessage("to_user", "from_user", voice)
		msgs = append(msgs, msg)
	})
	t.Run("video", func(t *testing.T) {
		video := NewMedia("video_id", "video_title", "video_desc")
		msg := NewVideoMessage("to_user", "from_user", video)
		msgs = append(msgs, msg)
	})
	t.Run("music", func(t *testing.T) {
		music := NewMusic("music_title", "music_desc", "music_url", "music_hqurl", "music_thumb")
		msg := NewMusicMessage("to_user", "from_user", music)
		msgs = append(msgs, msg)
	})
	t.Run("article", func(t *testing.T) {
		as := []*Article{
			{
				"a1_title",
				"a1_desc",
				"a1_picurl",
				"a1_url",
			},
			{
				"a2_title",
				"a2_desc",
				"a2_picurl",
				"a2_url",
			},
		}
		msg := NewArticlesMessage("to_user", "from_user", as)
		msgs = append(msgs, msg)
	})

	for i := 0; i < len(msgs); i++ {
		b, err := xml.MarshalIndent(msgs[i], "", "  ")
		if err != nil {
			t.Fatalf("%3d xml MarshalIndent %s", i, err)
		}
		t.Logf("%3d:\n %s\n------------------------\n", i, b)

	}
}
