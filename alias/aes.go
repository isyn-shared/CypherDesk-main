package alias

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"math/rand"

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
	fmt.Println("DEBUG2: ", key)
	return key
}
