package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	secret "main/secret"
)

var (
	SecretKey        = secret.Secret.SecretKey
	SecretRefreshKey = secret.Secret.SecretRefreshKey
)

func EncryptToken(tokenString string) (string, error) {
	block, err := aes.NewCipher(SecretKey)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(tokenString))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(tokenString))

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptToken(encryptedToken string) (string, error) {
	block, err := aes.NewCipher(SecretKey)
	if err != nil {
		return "", err
	}

	decodedCipherText, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	iv := decodedCipherText[:aes.BlockSize]
	decodedCipherText = decodedCipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decodedCipherText, decodedCipherText)

	return string(decodedCipherText), nil
}
