package cryptoutil

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 将 string 编码为 string 类型 MD5 校验值。
func MD5(str string) string {
	return hex.EncodeToString(MD5Bytes(str))
}

// MD5Bytes 将 string 编码为 []byte 类型 MD5 校验值。
func MD5Bytes(str string) []byte {
	m := md5.New()
	_, _ = m.Write([]byte(str))
	return m.Sum(nil)
}
