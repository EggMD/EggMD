package tool

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// HashEmail hashes email address to MD5 string.
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	h := md5.New()
	_, _ = h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

// AvatarLink returns relative avatar link to the site domain by given email or avatar hash.
func AvatarLink(input string) (url string) {
	if strings.ContainsRune(input, '@') {
		input = HashEmail(input)
	}
	url = "https://cdn.v2ex.com/gravatar/" + input + "?d=identicon"
	return url
}
