package alias

import (
	"crypto/md5"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
)

// HashPass returns encrypted user pass
func HashPass(password string) string {
	// admin pass: 2xEjQA7zyS4KkZ/znructRzO+QIk8o25uLa4+cHiTNbcPt8htog112gc28BSSSyKchkiIPSdNu00vT61BW8KbA==
	md5Str := MD5(password)
	return SHA256(md5Str)
}

// MD5 returns encrypted string using md5
func MD5(str string) string {
	md5Byte := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", md5Byte)
}

// SHA256 returns encrypted string using sha256
func SHA256(str string) string {
	shaByte := sha256.Sum256([]byte(str))
	return string(fmt.Sprintf("%x", shaByte))
}

// Base64Enc encrypts string into BASE64 format
func Base64Enc(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Dec decrypts BASE64-formated string
func Base64Dec(data string) string {
	dec, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("Error in alias: cannot decrypt base64 string")
		return ""
	}
	return string(dec)
}

// StandartRefact returns encrypted string using standart encryption for user struct
func StandartRefact(val string, dec bool, stInfoKey string) string {
	ak := new(AesKey)
	ak.Read(stInfoKey)

	var enc string

	if dec {
		enc = string(DecryptAES([]byte(Base64Dec(val)), ak))
	} else {
		enc = Base64Enc(string(EncryptAES([]byte(val), ak)))
	}

	return enc
}
