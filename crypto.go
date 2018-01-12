package weixin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	mrand "math/rand"
	"strconv"
	"time"
)

const (
	wxNonce       = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
	wxNonceLength = 10
	wxAESHeader   = 16
	wxAESLength   = 4
)

// NewAESKey根据生成AES密钥
func NewCipherBlock(encodingAESKey string) (cipher.Block, error) {
	// AES密钥： AESKey=Base64_Decode(EncodingAESKey + “=”), EncodingAESKey尾部填充一个字符的“=”, 用Base64_Decode生成32个字节的AESKey；
	key, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// MessageWithAES 加密消息格式
type EncryptMessage struct {
	// XMLName xml名称
	XMLName xml.Name `xml:"xml" json:"-"`
	// ToUserName 接收者
	ToUserName CDATA
	// Encrypt 加密的文本
	Encrypt CDATA
}

// ReponseWithAES 响应微信服务器的格式
type EncryptResponse struct {
	// XMLName xml名称
	XMLName xml.Name `xml:"xml" json:"-"`
	// Encrypt加密的字符串
	Encrypt CDATA
	// MsgSignature 消息签名
	MsgSignature CDATA
	// TimeStamp 发送事件
	TimeStamp string
	// Nonce 随机字符
	Nonce CDATA
}

func randomString() string {
	bs := make([]byte, wxNonceLength)
	mrand.Seed(time.Now().UnixNano())
	for i := 0; i < len(bs); i++ {
		b := wxNonce[mrand.Intn(len(bs))]
		bs[i] = b
	}
	return string(bs)
}

func NewEncryptResponse(appid, token, nonce, ciphertext string) *EncryptResponse {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if nonce == "" {
		nonce = randomString()
	}
	signature := Sign(token, timestamp, nonce, ciphertext)

	return &EncryptResponse{
		Encrypt:      CDATA(ciphertext),
		MsgSignature: CDATA(signature),
		TimeStamp:    timestamp,
		Nonce:        CDATA(nonce),
	}
}

// Decrypt 解密
func Decrypt(block cipher.Block, ciphertext string) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decrypt: ciphertext %s", err)
	}
	if len(bs) < block.BlockSize() {
		return nil, fmt.Errorf("decrypt: ciphertext too short: %d", len(bs))
	}
	if len(bs)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("decrypt: ciphertext is not a multiple of the block size")
	}
	iv := bs[:aes.BlockSize]
	bs = bs[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(bs, bs)
	return bs, nil
}

func ParseEncryptMessage(b []byte, appid string) ([]byte, string, error) {
	xmlStart := wxAESHeader + wxAESLength
	buf := bytes.NewReader(b[wxAESHeader:xmlStart])

	var length int32
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return nil, "", fmt.Errorf("decrypt: ciphertext when read plaintext length")
	}
	xmlEnd := xmlStart + int(length)
	return b[xmlStart:xmlEnd], string(b[xmlEnd : xmlEnd+len(appid)]), nil
}

// Encrypt 加密
func Encrypt(block cipher.Block, plaintext string, appid string) (string, error) {
	// 随机16位字符
	bs := []byte(plaintext)
	rb := make([]byte, wxAESHeader)
	if _, err := io.ReadFull(rand.Reader, rb); err != nil {
		return "", fmt.Errorf("encrypt when read 16 rand bytes %s", err)
	}
	fmt.Println("rb", len(rb))
	// 取消息的网络长度4字节
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, int32(len(bs))); err != nil {
		return "", fmt.Errorf("encrypt when generate len(plaintext)")
	}
	length := make([]byte, wxAESLength)
	copy(length, buf.Bytes())

	buf.Reset()
	// 写入16位随机字符
	buf.Write(rb)
	// 写入4位网络长度
	buf.Write(length)
	// 写入消息文本
	buf.Write(bs)
	// 写入appid
	buf.Write([]byte(appid))
	bs = buf.Bytes()

	n := aes.BlockSize - len(bs)%aes.BlockSize
	if n != 0 {
		padding := make([]byte, n)
		if _, err := io.ReadFull(rand.Reader, padding); err != nil {
			return "", fmt.Errorf("encrypt when add padding %s", err)
		}
		bs = append(bs, padding...)
	}
	ciphertext := make([]byte, aes.BlockSize+len(bs))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("encrypt when read iv %s", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], bs)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
