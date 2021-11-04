package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
)

func Key(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)[0:32]
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Encrypt(key []byte, rawData []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)

	cipherText := make([]byte, blockSize+len(rawData))

	iv := cipherText[:blockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}
