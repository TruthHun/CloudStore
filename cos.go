package CloudStore

import (
	"fmt"
	"net/http"

	"github.com/mozillazg/go-cos"
)

type COS struct {
	SecretId  string
	SecretKey string
	Bucket    string
	AppID     string
	Region    string
	Domain    string
	client    *cos.Client
}

func NewCOS(secretId, secretKey, bucket, appId, region string) (c *COS, err error) {
	var baseURL *cos.BaseURL

	c = &COS{
		SecretId:  secretId,
		SecretKey: secretKey,
		Bucket:    bucket,
		AppID:     appId,
		Region:    region,
	}

	//	https://wafer-1251298948.cos.ap-guangzhou.myqcloud.com/
	baseURL, err = cos.NewBaseURL(fmt.Sprintf("https://%v-%v.cos.%v.myqcloud.com", bucket, appId, region))
	if err != nil {
		return
	}
	c.client = cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		}})
	return
}

func (c *COS) IsExist(object string) (err error) {

	return
}

func (c *COS) Put(tmpFile, saveFile string, header ...map[string]string) (err error) {
	return
}

func (c *COS) Delete(object ...string) (err error) {
	return
}

func (c *COS) GetSignURL(object string, expire int64) (link string, err error) {
	return
}
