package weixin

import (
	"encoding/xml"
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
	from_xml  = `<xml><ToUserName><![CDATA[gh_10f6c3c3ac5a]]></ToUserName><FromUserName><![CDATA[oyORnuP8q7ou2gfYjqLzSIWZf0rs]]></FromUserName><CreateTime>1409735668</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[abcdteT]]></Content><MsgId>6054768590064713728</MsgId><Encrypt><![CDATA[hyzAe4OzmOMbd6TvGdIOO6uBmdJoD0Fk53REIHvxYtJlE2B655HuD0m8KUePWB3+LrPXo87wzQ1QLvbeUgmBM4x6F8PGHQHFVAFmOD2LdJF9FrXpbUAh0B5GIItb52sn896wVsMSHGuPE328HnRGBcrS7C41IzDWyWNlZkyyXwon8T332jisa+h6tEDYsVticbSnyU8dKOIbgU6ux5VTjg3yt+WGzjlpKn6NPhRjpA912xMezR4kw6KWwMrCVKSVCZciVGCgavjIQ6X8tCOp3yZbGpy0VxpAe+77TszTfRd5RJSVO/HTnifJpXgCSUdUue1v6h0EIBYYI1BD1DlD+C0CR8e6OewpusjZ4uBl9FyJvnhvQl+q5rv1ixrcpCumEPo5MJSgM9ehVsNPfUM669WuMyVWQLCzpu9GhglF2PE=]]></Encrypt></xml>`
)

func TestAES(t *testing.T) {
	t.Run("enc", func(t *testing.T) {
		fmt.Printf("enc to_xml:\n%s\n", to_xml)
		ciphertext, err := Encrypt(encodingAESKey, to_xml, appid)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("enc:\n%s\n", ciphertext)

		plaintext, err := Decrypt(encodingAESKey, ciphertext)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("- %s -\n---- %d ----\n", plaintext, plaintext[len(plaintext)-1])

		b, id, err := ParseEncryptMessage(plaintext)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("enc xml and appid:\n%s\n%s\n", b, id)

	})
	t.Run("dec", func(t *testing.T) {
		fmt.Printf("dec from_xml:\n%s\n", from_xml)

		var emsg EncryptMessage
		if err := xml.Unmarshal([]byte(from_xml), &emsg); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("dec emsg:\n%#v\n", emsg)
		plaintext, err := Decrypt(encodingAESKey, string(emsg.Encrypt))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("- %s -\n---- %d ----\n", plaintext, plaintext[len(plaintext)-1])
		b, id, err := ParseEncryptMessage(plaintext)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("dec xml and appid:\n%s\n%s\n", b, id)
	})
	t.Run("Wxdec", func(t *testing.T) {
		fmt.Println("---------------------")
		x := `<xml><Encrypt><![CDATA[3RhKt6TtdQN/H0QfbKewT1KU4xpxD5LIv1BFRyUxjed9mMwcg//sqyBxVehYzVslCxiw6aW46vnH1FZyDD5VeRJY/yLfKqGWkQNfeysoY+THiUfpDEtFmlzcZMQSAiAeUurtfLSO2PLrgqDlzvtRGhA+ZM0/FCGAJChDydr/YoXa7QQ/Q84C6TvFXiA/7FjWnSP8OoGmlh+ahkdyy/qBI2OD0D2Jh6nUFolYL0p0e8cFC2VknBnOZ3zn60ZvWaPWiZFKgDQcmTk9wYKDuj1gMj0WfVGjKTjfzYQu9f7xzXxcFYGA2kLHRS6p9ArvyGnSMP7k5tuwU+2TIaz/2AU72gXb1zF/Lg/L3et1eh9F6oxr/rfVGoSm2HK8JSns9kyX3/jTa/2+J57wzMyKE13yJg==]]></Encrypt><MsgSignature><![CDATA[18c9d2b44c8fe1404bb3cb52999fefe06464cd7c]]></MsgSignature><TimeStamp>1515845750</TimeStamp><Nonce><![CDATA[973497014]]></Nonce></xml>`
		const (
			tken = `tom00123`
			key  = `JWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C`
		)
		var eres EncryptResponse
		if err := xml.Unmarshal([]byte(x), &eres); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\n", eres)
		sign := Sign(tken, eres.TimeStamp, string(eres.Nonce), string(eres.Encrypt))
		if sign != string(eres.MsgSignature) {
			fmt.Printf("sign:\n%s\nmsg_sign:\n%s\n", sign, eres.MsgSignature)
			log.Fatalln("invalid")
		}
		plaintext, err := Decrypt(key, string(eres.Encrypt))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\n", plaintext)
		b, id, err := ParseEncryptMessage(plaintext)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("eres xml and appid:\n%s\n%s\n", b, id)

	})
}
