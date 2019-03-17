# 云储存

- [x] oss - 阿里云云存储
    - SDK:
    - 文档:
- [x] cos - 腾讯云云存储
    - SDK:
    - 文档:
- [x] bos - 百度云云存储
    - SDK:
    - 文档:
- qiniu - 七牛云存储
    - SDK:
    - 文档:
- [x] upyun - 又拍云存储
  - SDK:
  - 文档:  
- [x] obs - 华为云云存储
    - SDK：      https://support.huaweicloud.com/devg-obs_go_sdk_doc_zh/zh-cn_topic_0142815182.html
    - 文档：      https://support.huaweicloud.com/api-obs_go_sdk_api_zh/zh-cn_topic_0142812005.html
- minio
    - SDK:
    - 文档:
- local - 本地化存储




TODO: 
- [ ] 注意，domain 参数要处理一下，最后统一不带"/"
- [ ] 最后获取的签名链接，替换成绑定的域名
- [ ] 压缩一个gzip的svg图片，测试下能否正常访问
- [ ] timeout 时间要处理一下，因为一些非内网方式上传文件，在大文件的时候，5分钟或者10分钟都有可能会超时
- [ ] `Lists`方法在查询列表的时候，需要对prefix参数做下处理

## 注意
所有云存储的`endpoint`，在配置的时候都是不带 `http://`或者`https://`的

## 为什么要有这个项目？

自己开源的两个项目 [BookStack](https://github.com/TruthHun/BookStack)、
[DocHub](https://github.com/TruthHun/DocHub) 使用到云存储，
但是各个云存储服务商的接口都各不相同，使用起来非常麻烦，所以直接弄了这么一个云存储的集合，
设计统一的接口规范，以方便使用。

你不需要看各大云存储的文档，只需要看这里的集成，统一使用各大云存储

目前只是对云存储做一个简单的集成，更多功能集成，有待不断升级迭代.

## 目前实现的接口
以`bucket`为操作单位
- 上传文件对象（同时返回文件信息，如文件大小等）
- 获取文件对象链接地址 (分内网和外网以及签名链接地址)
- 获取文件对象基本信息 (文档大小、扩展名、更新时间)
- 获取文件对象列表
- `header` 设置 （仅对云存储邮箱）
- 删除文件对象
- 下载文件到云存
