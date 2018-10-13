package alias

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

// HashPass returns encrypted user pass
func HashPass(password string) string {
	// admin pass: 465c194afb65670f38322df087f0a9bb225cc257e43eb4ac5a0c98ef5b3173ac
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
