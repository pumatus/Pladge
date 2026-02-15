package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5加密字符串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}
