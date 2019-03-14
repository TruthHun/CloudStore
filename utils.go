package CloudStore

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 绝对路径，abs => absolute
func objectAbs(object string) string {
	return "/" + strings.TrimLeft(object, " ./")
}

// 相对路径 rel => relative
func objectRel(object string) string {
	return strings.TrimLeft(object, " ./")
}

// MD5 Crypt
func MD5Crypt(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
