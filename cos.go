package CloudStore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type COS struct {
	AccessKey string
	SecretKey string
	Bucket    string
	AppID     string
	Region    string
	Domain    string
	Client    *cos.Client
}

func NewCOS(accessKey, secretKey, bucket, appId, region, domain string) (c *COS, err error) {
	c = &COS{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		AppID:     appId,
		Region:    region,
	}
	u, _ := url.Parse(fmt.Sprintf("https://%v-%v.cos.%v.myqcloud.com", bucket, appId, region))
	if domain == "" {
		domain = u.String()
	}
	c.Domain = strings.TrimRight(domain, "/ ")
	c.Client = cos.NewClient(
		&cos.BaseURL{BucketURL: u},
		&http.Client{
			Timeout: 1800 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:  accessKey,
				SecretKey: secretKey,
			},
		})
	return
}

func (c *COS) IsExist(object string) (err error) {
	_, err = c.GetInfo(object)
	return
}

func (c *COS) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	var reader *os.File
	reader, err = os.Open(tmpFile)
	if err != nil {
		return
	}
	defer reader.Close()
	objHeader := &cos.ObjectPutHeaderOptions{}
	for _, header := range headers {
		for k, v := range header {
			switch strings.ToLower(k) {
			case "content-encoding":
				objHeader.ContentEncoding = v
			case "content-type":
				objHeader.ContentType = v
			case "content-disposition":
				objHeader.ContentDisposition = v
			}
		}
	}
	opt := &cos.ObjectPutOptions{ObjectPutHeaderOptions: objHeader}
	_, err = c.Client.Object.Put(context.Background(), objectRel(saveFile), reader, opt)
	return
}

func (c *COS) Delete(objects ...string) (err error) {
	var errs []string
	for _, object := range objects {
		_, err = c.Client.Object.Delete(context.Background(), objectRel(object))
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}
	return
}

func (c *COS) GetSignURL(object string, expire int64) (link string, err error) {
	if expire <= 0 {
		link = c.Domain + objectAbs(object)
		return
	}

	var u *url.URL
	exp := time.Duration(expire) * time.Second
	u, err = c.Client.Object.GetPresignedURL(context.Background(),
		http.MethodGet, objectRel(object),
		c.AccessKey, c.SecretKey,
		exp, nil)
	if err != nil {
		return
	}
	link = u.String()
	if !strings.HasPrefix(link, c.Domain) {
		link = c.Domain + u.RequestURI()
	}
	return
}

func (c *COS) Download(object string, savePath string) (err error) {
	_, err = c.Client.Object.GetToFile(context.Background(), objectRel(object), savePath, nil)
	return
}

func (c *COS) GetInfo(object string) (info File, err error) {
	var resp *cos.Response
	path := objectRel(object)
	resp, err = c.Client.Object.Get(context.Background(), path, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	header := make(map[string]string)
	for k, _ := range resp.Header {
		header[k] = resp.Header.Get(k)
	}
	info = File{
		Header: header,
		Name:   path,
	}
	info.ModTime, _ = time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	info.Size, _ = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	info.IsDir = info.Size == 0
	return
}

func (c *COS) Lists(prefix string) (files []File, err error) {
	// TODO: 腾讯云的SDK中暂时没开放这个功能
	return
}
