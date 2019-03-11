package CloudStore

import (
	"os"
)

// 本地存储
type Local struct {
	Folder    string
	Domain    string
	HeaderExt string
}

// 新建本次存储的文件夹
func NewLocal(folder, domain string) (local *Local, err error) {
	err = os.MkdirAll(folder, os.ModePerm)
	local = &Local{
		Folder:    folder,
		Domain:    domain,
		HeaderExt: ".header.json",
	}
	return
}

func (l *Local) IsExist(object string) (err error) {

	return
}

func (l *Local) Put(tmpFile, saveFile string, header ...map[string]string) (err error) {
	return
}

func (l *Local) Delete(object ...string) (err error) {
	return
}

func (l *Local) GetSignURL(object string, expire int64) (link string, err error) {
	return
}
