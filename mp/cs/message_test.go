package cs

import (
	"encoding/json"
	"testing"
)

func TestMsg(t *testing.T) {
	var msgs = make([]*Message, 0)
	t.Run("text", func(t *testing.T) {
		msg := NewTextMessage("to_user", "from_user", "text_msg")
		msgs = append(msgs, msg)
	})
	t.Run("image", func(t *testing.T) {
		image := NewMedia("image_id", "", "", "")
		msg := NewImageMessage("to_user", "from_user", image)
		msgs = append(msgs, msg)
	})
	t.Run("voice", func(t *testing.T) {
		voice := NewMedia("voice_id", "", "", "")
		msg := NewVoiceMessage("to_user", "from_user", voice)
		msgs = append(msgs, msg)
	})
	t.Run("video", func(t *testing.T) {
		video := NewMedia("video_id", "video_title", "video_desc", "video_thumb")
		msg := NewVideoMessage("to_user", "from_user", video)
		msgs = append(msgs, msg)
	})
	t.Run("music", func(t *testing.T) {
		music := NewMusic("music_title", "music_desc", "music_url", "music_hqurl", "music_thumb")
		msg := NewMusicMessage("to_user", "from_user", music)
		msgs = append(msgs, msg)
	})
	t.Run("news", func(t *testing.T) {
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
		msg := NewNewsMessage("to_user", "from_user", as)
		msgs = append(msgs, msg)
	})

	t.Run("mpnews", func(t *testing.T) {
		mpnews := NewMedia("mpnews_id", "", "", "")
		msg := NewMpNewsMessage("to_user", "from_user", mpnews)
		msgs = append(msgs, msg)
	})

	t.Run("wxcard", func(t *testing.T) {
		wxcard := NewWxCard("wxcard_id")
		msg := NewWxCardMessage("to_user", "from_user", wxcard)
		msgs = append(msgs, msg)
	})

	t.Run("miniprogrampage", func(t *testing.T) {
		mini := NewMiniProgramPage("miniprogrampage_title", "appid",
			"pagepath", "thumb")
		msg := NewMiniProgramPageMessage("to_user", "from_user", mini)
		msgs = append(msgs, msg)
	})

	for i := 0; i < len(msgs); i++ {
		b, err := json.MarshalIndent(msgs[i], "", "  ")
		if err != nil {
			t.Fatalf("%3d xml MarshalIndent %s", i, err)
		}
		t.Logf("%3d:\n %s\n------------------------\n", i, b)

	}
}
