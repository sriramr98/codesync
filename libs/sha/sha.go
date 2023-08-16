package sha

import (
	"crypto/sha1"
	"encoding/hex"
)

func ConvertToShaHex(content []byte) string {
	hash := sha1.New()
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil))
}
