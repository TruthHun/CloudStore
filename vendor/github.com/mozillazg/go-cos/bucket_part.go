package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// ListMultipartUploadsResult ...
//
// https://cloud.tencent.com/document/product/436/7736
type ListMultipartUploadsResult struct {
	XMLName xml.Name `xml:"ListMultipartUploadsResult"`
	// 分块上传的目标 Bucket,由用户自定义字符串和系统生成appid数字串由中划线连接而成，
	// 如：mybucket-1250000000
	Bucket string `xml:"Bucket"`
	// 规定返回值的编码格式，合法值：url
	EncodingType string `xml:"Encoding-Type"`
	// 列出条目从该 key 值开始
	KeyMarker string
	// 列出条目从该 UploadId 值开始
	UploadIDMarker string `xml:"UploadIdMarker"`
	// 假如返回条目被截断，则返回 NextKeyMarker 就是下一个条目的起点
	NextKeyMarker string
	// 假如返回条目被截断，则返回 UploadId 就是下一个条目的起点
	NextUploadIDMarker string `xml:"NextUploadIdMarker"`
	// 设置最大返回的 multipart 数量，合法取值从 0 到 1000
	MaxUploads int
	// 返回条目是否被截断
	IsTruncated bool
	// Upload 列表
	Uploads []MultipartUpload `xml:"Upload,omitempty"`
	// 限定返回的 Object key 必须以 Prefix 作为前缀。
	// 注意使用 prefix 查询时，返回的 key 中仍会包含 Prefix
	Prefix string
	// 定界符为一个符号，对 object 名字包含指定前缀且第一次出现 delimiter 字符之间的
	// object 作为一组元素：common prefix。如果没有 prefix，则从路径起点开始
	Delimiter string `xml:"delimiter,omitempty"`
	// 将 prefix 到 delimiter 之间的相同路径归为一类，定义为 Common Prefix
	CommonPrefixes []string `xml:"CommonPrefixs>Prefix,omitempty"`
}

// ListMultipartUploadsOptions ...
//
// https://cloud.tencent.com/document/product/436/7736
type ListMultipartUploadsOptions struct {
	// https://cloud.tencent.com/document/product/436/7736
	Delimiter string `url:"delimiter,omitempty"`
	// 规定返回值的编码格式，合法值：url
	EncodingType string `url:"encoding-type,omitempty"`
	// 限定返回的 Object key 必须以 Prefix 作为前缀。
	// 注意使用 prefix 查询时，返回的 key 中仍会包含 Prefix
	Prefix string `url:"prefix,omitempty"`
	// 设置最大返回的 multipart 数量，合法取值从1到1000，默认1000
	MaxUploads int `url:"max-uploads,omitempty"`
	// 与 UploadIDMarker 一起使用
	// 当 UploadIDMarker 未被指定时，ObjectName 字母顺序大于 KeyMarker 的条目将被列出
	// 当 UploadIDMarker 被指定时，ObjectName 字母顺序大于 KeyMarker 的条目被列出，
	// ObjectName 字母顺序等于 KeyMarker 同时 UploadID 大于 UploadIDMarker 的条目将被列出。
	KeyMarker      string `url:"key-marker,omitempty"`
	UploadIDMarker string `url:"upload-id-marker,omitempty"`
}

// ListMultipartUploads ...
//
// List Multipart Uploads 用来查询正在进行中的分块上传。单次请求操作最多列出 1000 个正在进行中的分块上传。
//
// 注意：该请求需要有 Bucket 的读权限。
//
// https://cloud.tencent.com/document/product/436/7736
func (s *BucketService) ListMultipartUploads(ctx context.Context, opt *ListMultipartUploadsOptions) (*ListMultipartUploadsResult, *Response, error) {
	var res ListMultipartUploadsResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.BucketURL,
		uri:      "/?uploads",
		method:   http.MethodGet,
		result:   &res,
		optQuery: opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// MultipartUpload 每个 Multipart Upload 的信息
type MultipartUpload struct {
	// Object 的名称
	Key string
	// 标示本次分块上传的 ID
	UploadID string `xml:"UploadID"`
	// 用来表示分块的存储级别，枚举值：STANDARD，STANDARD_IA，ARCHIVE
	StorageClass string
	// 用来表示本次上传发起者的信息
	Initiator *Initiator
	// 用来表示这些分块所有者的信息
	Owner *Owner
	// 分块上传的起始时间
	Initiated string
}
