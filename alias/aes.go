package alias

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"github.com/BurntSushi/toml"
)

// AesKey structure includes key and iv val
type AesKey struct {
	Key string
	IV  string
}

func (ak *AesKey) Read(path string) {
	if _, err := toml.DecodeFile(path, &ak); err != nil {
		log.Fatal(err)
	}
}

func createCipher(key []byte) cipher.Block {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to create the AES cipher: %s", err)
	}
	return c
}

// EncryptAES
func EncryptAES(src []byte, ak *AesKey) []byte {
	key := []byte(ak.Key)
	iv := []byte(ak.IV)[:aes.BlockSize]
	blockCipher := createCipher(key)
	stream := cipher.NewCTR(blockCipher, iv)
	stream.XORKeyStream(src, src)
	return src
}

// DecryptAES
func DecryptAES(src []byte, ak *AesKey) []byte {
	key := []byte(ak.Key)
	iv := []byte(ak.IV)[:aes.BlockSize]
	blockCipher := createCipher(key)
	stream := cipher.NewCTR(blockCipher, iv)
	stream.XORKeyStream(src, src)
	return src
}

func GenAESKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

var typicalIV = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}

func EncodeBase64BytesToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64StringToBytes(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// func Encrypt(key, text string) string {
//     block, err := aes.NewCipher([]byte(key))
//     if err != nil { panic(err) }
//     plaintext := []byte(text)
//     cfb := cipher.NewCFBEncrypter(block, iv)
//     ciphertext := make([]byte, len(plaintext))
//     cfb.XORKeyStream(ciphertext, plaintext)
//     return EncodeBase64BytesToString(ciphertext)
// }

func DecryptAESCBC(b64cipherText string, key []byte) ([]byte, error) {
	cipherText := []byte(Base64Dec(b64cipherText))
	block, _ := aes.NewCipher(key)
	if len(cipherText) < aes.BlockSize {
		return cipherText, errors.New("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	if len(cipherText)%aes.BlockSize != 0 {
		return cipherText, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	return cipherText, nil
}

func EncryptAESCBC(plaintext, key []byte) ([]byte, error) {
	if len(plaintext)%aes.BlockSize != 0 {
		return plaintext, errors.New("plaintext is not a multiple of the block size")
	}
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := []byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}
