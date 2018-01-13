package weixin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	mrand "math/rand"
	"strconv"
	"time"
)

const (
	wxNonceLength = 10
	wxAESHeader   = 16
	wxAESLength   = 4

	randString = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
)

// NewAESKey根据生成AES密钥
func NewCipherBlock(encodingAESKey string) (cipher.Block, error) {
	if len(encodingAESKey) != 43 {
		return nil, errors.New("EncodingAESKey must be 43 bytes")
	}
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

// EncryptMessage微信服务器推送的加密消息
type EncryptMessage struct {
	// XMLName xml名称
	XMLName xml.Name `xml:"xml" json:"-"`
	// ToUserName 接收者
	ToUserName CDATA
	// Encrypt 加密的文本
	Encrypt CDATA
}

// EncryptResponse响应微信服务器的格式
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

// random生成随机字符串
func random(n int) []byte {
	bs := make([]byte, n)
	mrand.Seed(time.Now().UnixNano())
	for i := 0; i < len(bs); i++ {
		bs[i] = randString[mrand.Intn(len(randString))]
	}
	return bs
}

// NewEncryptResponse 使用appid、token、nonce和ciphertext生成加密的应答消息
func NewEncryptResponse(appid, token, nonce, ciphertext string) *EncryptResponse {
	// 获取当前的unix时间戳的字符串
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if nonce == "" {
		nonce = string(random(wxNonceLength))
	}
	// 生成签名
	signature := Sign(token, timestamp, nonce, ciphertext)

	return &EncryptResponse{
		Encrypt:      CDATA(ciphertext),
		MsgSignature: CDATA(signature),
		TimeStamp:    timestamp,
		Nonce:        CDATA(nonce),
	}
}

// Decrypt 解密
func Decrypt(key, ciphertext string) ([]byte, error) {
	block, err := NewCipherBlock(key)
	if err != nil {
		return nil, err
	}

	// base64解码
	bs, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decrypt ciphertext %s", err)
	}
	// 加密文本长度必须大于BlockSize
	if len(bs) < block.BlockSize() {
		return nil, fmt.Errorf("ciphertext too short: %d", len(bs))
	}
	// 加密文本的长度必须是BlockSize的正数倍
	if len(bs)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	iv := bs[:aes.BlockSize]
	//bs = bs[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(bs, bs)
	return bs, nil
}

// ParseEncryptMessage 解析解密后的加密消息的主体和appid
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
func Encrypt(key, plaintext, appid string) (string, error) {
	block, err := NewCipherBlock(key)
	if err != nil {
		return "", err
	}
	// 随机16位字符
	bs := []byte(plaintext)
	rb := random(wxAESHeader)
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
		bs = append(bs, random(n)...)
	}
	ciphertext := make([]byte, len(bs))
	iv := random(aes.BlockSize)
	/*
		iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", fmt.Errorf("encrypt when read iv %s", err)
		}
	*/
	mode := cipher.NewCBCEncrypter(block, iv)
	//mode.CryptBlocks(ciphertext[aes.BlockSize:], bs)
	mode.CryptBlocks(ciphertext, bs)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
