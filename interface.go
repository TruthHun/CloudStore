package CloudStore

type CloudStore interface {
	Put(string, string, ...map[string]string) error // 上传文件
	Delete(...string) error                         // 删除文件
	GetSignURL(string, int64) (string, error)       // 文件访问签名
	IsExist(string) error                           // 判断文件是否存在
}
