package CloudStore

import (
	"github.com/baidubce/bce-sdk-go/services/bos"
)

type BOS struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Domain    string
	client    *bos.Client
}

// new bos
func NewBOS(accessKey, secretKey, bucket, endpoint string) (b *BOS, err error) {
	b = &BOS{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Endpoint:  endpoint,
	}
	b.client, err = bos.NewClient(accessKey, secretKey, endpoint)
	if err != nil {
		return
	}
	return
}

func (b *BOS) IsExist(object string) (err error) {

	return
}

func (b *BOS) Put(tmpFile, saveFile string, header ...map[string]string) (err error) {
	return
}

func (b *BOS) Delete(object ...string) (err error) {
	return
}

func (b *BOS) GetSignURL(object string, expire int64) (link string, err error) {
	return
}
