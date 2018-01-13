package weixin

import (
	"fmt"
	"log"
	"testing"
)

const (
	encodingAESKey = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
	to_xml         = `<xml><ToUserName><![CDATA[oia2TjjewbmiOUlr6X-1crbLOvLw]]></ToUserName><FromUserName><![CDATA[gh_7f083739789a]]></FromUserName><CreateTime>1407743423</CreateTime><MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId><Title><![CDATA[testCallBackReplyVideo]]></Title><Description><![CDATA[testCallBackReplyVideo]]></Description></Video></xml>`
	token          = "spamtest"
	nonce          = "1320562132"
	appid          = "wx2c2769f8efd9abc2"

	//
	//
	timestamp = "1409735669"
	msg_sign  = "5d197aaffba7e9b25a30732f161a50dee96bd5fa"
	from_xml  = `hyzAe4OzmOMbd6TvGdIOO6uBmdJoD0Fk53REIHvxYtJlE2B655HuD0m8KUePWB3+LrPXo87wzQ1QLvbeUgmBM4x6F8PGHQHFVAFmOD2LdJF9FrXpbUAh0B5GIItb52sn896wVsMSHGuPE328HnRGBcrS7C41IzDWyWNlZkyyXwon8T332jisa+h6tEDYsVticbSnyU8dKOIbgU6ux5VTjg3yt+WGzjlpKn6NPhRjpA912xMezR4kw6KWwMrCVKSVCZciVGCgavjIQ6X8tCOp3yZbGpy0VxpAe+77TszTfRd5RJSVO/HTnifJpXgCSUdUue1v6h0EIBYYI1BD1DlD+C0CR8e6OewpusjZ4uBl9FyJvnhvQl+q5rv1ixrcpCumEPo5MJSgM9ehVsNPfUM669WuMyVWQLCzpu9GhglF2PE=`
)

func TestAES(t *testing.T) {
	t.Run("enc", func(t *testing.T) {
		ciphertext, err := Encrypt(encodingAESKey, to_xml, appid)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\n", ciphertext)
	})
	t.Run("dec", func(t *testing.T) {
		plaintext, err := Decrypt(encodingAESKey, from_xml)
		if err != nil {
			log.Fatalln(err)
		}
		b, id, err := ParseEncryptMessage(plaintext, appid)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("xml: %s\nid: %s\n", b, id)
	})
}
