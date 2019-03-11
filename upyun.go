package CloudStore

import (
	"fmt"
	"os"
	"strings"

	"github.com/upyun/go-sdk/upyun"
)

type UpYun struct {
	Bucket   string
	Operator string
	Password string
	Domain   string
	Client   *upyun.UpYun
}

func NewUpYun(bucket, operator, password, domain string) *UpYun {
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}
	domain = strings.TrimRight(domain, "/")
	client := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   bucket,
		Operator: operator,
		Password: password,
	})
	return &UpYun{
		Bucket:   bucket,
		Operator: operator,
		Password: password,
		Domain:   domain,
		Client:   client,
	}
}

func (u *UpYun) IsExist(object string) (err error) {

	return
}

func (u *UpYun) Put(tmpFile, saveFile string, header ...map[string]string) (err error) {
	_, err = os.Stat(tmpFile)
	// TODO: 设置header
	headers := make(map[string]string)
	if err != nil {
		return
	}
	err = u.Client.Put(&upyun.PutObjectConfig{
		Path:      saveFile,
		LocalPath: tmpFile,
		Headers:   headers,
	})
	return
}

func (u *UpYun) Delete(object ...string) (err error) {
	return
}

func (u *UpYun) GetSignURL(object string, expire int64) (link string, err error) {
	return
}

func (u *UpYun) Lists() {
	objsChan := make(chan *upyun.FileInfo, 10)
	u.Client.List(&upyun.GetObjectsConfig{
		Path:        "/",
		ObjectsChan: objsChan,
	})
	for obj := range objsChan {
		fmt.Printf("%+v\n", obj)
	}
}
