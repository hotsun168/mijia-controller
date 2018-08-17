package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

var iv = []byte{0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58, 0x56, 0x2e}

func Encrypt(str string, key []byte) ([]byte, error) {
	encrypted := make([]byte, len(str))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.CryptBlocks(encrypted, []byte(str))
	return encrypted, nil
}

func Decrypt(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	return string(data), nil
}
