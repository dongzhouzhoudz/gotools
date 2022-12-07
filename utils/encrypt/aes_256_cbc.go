package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type Aes256Cbc struct {
	Key string
	Iv  string
}

// Encrypt aes-256-cbc 加密方法
func (aes256Cbc *Aes256Cbc) Encrypt(plaintext string) string {
	bKey := []byte(aes256Cbc.Key)
	bIV := []byte(aes256Cbc.Iv)
	bPlaintext := pKCS5Padding([]byte(plaintext), aes.BlockSize, len(plaintext))
	block, _ := aes.NewCipher(bKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt aes-256-cbc 解密方法
func (aes256Cbc *Aes256Cbc) Decrypt(encryptString string) string {
	originData, _ := base64.StdEncoding.DecodeString(encryptString)
	ivBytes := []byte(aes256Cbc.Iv)
	cipherBlock, err := aes.NewCipher([]byte(aes256Cbc.Key))
	if err != nil {
		fmt.Println(err)
	}
	cipher.NewCBCDecrypter(cipherBlock, ivBytes).CryptBlocks(originData, originData)
	return string(pKCS5UnPadding(originData))

}

// pKCS5Padding 加密字符串结果补位
func pKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

// pKCS5UnPadding 解密字符串结果补位
func pKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	if length-unPadding < 0 {
		return []byte("")
	}
	return src[:(length - unPadding)]
}

// GetAes256CbcInstance 获取Aes256Cbc 加密实例
func GetAes256CbcInstance(key string, iv string) *Aes256Cbc {
	return &Aes256Cbc{
		Key: key,
		Iv:  iv,
	}
}
