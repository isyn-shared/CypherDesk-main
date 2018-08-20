package alias

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

// HashPass returns encrypted user pass
func HashPass(password string) string {
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
