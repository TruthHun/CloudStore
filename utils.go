package CloudStore

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

var (
	sevenDays int64 = 7 * 24 * 3600
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

func CompressByGzip(tmpFile, saveFile string) (err error) {
	var (
		input []byte
		buf   bytes.Buffer
	)
	input, err = ioutil.ReadFile(tmpFile)
	if err != nil {
		return
	}

	writer, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	defer writer.Close()

	writer.Write(input)
	writer.Flush()

	err = ioutil.WriteFile(saveFile, buf.Bytes(), os.ModePerm)

	return
}

func toJSON(v interface{}) (jsonStr string) {
	p, err := json.Marshal(v)
	if err != nil {
		return
	}
	return string(p)
}
