package mp

import (
	"encoding/xml"
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
		t.Logf("enc to_xml:\n%s\n", to_xml)
		ciphertext, err := Encrypt(encodingAESKey, appid, []byte(to_xml))
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("enc:\n%s\n", ciphertext)

		plaintext, err := Decrypt(encodingAESKey, ciphertext)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("- %s -\n---- %d ----\n", plaintext, plaintext[len(plaintext)-1])

		b, id, err := ParseDecryptMessage(plaintext)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("enc xml and appid:\n%s\n%s\n", b, id)

	})
	t.Run("dec", func(t *testing.T) {
		t.Logf("dec from_xml:\n%s\n", from_xml)

		var emsg EncryptMessage
		if err := xml.Unmarshal([]byte(from_xml), &emsg); err != nil {
			t.Fatal(err)
		}
		t.Logf("dec emsg:\n%#v\n", emsg)
		plaintext, err := Decrypt(encodingAESKey, string(emsg.Encrypt))
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("- %s -\n---- %d ----\n", plaintext, plaintext[len(plaintext)-1])
		b, id, err := ParseDecryptMessage(plaintext)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("dec xml and appid:\n%s\n%s\n", b, id)
	})
	t.Run("Wxdec", func(t *testing.T) {
		t.Logf("---------------------\n")
		x := `<xml><Encrypt><![CDATA[TYISAfTaVqF97gn22+BQrYOtJcf4GU360iTvjzdLBAp800rTdKFOj+nUhAJKUJ82YA5zbRHPZ/F6P8ok8dYMWhu2zBWwv+xIWlERDlaIKp2CKzbSa5FZ2gl1EWrZzn/GDkKEuDIEY7GyjJaVfiduMg8N6oBlxx6xYz0tuyNlVWoAbgDvIxYYJwkN7CRADuD0IPTE7mkY4fGc56fxFc2D58vnsIOQ8ys28m81fhHZ4g0UjcJKJXofj0N5QTJxO9RmHVP39+b0KcevUacw4Dmi8c72/S0IqIl/eBgZVG4IVVPgpii/7gojmbFDHFi9/RTA4nwQJAe9JxNx+76RvnmXvfCW8PigFlvtsJvY31Nv/ZB97ZxIUiJYk2Y48JoH3wZEs/9NKzDHOaOT+SvRFLaIWeD9DUayXXG5g0vEfugrZPM=]]></Encrypt><MsgSignature><![CDATA[2a57bca3d93100e42599d453d97ebde9f5a57eca]]></MsgSignature><TimeStamp>1515926128</TimeStamp><Nonce><![CDATA[143125540]]></Nonce></xml>`
		const (
			tken = `tom00123`
			key  = `jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C`
		)
		var eres EncryptResponse
		if err := xml.Unmarshal([]byte(x), &eres); err != nil {
			t.Fatal(err)
		}
		t.Logf("%s\n", eres)
		sign := Sign(tken, eres.TimeStamp, string(eres.Nonce), string(eres.Encrypt))
		if sign != string(eres.MsgSignature) {
			t.Logf("sign:\n%s\nmsg_sign:\n%s\n", sign, eres.MsgSignature)
			t.Fatal("invalid")
		}
		plaintext, err := Decrypt(key, string(eres.Encrypt))
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s\n", plaintext)
		b, id, err := ParseDecryptMessage(plaintext)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("eres xml and appid:\n%s\n%s\n", b, id)

	})
}
