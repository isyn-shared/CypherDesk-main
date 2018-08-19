package alias

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

// HashPass returns encrypted user pass
func HashPass(password string) string {
	md5Byte := md5.Sum([]byte(password))
	md5Str := fmt.Sprintf("%x", md5Byte)

	shaByte := sha256.Sum256([]byte(md5Str))
	return string(fmt.Sprintf("%x", shaByte))
}
