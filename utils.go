package upms

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"
)

func MD5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Decode(str string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func GetRandomSalt() string {
	return GetRandomString(32)
}

func GetRandomString(ln int32) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	strbyte := []byte(str)
	result := make([]byte, ln)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < int(ln); i++ {
		result[i] = str[r.Intn(len(strbyte))]
	}
	return string(result)
}
