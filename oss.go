package CloudStore

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type OSS struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Endpoint        string
	Domain          string      //绑定的域名
	bk              *oss.Bucket //bucket
}

func NewOSS(key, secret, bucket, endpoint string, domain ...string) (o *OSS, err error) {
	var (
		client *oss.Client
		bk     *oss.Bucket
	)
	dm := "https://" + endpoint + ".aliyun.com" //TODO
	if len(domain) > 0 {
		dm = domain[0]
	}
	o = &OSS{
		AccessKeyId:     key,
		AccessKeySecret: secret,
		Bucket:          bucket,
		Endpoint:        endpoint,
		Domain:          dm,
	}

	// create bucket object
	client, err = oss.New(endpoint, key, secret)
	if err != nil {
		return
	}
	bk, err = client.Bucket(bucket)
	if err != nil {
		return
	}
	o.bk = bk
	return
}

//PutObject(source, destination string, returnInfo ...bool) (ObjectInfo, error) //上传文件
//RenameObject(source, destination string) (ObjectInfo, error)                  //重命名文件对象
//CopyObject(source, destination string) (ObjectInfo, error)                    //拷贝文件对象
//ListObjects(page, pageSize int)                                               //罗列文件
//DelObjects(objects ...string) Errors                                          //删除文件
//GetObjectURL(object string) (string, error)                                   //获取文件URL链接
//GetObjectURLWithSign(object string) (string, error)                           //获取文件URL链接
//GetObjectMeta()                                                               //获取文件的元信息
//SetObjectMeta(object string, meta map[string]string) error                    //设置文件的meta
//Down2Object(link, object string) (ObjectInfo, error)                          //下载文件到云存储，支持自定义文件名，否则返回默认的文件名
//Migrate(source, destination interface{}, kv map[string]string, thread ...int) /

func (o *OSS) DelObjects(objects ...string) (errs map[string]error) {
	errs = make(map[string]error)
	o.bk.DeleteObjects(objects)
	return
}
