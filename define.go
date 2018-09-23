package CloudStore

import "time"

type ObjectInfo struct {
	Object  string
	Size    int64
	Ext     string
	MD5     string
	ModTime time.Time
	URL     string
}

type CloudStore interface {
	PutObject(source, destination string, returnInfo ...bool) (ObjectInfo, error) //上传文件
	RenameObject(source, destination string) (ObjectInfo, error)                  //重命名文件对象
	CopyObject(source, destination string) (ObjectInfo, error)                    //拷贝文件对象
	ListObjects(page, pageSize int)                                               //罗列文件
	DelObjects(objects ...string) map[string]error                                //删除文件
	GetObjectURL(object string) (string, error)                                   //获取文件URL链接
	GetObjectURLWithSign(object string) (string, error)                           //获取文件URL链接
	GetObjectMeta()                                                               //获取文件的元信息
	SetObjectMeta(object string, meta map[string]string) error                    //设置文件的meta
	Down2Object(link, object string) (ObjectInfo, error)                          //下载文件到云存储，支持自定义文件名，否则返回默认的文件名
	Migrate(source, destination interface{}, kv map[string]string, thread ...int) //云存储迁移，比如从OSS到COS或者从COS到OSS等（可以使用down2object来替换。从A的object列表中获取文件，然后通过down2object输入下载连接进行迁移）
}
