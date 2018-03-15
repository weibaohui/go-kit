package strkit

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func IsBlank(str string) bool {
	return len(str) == 0
}
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return GetRandomString(32)
	}
	return id.String()
}
