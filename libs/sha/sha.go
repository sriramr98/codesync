package sha

import (
	"crypto/sha1"
	"encoding/base64"
)

func ConvertToShaBase64(content string) string {
	hash := sha1.New()
	hash.Write([]byte(content))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
