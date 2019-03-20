package CloudStore

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type QINIU struct {
	AccessKey     string
	SecretKey     string
	Bucket        string
	Domain        string
	Zone          *storage.Zone
	mac           *qbox.Mac
	BucketManager *storage.BucketManager
}

func NewQINIU(accessKey, secretKey, bucket, domain string) (q *QINIU, err error) {
	q = &QINIU{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Domain:    domain,
	}
	q.Domain = strings.TrimRight(q.Domain, "/")
	q.mac = qbox.NewMac(accessKey, secretKey)
	q.Zone, err = storage.GetZone(accessKey, bucket)
	if err != nil {
		return
	}
	q.BucketManager = storage.NewBucketManager(q.mac, &storage.Config{Zone: q.Zone})
	return
}

func (q *QINIU) IsExist(object string) (err error) {
	_, err = q.GetInfo(object)
	return
}

// TODO: 目前没发现有可以设置header的地方
func (q *QINIU) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	policy := storage.PutPolicy{Scope: q.Bucket}
	token := policy.UploadToken(q.mac)
	cfg := &storage.Config{
		Zone: q.Zone,
	}
	form := storage.NewFormUploader(cfg)
	ret := &storage.PutRet{}
	params := make(map[string]string)
	for _, header := range headers {
		for k, v := range header {
			params["x:"+k] = v
		}
	}
	extra := &storage.PutExtra{
		Params: params,
	}
	saveFile = objectRel(saveFile)
	// 需要先删除，文件已存在的话，没法覆盖
	q.Delete(saveFile)
	err = form.PutFile(context.Background(), ret, token, saveFile, tmpFile, extra)
	return
}

func (q *QINIU) Delete(objects ...string) (err error) {
	length := len(objects)
	if length == 0 {
		return
	}

	defer func() {
		// 被删除文件不存在的时候，err值为空但不为nil，这里处理一下
		if err != nil && err.Error() == "" {
			err = nil
		}
	}()

	deleteOps := make([]string, 0, length)
	for _, object := range objects {
		deleteOps = append(deleteOps, storage.URIDelete(q.Bucket, objectRel(object)))
	}
	cfg := &storage.Config{
		Zone: q.Zone,
	}
	manager := storage.NewBucketManager(q.mac, cfg)
	var res []storage.BatchOpRet
	res, err = manager.Batch(deleteOps)
	if err != nil {
		return
	}

	var errs []string
	for _, item := range res {
		if item.Code != http.StatusOK {
			errs = append(errs, fmt.Errorf("%+v: %v", item.Data, item.Code).Error())
		}
	}

	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}

	return
}

func (q *QINIU) GetSignURL(object string, expire int64) (link string, err error) {
	object = objectRel(object)
	if expire > 0 {
		deadline := time.Now().Add(time.Second * time.Duration(expire)).Unix()
		link = storage.MakePrivateURL(q.mac, q.Domain, object, deadline)
	} else {
		link = storage.MakePublicURL(q.Domain, object)
	}

	if !strings.HasPrefix(link, q.Domain) {
		if u, errU := url.Parse(link); errU == nil {
			link = q.Domain + u.RequestURI()
		}
	}

	return
}

func (q *QINIU) Download(object string, savePath string) (err error) {
	var link string
	link, err = q.GetSignURL(object, 3600)
	if err != nil {
		return
	}
	req := httplib.Get(link).SetTimeout(30*time.Minute, 30*time.Minute)
	if strings.HasPrefix(strings.ToLower(link), "https://") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	var resp *http.Response

	resp, err = req.Response()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%v: %v", resp.Status, string(data))
	}

	err = ioutil.WriteFile(savePath, data, os.ModePerm)

	return
}

func (q *QINIU) GetInfo(object string) (info File, err error) {
	var fileInfo storage.FileInfo

	object = objectRel(object)
	fileInfo, err = q.BucketManager.Stat(q.Bucket, object)
	if err != nil {
		return
	}
	info = File{
		Name:    object,
		Size:    fileInfo.Fsize,
		ModTime: storage.ParsePutTime(fileInfo.PutTime),
		IsDir:   fileInfo.Fsize == 0,
	}
	return
}

func (q *QINIU) Lists(prefix string) (files []File, err error) {
	var items []storage.ListItem

	prefix = objectRel(prefix)
	limit := 1000
	cfg := &storage.Config{
		Zone: q.Zone,
	}

	manager := storage.NewBucketManager(q.mac, cfg)
	items, _, _, _, err = manager.ListFiles(q.Bucket, prefix, "", "", limit)
	if err != nil {
		return
	}

	for _, item := range items {
		files = append(files, File{
			ModTime: storage.ParsePutTime(item.PutTime),
			Name:    objectRel(item.Key),
			Size:    item.Fsize,
			IsDir:   item.Fsize == 0,
		})
	}

	return
}
