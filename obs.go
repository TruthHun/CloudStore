package CloudStore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"obs"
)

type OBS struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Domain    string
	Client    *obs.ObsClient
}

func NewOBS(accessKey, secretKey, bucket, endpoint, domain string) (o *OBS, err error) {
	o = &OBS{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
		Bucket:    bucket,
	}

	if domain == "" {
		domain = fmt.Sprintf("https://%v.%v", bucket, endpoint)
	}
	o.Domain = strings.TrimRight(domain, "/")
	o.Client, err = obs.New(accessKey, secretKey, endpoint)
	return
}

func (o *OBS) IsExist(object string) (err error) {
	_, err = o.GetInfo(object)
	return
}

func (o *OBS) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	var p []byte
	p, err = ioutil.ReadFile(tmpFile)
	if err != nil {
		return
	}

	input := &obs.PutObjectInput{}
	input.Bucket = o.Bucket
	input.Key = objectRel(saveFile)
	input.Metadata = make(map[string]string)
	input.Body = bytes.NewBuffer(p)

	for _, header := range headers {
		for k, v := range header {
			switch strings.ToLower(k) {
			case "content-type":
				input.ContentType = v
			case "content-encoding":
				input.ContentEncoding = v
			case "content-disposition":
				input.ContentDisposition = v
			default:
				input.Metadata[k] = v
			}
		}
	}
	_, err = o.Client.PutObject(input)
	return
}

func (o *OBS) Delete(objects ...string) (err error) {
	if len(objects) <= 0 {
		return
	}
	var objs []obs.ObjectToDelete
	for _, object := range objects {
		objs = append(objs, obs.ObjectToDelete{
			Key: objectRel(object),
		})
	}
	input := &obs.DeleteObjectsInput{
		Bucket:  o.Bucket,
		Objects: objs,
	}
	_, err = o.Client.DeleteObjects(input)
	return
}

func (o *OBS) GetSignURL(object string, expire int64) (link string, err error) {
	if expire <= 0 {
		link = o.Domain + objectAbs(object)
		return
	}
	input := &obs.CreateSignedUrlInput{
		Method:  http.MethodGet,
		Bucket:  o.Bucket,
		Key:     objectRel(object),
		Expires: int(expire),
	}
	output := &obs.CreateSignedUrlOutput{}
	output, err = o.Client.CreateSignedUrl(input)
	if err != nil {
		return
	}
	link = output.SignedUrl
	if !strings.HasPrefix(link, o.Domain) {
		if u, errU := url.Parse(link); errU == nil {
			link = o.Domain + u.RequestURI()
		}
	}
	return
}

func (o *OBS) Download(object string, savePath string) (err error) {
	input := &obs.GetObjectInput{}
	input.Key = objectRel(object)
	input.Bucket = o.Bucket

	output := &obs.GetObjectOutput{}
	output, err = o.Client.GetObject(input)
	if err != nil {
		return
	}
	defer output.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(output.Body)
	if err != nil {
		return
	}

	return ioutil.WriteFile(savePath, b, os.ModePerm)
}

func (o *OBS) GetInfo(object string) (info File, err error) {
	input := &obs.GetObjectMetadataInput{
		Bucket: o.Bucket,
		Key:    objectRel(object),
	}
	output := &obs.GetObjectMetadataOutput{}
	output, err = o.Client.GetObjectMetadata(input)
	if err != nil {
		return
	}
	info = File{
		Name:    objectRel(object),
		Size:    output.ContentLength,
		IsDir:   output.ContentLength == 0,
		ModTime: output.LastModified,
	}
	return
}

func (o *OBS) Lists(prefix string) (files []File, err error) {
	prefix = objectRel(prefix)
	input := &obs.ListObjectsInput{}
	input.Prefix = prefix
	input.Bucket = o.Bucket
	output := &obs.ListObjectsOutput{}
	output, err = o.Client.ListObjects(input)
	if err != nil {
		return
	}

	for _, item := range output.Contents {
		files = append(files, File{
			ModTime: item.LastModified,
			Name:    objectRel(item.Key),
			Size:    item.Size,
			IsDir:   item.Size == 0,
		})
	}

	return
}
