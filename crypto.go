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
)

const (
	wxAESHeader = 16
	wxAESLength = 4
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
type MessageWithAES struct {
	// XMLName xml名称
	XMLName xml.Name `xml:"xml" json:"-"`
	// ToUserName 接收者
	ToUserName CDATA
	// Encrypt 加密的文本
	Encrypt CDATA
}

// ReponseWithAES 响应微信服务器的格式
type ReponseWithAES struct {
	// XMLName xml名称
	XMLName xml.Name `xml:"xml" json:"-"`
	// Encrypt加密的字符串
	Encrypt CDATA
	// MsgSignature 消息签名
	MsgSignature CDATA
	// TimeStamp 发送事件
	TimeStamp int64
	// Nonce 随机字符
	Nonce CDATA
}

// Decrypt 解密
func Decrypt(block cipher.Block, ciphertext string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("decrypt: ciphertext %s", err)
	}
	if len(text) < block.BlockSize() {
		return "", fmt.Errorf("decrypt: ciphertext too short: %d", len(text))
	}
	if len(text)%block.BlockSize() != 0 {
		return "", fmt.Errorf("decrypt: ciphertext is not a multiple of the block size")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(text, text)

	start := wxAESHeader + wxAESLength
	b := bytes.NewReader(text[wxAESHeader:start])

	var length int32
	if err = binary.Read(b, binary.BigEndian, &length); err != nil {
		return "", fmt.Errorf("decrypt: ciphertext when read plaintext length")
	}
	fmt.Println(length)
	fmt.Printf("%s\n", text)
	return string(text[start : start+int(length)]), nil
}

// Encrypt 加密
func Encrypt(block cipher.Block, plaintext string, appid string) (string, error) {
	// 随机16位字符
	text := []byte(plaintext)
	rb := make([]byte, wxAESHeader)
	if _, err := io.ReadFull(rand.Reader, rb); err != nil {
		return "", fmt.Errorf("encrypt when read 16 rand bytes %s", err)
	}
	// 取消息的网络长度4字节
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, int32(len(text))); err != nil {
		return "", fmt.Errorf("encrypt when generate len(plaintext)")
	}
	length := buf.Bytes()

	buf.Reset()
	// 写入16位随机字符
	buf.Write(rb)
	// 写入4位网络长度
	buf.Write(length)
	// 写入消息文本
	buf.Write(text)
	// 写入appid
	buf.Write([]byte(appid))
	text = buf.Bytes()

	n := aes.BlockSize - len(text)%aes.BlockSize
	if n != 0 {
		padding := make([]byte, n)
		if _, err := io.ReadFull(rand.Reader, padding); err != nil {
			return "", fmt.Errorf("encrypt when add padding %s", err)
		}
		text = append(text, padding...)
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("encrypt when read iv %s", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], text)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
