# CloudStore - 云储存集成

国内各大云存储服务接口集成，让云存储使用更方便简单。

目前集成的有：`阿里云OSS`,`百度云BOS`、`腾讯云COS`、`华为云OBS`、`七牛云`、`又拍云`、[Minio](https://www.bookstack.cn/books/MinioCookbookZH)

## 为什么要有这个项目？

为了一劳永逸...

为了变得更懒...

如果上传文件到各大云存储，都变成下面这样:
```
clientBOS.Upload(tmpFile, saveFile)     // 百度云
clientCOS.Upload(tmpFile, saveFile)     // 腾讯云
clientMinio.Upload(tmpFile, saveFile)   // Minio
clientOBS.Upload(tmpFile, saveFile)     // 华为云
clientOSS.Upload(tmpFile, saveFile)     // 阿里云
clientUpYun.Upload(tmpFile, saveFile)   // 又拍云
clientQiniu.Upload(tmpFile, saveFile)   // 七牛云
```

如果各大云存储删除文件对象，都变成下面这样：
```
clientXXX.Delete(file1, file2, file3, ...)
```

不需要翻看各大云存储服务的一大堆文档，除了创建的客户端对象不一样之外，调用的方法和参数都一毛一样，会不会很爽？



## 目前初步实现的功能接口

```
type CloudStore interface {
	Delete(objects ...string) (err error)                                             // 删除文件
	GetSignURL(object string, expire int64) (link string, err error)                  // 文件访问签名
	IsExist(object string) (err error)                                                // 判断文件是否存在
	Lists(prefix string) (files []File, err error)                                    // 文件前缀，列出文件
	Upload(tmpFile string, saveFile string, headers ...map[string]string) (err error) // 上传文件
	Download(object string, savePath string) (err error)                              // 下载文件
	GetInfo(object string) (info File, err error)                                     // 获取指定文件信息
}
```


## 目前集成和实现的功能

- [x] oss - 阿里云云存储 [SDK](https://github.com/aliyun/aliyun-oss-go-sdk) && [文档](https://www.bookstack.cn/books/aliyun-oss-go-sdk)
- [x] cos - 腾讯云云存储 [SDK](https://github.com/tencentyun/cos-go-sdk-v5) && [文档](https://www.bookstack.cn/books/tencent-cos-go-sdk)
- [x] bos - 百度云云存储 [SDK](https://github.com/baidubce/bce-sdk-go) && [文档](https://www.bookstack.cn/books/bos-go-sdk)
- [x] qiniu - 七牛云存储 [SDK](https://github.com/qiniu/api.v7) && [文档](https://www.bookstack.cn/books/qiniu-go-sdk)
- [x] upyun - 又拍云存储 [SDK](https://github.com/upyun/go-sdk) && [文档]()
- [x] obs - 华为云云存储 [SDK](https://support.huaweicloud.com/devg-obs_go_sdk_doc_zh/zh-cn_topic_0142815182.html) && [文档](https://www.bookstack.cn/books/obs-go-sdk)
- [x] minio [SDK](https://github.com/minio/minio-go) && [文档](https://www.bookstack.cn/books/MinioCookbookZH)




TODO: 
- [x] 注意，domain 参数要处理一下，最后统一不带"/"
- [x] 最后获取的签名链接，替换成绑定的域名
- [x] timeout 时间要处理一下，因为一些非内网方式上传文件，在大文件的时候，5分钟或者10分钟都有可能会超时
- [x] `Lists`方法在查询列表的时候，需要对prefix参数做下处理

## 注意
所有云存储的`endpoint`，在配置的时候都是不带 `http://`或者`https://`的

## DocHub 可用云存储
- [x] 百度云 BOS，需要自行压缩svg文件为gzip
- [x] 腾讯云 COS，需要自行压缩svg文件为gzip
- [x] 阿里云 OSS，需要自行压缩svg文件为gzip
- [x] Minio，需要自行压缩svg文件为gzip
- [x] 七牛云存储，在上传svg的时候不需要压缩，svg访问的时候，云存储自行压缩了
- [x] 又拍云，在上传svg的时候不需要压缩，svg访问的时候，云存储自行压缩了
- [x] 华为云 OBS，在上传svg的时候不需要压缩，svg访问的时候，云存储自行压缩了





