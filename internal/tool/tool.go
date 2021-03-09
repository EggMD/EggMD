package tool

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// HashEmail 将电子邮箱地址转换成 MD5 哈希。
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	h := md5.New()
	_, _ = h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

// AvatarLink 根据输入返回对应的 Avatar 头像链接。
func AvatarLink(input string) (url string) {
	if strings.ContainsRune(input, '@') {
		input = HashEmail(input)
	}
	url = "https://cdn.v2ex.com/gravatar/" + input + "?d=identicon"
	return url
}
