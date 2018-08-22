package alias

import (
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"time"
	"unicode/utf8"
)

// NullStr returns converted strign into sql.NullString
func NullStr(str string) sql.NullString {
	var ns sql.NullString
	ns.String = str
	ns.Valid = true
	return ns
}

// StrNull return converted sql.NullString into string
func StrNull(str sql.NullString) (string, error) {
	if str.Valid {
		return str.String, nil
	}
	return "", errors.New("nullString is not valid")
}

// StrLen returns size of string in rune`s
func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

// EmptyStr check if string is empty
func EmptyStr(s string) bool {
	if StrLen(s) == 0 {
		return true
	}
	return false
}

func EmptyStrArr(arr []string) bool {
	for _, str := range arr {
		if EmptyStr(str) {
			return true
		}
	}
	return false
}

// STI returns converted string to int
func STI(str string) (int, error) {
	i, err := strconv.Atoi(str)
	return i, err
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset returns random string with length using alphabet
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
