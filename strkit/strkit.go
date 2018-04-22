package strkit

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"crypto/sha1"
	"strings"
)

func IsBlank(str string) bool {
	return len(str) == 0
}

func RandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandomNumber(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成32位md5字串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 encodes string to SHA1 hex value
func SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// ShortSHA1 truncates SHA1 string to at most 10
func ShortSHA1(sha1 string) string {
	if len(sha1) > 10 {
		return sha1[:10]
	}
	return sha1

}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func ToUnderLine(s string) string {

	data := make([]byte, 0, len(s)*2)
	num := len(s)
	lastIsUpper := true
	for i := 0; i < num; i++ {
		d := s[i]
		if d >= 'A' && d <= 'Z' {
			if !lastIsUpper {
				data = append(data, '_')
				lastIsUpper = true
			}
		} else {
			lastIsUpper = false
		}

		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
