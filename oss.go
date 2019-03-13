package CloudStore

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	Key       string
	Secret    string
	Endpoint  string
	Bucket    string
	Domain    string
	bucketObj *oss.Bucket
}

// New OSS
func NewOSS(key, secret, endpoint, bucket, domain string) (o *OSS, err error) {
	var client *oss.Client
	if domain == "" {
		domain = "https://" + bucket + "." + endpoint
	}
	domain = strings.TrimRight(domain, "/")
	o = &OSS{
		Key:      key,
		Secret:   secret,
		Endpoint: endpoint,
		Bucket:   bucket,
		Domain:   domain,
	}
	client, err = oss.New(endpoint, key, secret)
	if err != nil {
		return
	}
	o.bucketObj, err = client.Bucket(bucket)
	return
}

// 注意：阿里云OSS前面，object 前不能是"/"，所以这里需要处理
func (o *OSS) objectToPath(object string) (path string) {
	return strings.TrimLeft(object, "./")
}

func (o *OSS) IsExist(object string) (err error) {
	var b bool
	b, err = o.bucketObj.IsObjectExist(o.objectToPath(object))
	if err != nil {
		return
	}
	if !b {
		return errors.New("file is not exist")
	}
	return
}

func (o *OSS) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	var opts []oss.Option
	for _, header := range headers {
		for k, v := range header {
			switch strings.ToLower(k) {
			case "content-type":
				opts = append(opts, oss.ContentType(v))
			case "content-encoding":
				opts = append(opts, oss.ContentEncoding(v))
			case "content-disposition":
				opts = append(opts, oss.ContentDisposition(v))
				// TODO: more
			}
		}
	}
	err = o.bucketObj.PutObjectFromFile(strings.TrimLeft(saveFile, "./"), tmpFile, opts...)
	return
}

func (o *OSS) Delete(objects ...string) (err error) {
	_, err = o.bucketObj.DeleteObjects(objects)
	return
}

func (o *OSS) GetSignURL(object string, expire int64) (link string, err error) {
	path := o.objectToPath(object)
	if expire <= 0 {
		return o.Domain + "/" + path, nil
	}
	return o.bucketObj.SignURL(path, http.MethodGet, expire)
}

func (o *OSS) Download(object string, savePath string) (err error) {
	err = o.bucketObj.DownloadFile(o.objectToPath(object), savePath, 1048576)
	return
}

func (o *OSS) GetInfo(object string) (info File, err error) {
	// https://help.aliyun.com/document_detail/31859.html?spm=a2c4g.11186623.2.10.713d1592IKig7s#concept-lkf-swy-5db
	//Cache-Control	指定该 Object 被下载时的网页的缓存行为
	//Content-Disposition	指定该 Object 被下载时的名称
	//Content-Encoding	指定该 Object 被下载时的内容编码格式
	//Content-Language	指定该 Object 被下载时的内容语言编码
	//Expires	过期时间
	//Content-Length	该 Object 大小
	//Content-Type	该 Object 文件类型
	//Last-Modified	最近修改时间

	var header http.Header

	path := o.objectToPath(object)
	header, err = o.bucketObj.GetObjectMeta(path)
	if err != nil {
		return
	}

	headerMap := make(map[string]string)

	for k, _ := range header {
		headerMap[k] = header.Get(k)
	}

	info.Header = headerMap
	info.Size, _ = strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	info.ModTime, _ = time.Parse(http.TimeFormat, header.Get("Last-Modified"))
	info.Name = path
	info.IsDir = false
	return
}

func (o *OSS) Lists(prefix string) (files []File, err error) {
	var res oss.ListObjectsResult
	res, err = o.bucketObj.ListObjects(oss.Prefix(o.objectToPath(prefix)))
	if err != nil {
		return
	}
	for _, object := range res.Objects {
		files = append(files, File{
			ModTime: object.LastModified,
			Name:    object.Key,
			Size:    object.Size,
			IsDir:   false,
			Header:  map[string]string{},
		})
	}
	return
}
