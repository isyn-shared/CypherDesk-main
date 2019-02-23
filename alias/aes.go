package alias

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
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

func EncryptAESWithRandomIV(key []byte, message string) (encMsg string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	encMsg = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func DecryptAESWithRandomIV(key []byte, cipherText []byte) []byte {
	ak := &AesKey{
		Key: string(key),
		IV:  "104dd450a8536365",
	}
	return DecryptAES(cipherText, ak)
}
