package alias

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"github.com/BurntSushi/toml"
)

// AesKey structure includes key and iv val

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

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
	plaintext, err := pkcs7Pad(plaintext, aes.BlockSize)
	if err != nil {
		return make([]byte, 0), err
	}
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := []byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return []byte(Base64Enc(string(ciphertext))), nil
}

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}
