package CloudStore

import (
	"errors"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

type MinIO struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Domain    string
	Client    *minio.Client
}

func NewMinIO(accessKey, secretKey, bucket, endpoint, domain string) (m *MinIO, err error) {
	if domain == "" {
		domain = "http://" + endpoint
	}
	m = &MinIO{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Endpoint:  endpoint,
		Domain:    domain,
	}
	m.Client, err = minio.New(endpoint, accessKey, secretKey, false)
	m.Domain = strings.TrimRight(m.Domain, "/")
	return
}

func (m *MinIO) IsExist(object string) (err error) {
	_, err = m.GetInfo(object)
	return
}

func (m *MinIO) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	var (
		fp   *os.File
		info os.FileInfo
	)
	fp, err = os.Open(tmpFile)
	if err != nil {
		return
	}
	defer fp.Close()

	info, err = fp.Stat()
	if err != nil {
		return
	}

	opts := minio.PutObjectOptions{
		UserMetadata: make(map[string]string),
	}

	for _, header := range headers {
		for k, v := range header {
			switch strings.ToLower(k) {
			case "content-disposition":
				opts.ContentDisposition = v
			case "content-encoding":
				opts.ContentEncoding = v
			case "content-type":
				opts.ContentType = v
			default:
				opts.UserMetadata[k] = v
			}
		}
	}

	_, err = m.Client.PutObject(m.Bucket, objectRel(saveFile), fp, info.Size(), opts)
	return
}

func (m *MinIO) Delete(objects ...string) (err error) {
	if len(objects) == 0 {
		return
	}

	var errs []string

	objectsChan := make(chan string)
	go func() {
		defer close(objectsChan)
		for _, object := range objects {
			objectsChan <- objectRel(object)
		}
	}()
	for errRm := range m.Client.RemoveObjects(m.Bucket, objectsChan) {
		if errRm.Err != nil {
			errs = append(errs, errRm.Err.Error())
		}
	}
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}
	return
}

func (m *MinIO) GetSignURL(object string, expire int64) (link string, err error) {
	if expire <= 0 {
		link = m.Domain + objectAbs(object)
		return
	}
	if expire > sevenDays {
		expire = sevenDays
	}
	exp := time.Duration(expire) * time.Second
	u := &url.URL{}
	u, err = m.Client.PresignedGetObject(m.Bucket, objectRel(object), exp, nil)
	if err != nil {
		return
	}
	link = u.String()
	if !strings.HasPrefix(link, m.Domain) {
		link = m.Domain + u.RequestURI()
	}
	return
}

func (m *MinIO) Download(object string, savePath string) (err error) {
	obj := &minio.Object{}
	obj, err = m.Client.GetObject(m.Bucket, objectRel(object), minio.GetObjectOptions{})
	if err != nil {
		return
	}

	_, err = obj.Stat()
	if err != nil {
		return
	}

	var fp *os.File
	fp, err = os.Create(savePath)
	if err != nil {
		return
	}
	defer fp.Close()

	_, err = io.Copy(fp, obj)

	return
}

func (m *MinIO) GetInfo(object string) (info File, err error) {
	var objInfo minio.ObjectInfo
	opts := minio.StatObjectOptions{}
	object = objectRel(object)
	objInfo, err = m.Client.StatObject(m.Bucket, object, opts)
	if err != nil {
		return
	}
	info = File{
		ModTime: objInfo.LastModified,
		Name:    object,
		Size:    objInfo.Size,
		IsDir:   objInfo.Size == 0,
		Header:  make(map[string]string),
	}
	for k, _ := range objInfo.Metadata {
		info.Header[k] = objInfo.Metadata.Get(k)
	}
	return
}

func (m *MinIO) Lists(prefix string) (files []File, err error) {
	prefix = objectRel(prefix)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objects := m.Client.ListObjectsV2(m.Bucket, prefix, true, doneCh)
	for object := range objects {
		header := make(map[string]string)
		file := File{
			ModTime: object.LastModified,
			Size:    object.Size,
			IsDir:   object.Size == 0,
			Name:    objectRel(object.Key),
		}
		for k, _ := range object.Metadata {
			header[k] = object.Metadata.Get(k)
		}
		files = append(files, file)
	}
	return
}
