package CloudStore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/upyun/go-sdk/upyun"
)

type UpYun struct {
	Bucket   string
	Operator string
	Password string
	Domain   string
	Client   *upyun.UpYun
	secret   string
}

func NewUpYun(bucket, operator, password, domain, secret string) *UpYun {
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
		secret:   secret,
	}
}

func (u *UpYun) IsExist(object string) (err error) {
	_, err = u.Client.GetInfo(objectAbs(object))
	return
}

func (u *UpYun) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	_, err = os.Stat(tmpFile)
	h := make(map[string]string)
	if err != nil {
		return
	}
	for _, header := range headers {
		for k, v := range header {
			h[k] = v
		}
	}
	err = u.Client.Put(&upyun.PutObjectConfig{
		Path:      objectAbs(saveFile),
		LocalPath: tmpFile,
		Headers:   h,
	})
	return
}

func (u *UpYun) Delete(objects ...string) (err error) {
	var errs []string
	for _, object := range objects {
		err = u.Client.Delete(&upyun.DeleteObjectConfig{
			Path: objectAbs(object),
		})
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}
	return
}

// https://help.upyun.com/knowledge-base/cdn-token-limite/
func (u *UpYun) GetSignURL(object string, expire int64) (link string, err error) {
	path := objectAbs(object)
	if expire <= 0 {
		return u.Domain + path, nil
	}
	endTime := time.Now().Unix() + expire
	sign := MD5Crypt(fmt.Sprintf("%v&%v&%v", u.secret, endTime, path))
	sign = strings.Join(strings.Split(sign, "")[12:20], "") + fmt.Sprint(endTime)
	return u.Domain + path + "?_upt=" + sign, nil
}

func (u *UpYun) Lists(prefix string) (files []File, err error) {
	chans := make(chan *upyun.FileInfo, 1000)
	prefix = objectRel(prefix)
	u.Client.List(&upyun.GetObjectsConfig{
		Path:        prefix,
		ObjectsChan: chans,
	})
	var file File
	for obj := range chans {
		file = File{
			ModTime: obj.Time,
			Size:    obj.Size,
			IsDir:   obj.IsDir,
			Header:  obj.Meta, // 注意：这里获取不到文件的header
			Name:    filepath.Join(prefix, objectRel(obj.Name)),
		}
		files = append(files, file)
	}
	return
}

func (u *UpYun) Download(object string, savePath string) (err error) {
	_, err = u.Client.Get(&upyun.GetObjectConfig{
		Path:      objectAbs(object),
		LocalPath: savePath,
	})
	return
}

func (u *UpYun) GetInfo(object string) (info File, err error) {
	var fileInfo *upyun.FileInfo
	fileInfo, err = u.Client.GetInfo(objectAbs(object))
	if err != nil {
		return
	}
	info = File{
		ModTime: fileInfo.Time,
		Name:    objectRel(fileInfo.Name),
		Size:    fileInfo.Size,
		IsDir:   fileInfo.IsDir,
		Header:  fileInfo.Meta,
	}
	return
}
