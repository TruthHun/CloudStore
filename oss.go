package CloudStore

import (
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
func NewOSS(key, secret, endpoint, bucket string) (o *OSS, err error) {
	var client *oss.Client
	o = &OSS{
		Key:      key,
		Secret:   secret,
		Endpoint: endpoint,
		Bucket:   bucket,
	}
	client, err = oss.New(endpoint, key, secret)
	if err != nil {
		return
	}
	o.bucketObj, err = client.Bucket(bucket)
	return
}

func (o *OSS) IsExist(object string) (err error) {

	return
}

func (o *OSS) Put(tmpFile, saveFile string, header ...map[string]string) (err error) {
	return
}

func (o *OSS) Delete(object ...string) (err error) {
	return
}

func (o *OSS) GetSignURL(object string, expire int64) (link string, err error) {
	return
}
