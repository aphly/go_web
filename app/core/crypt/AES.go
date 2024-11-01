package crypt

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"go_web/app"
)

func AesEn(plaintext string) (string, error) {
	s, err := EncryptAES256ECB(plaintext, app.Config.Http.Appkey)
	if err != nil {
		return "", err
	}
	return s, nil
}

func AesDe(ciphertext string) (string, error) {
	s, err := DecryptAES256ECB(ciphertext, app.Config.Http.Appkey)
	if err != nil {
		return "", err
	}
	return s, nil
}

func EncryptAES256ECB(text string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	plaintext = pad(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	for start := 0; start < len(plaintext); start += block.BlockSize() {
		block.Encrypt(ciphertext[start:start+block.BlockSize()], plaintext[start:start+block.BlockSize()])
	}
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES256ECB(ciphertext string, key []byte) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := make([]byte, len(data))
	// 解密
	for start := 0; start < len(data); start += block.BlockSize() {
		block.Decrypt(plaintext[start:start+block.BlockSize()], data[start:start+block.BlockSize()])
	}
	plaintext = unpad(plaintext)
	return string(plaintext), nil
}

// pad 对数据进行PKCS#7填充
func pad(buf []byte, blockSize int) []byte {
	padding := blockSize - (len(buf) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(buf, padtext...)
}

// unpad 去除PKCS#7填充
func unpad(buf []byte) []byte {
	length := len(buf)
	if length == 0 {
		return buf
	}
	unpadding := int(buf[length-1])
	return buf[:length-unpadding]
}
