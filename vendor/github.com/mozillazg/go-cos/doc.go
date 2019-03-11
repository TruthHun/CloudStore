/*
Package cos 腾讯云对象存储服务 COS(Cloud Object Storage) Go SDK。


COS API Version

封装了 V5 版本的 XML API 。


Usage

在项目的 _example 目录下有各个 API 的使用示例[1] 。

[1]: 示例文件所对应的 API 说明 https://github.com/mozillazg/go-cos#todo

备注

* SDK 不会自动设置超时时间，用户根据需要设置合适的超时时间（比如，设置 `http.Client` 的 `Timeout` 字段之类的）或在需要时实现所需的超时机制（比如，通过 `context` 包实现）。

* 所有的 API 在 _example 目录下都有对应的使用示例[1]（示例程序中用到的 `debug` 包只是调试用的不是必需的依赖）。

[1]: 示例文件所对应的 API 说明 https://github.com/mozillazg/go-cos#todo


Authentication

默认所有 API 都是匿名访问. 如果想添加认证信息的话,可以通过自定义一个 http.Client 来添加认证信息.

比如, 使用内置的 AuthorizationTransport 来为请求增加 Authorization Header 签名信息:

	client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  "COS_SECRETID",
				SecretKey: "COS_SECRETKEY",
			},
		})

*/
package cos
