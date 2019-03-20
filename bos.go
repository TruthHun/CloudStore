package CloudStore

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
)

type BOS struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Domain    string
	Client    *bos.Client
}

// new bos
func NewBOS(accessKey, secretKey, bucket, endpoint, domain string) (b *BOS, err error) {
	b = &BOS{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Endpoint:  endpoint,
	}
	if domain == "" {
		domain = "https://" + bucket + "." + endpoint
	}
	b.Domain = strings.TrimRight(domain, "/")
	b.Client, err = bos.NewClient(accessKey, secretKey, endpoint)
	return
}

func (b *BOS) IsExist(object string) (err error) {
	_, err = b.GetInfo(object)
	return
}

func (b *BOS) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	var args = &api.PutObjectArgs{
		UserMeta: make(map[string]string),
	}
	for _, header := range headers {
		for k, v := range header {
			switch strings.ToLower(k) {
			case "content-disposition":
				args.ContentDisposition = v
			case "content-type":
				args.ContentType = v
			case "content-encoding":
				args.ContentEncoding = v
			default:
				args.UserMeta[k] = v
			}
		}
	}
	_, err = b.Client.PutObjectFromFile(b.Bucket, objectRel(saveFile), tmpFile, args)
	return
}

func (b *BOS) Delete(objects ...string) (err error) {
	if len(objects) == 0 {
		return
	}
	for idx, object := range objects {
		objects[idx] = objectRel(object)
	}
	res, _ := b.Client.DeleteMultipleObjectsFromKeyList(b.Bucket, objects)
	if res != nil && len(res.Errors) > 0 {
		err = fmt.Errorf("%+v", res)
	}
	return
}

func (b *BOS) GetSignURL(object string, expire int64) (link string, err error) {
	if expire <= 0 {
		link = b.Domain + objectAbs(object)
	} else {
		link = b.Client.BasicGeneratePresignedUrl(b.Bucket, objectRel(object), int(expire))
		if !strings.HasPrefix(link, b.Domain) {
			if u, errU := url.Parse(link); errU == nil {
				link = b.Domain + u.RequestURI()
			}
		}
	}
	return
}

func (b *BOS) Download(object string, savePath string) (err error) {
	err = b.Client.DownloadSuperFile(b.Bucket, objectRel(object), savePath)
	return
}

func (b *BOS) GetInfo(object string) (info File, err error) {
	var resp *api.GetObjectMetaResult
	resp, err = b.Client.GetObjectMeta(b.Bucket, objectRel(object))
	if err != nil {
		return
	}
	info = File{
		Name:   objectRel(object),
		Size:   resp.ContentLength,
		IsDir:  resp.ContentLength == 0,
		Header: resp.UserMeta,
	}
	info.ModTime, _ = time.Parse(http.TimeFormat, resp.LastModified)
	return
}

func (b *BOS) Lists(prefix string) (files []File, err error) {
	var resp *api.ListObjectsResult
	args := &api.ListObjectsArgs{
		Prefix:  objectRel(prefix),
		MaxKeys: 1000,
	}
	resp, err = b.Client.ListObjects(b.Bucket, args)
	if err != nil {
		return
	}

	for _, object := range resp.Contents {
		file := File{
			Size:  int64(object.Size),
			Name:  objectRel(object.Key),
			IsDir: object.Size == 0,
		}
		file.ModTime, _ = time.Parse(http.TimeFormat, object.LastModified)
		files = append(files, file)
	}
	return
}
