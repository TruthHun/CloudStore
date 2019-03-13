package CloudStore

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func objectToPath(object string) (path string) {
	return "/" + strings.TrimLeft(object, " ./")
}

// MD5 Crypt
func MD5Crypt(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
