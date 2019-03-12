package CloudStore

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type File struct {
	ModTime time.Time
	Name    string
	Size    int64
	IsDir   bool
	Header  map[string]string
}

type CloudStore interface {
	Delete(...string) error                            // 删除文件
	GetSignURL(string, int64) (string, error)          // 文件访问签名
	IsExist(string) error                              // 判断文件是否存在
	Lists(prefix string) ([]File, error)               // 文件前缀，列出文件
	Upload(string, string, ...map[string]string) error // 上传文件
	//Download(string)//下载文件
	//GetInfo   //获取指定文件信息
}

// MD5 Crypt
func MD5Crypt(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
