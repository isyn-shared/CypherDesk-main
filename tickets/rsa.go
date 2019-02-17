package tickets

import (
	"CypherDesk-main/alias"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func getPublicKeyFromPem(pemPubKey string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemPubKey))
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err.Error)
	}
	return pubKey
}

func getPrivateKeyFromPem(pemPrivKey string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemPrivKey))
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err.Error)
	}
	return privKey
}

func genKeys(fileName string) {
	fileName = "keys/rsa/" + fileName
	privKey := chk(rsa.GenerateKey(rand.Reader, 512)).(*rsa.PrivateKey)
	pubKey := &privKey.PublicKey

	var pemPrivateBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	var pemPublicBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	}
	pemPrivFile := chk(os.Create(fileName + "_pub.pem")).(*os.File)
	pemPubFile := chk(os.Create(fileName + "_priv.pem")).(*os.File)

	err := pem.Encode(pemPrivFile, pemPrivateBlock)
	if err != nil {
		panic("Server error (error code: 09654)")
	}
	pemPrivFile.Close()

	err = pem.Encode(pemPubFile, pemPublicBlock)
	if err != nil {
		panic("Server error (error code: 32423)")
	}
	pemPubFile.Close()
}

func getPemPubKey(fileName string) string {
	return chk(alias.ReadFile("keys/rsa/" + fileName + "_pub.pem")).(string)
}

func getPemPrivKey(fileName string) string {
	return chk(alias.ReadFile("keys/rsa/" + fileName + "_pub.pem")).(string)
}

/* Methods that encrypts using random hash string */
// EncryptWithPublicKey encrypts data with public key
func encryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error on server (error code: 12334)")
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func decryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		panic("Error on server (error code: 12335)")

	}
	return plaintext
}
