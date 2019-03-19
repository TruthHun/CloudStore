package CloudStore

import (
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
	Delete(objects ...string) (err error)                                             // 删除文件
	GetSignURL(object string, expire int64) (link string, err error)                  // 文件访问签名
	IsExist(object string) (err error)                                                // 判断文件是否存在
	Lists(prefix string) (files []File, err error)                                    // 文件前缀，列出文件
	Upload(tmpFile string, saveFile string, headers ...map[string]string) (err error) // 上传文件
	Download(object string, savePath string) (err error)                              // 下载文件
	GetInfo(object string) (info File, err error)                                     // 获取指定文件信息
}
