package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// InitiateMultipartUploadOptions ...
type InitiateMultipartUploadOptions struct {
	*ACLHeaderOptions
	*ObjectPutHeaderOptions
}

// InitiateMultipartUploadResult ...
type InitiateMultipartUploadResult struct {
	XMLName xml.Name `xml:"InitiateMultipartUploadResult"`
	// 分片上传的目标 Bucket，由用户自定义字符串和系统生成appid数字串由中划线连接而成，如：mybucket-1250000000
	Bucket string
	// Object 的名称
	Key string
	// 在后续上传中使用的 ID
	UploadID string `xml:"UploadId"`
}

// InitiateMultipartUpload ...
//
// Initiate Multipart Upload请求实现初始化分片上传，成功执行此请求以后会返回Upload ID用于后续的Upload Part请求。
//
// https://cloud.tencent.com/document/product/436/7746
func (s *ObjectService) InitiateMultipartUpload(ctx context.Context, name string, opt *InitiateMultipartUploadOptions) (*InitiateMultipartUploadResult, *Response, error) {
	var res InitiateMultipartUploadResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name) + "?uploads",
		method:    http.MethodPost,
		optHeader: opt,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// ObjectUploadPartOptions ...
type ObjectUploadPartOptions struct {
	// RFC 2616 中定义的 HTTP 请求内容长度（字节）
	Expect          string `header:"Expect,omitempty" url:"-"`
	XCosContentSHA1 string `header:"x-cos-content-sha1" url:"-"`
	// RFC 1864 中定义的经过Base64编码的128-bit 内容 MD5 校验值。此头部用来校验文件内容是否发生变化
	ContentMD5 string `header:"Content-MD5" url:"-"`
	// RFC 2616 中定义的 HTTP 请求内容长度（字节）
	ContentLength int `header:"Content-Length,omitempty" url:"-"`
}

// UploadPart ...
//
// Upload Part 接口请求实现将对象按照分块的方式上传到 COS。
// 最多支持 10000 分块，每个分块大小为 1 MB 到 5 GB ，最后一个分块可以小于 1 MB。
//
// 细节分析
//
// * 分块上传首先需要进行初始化，使用 Initiate Multipart Upload 接口实现，初始化后会得到一个 uploadId ，唯一标识本次上传；
// * 在每次请求 Upload Part 时，需要携带 partNumber 和 uploadId，partNumber 为块的编号，支持乱序上传；
// * 当传入 uploadId 和 partNumber 都相同的时候，后传入的块将覆盖之前传入的块。当 uploadId 不存在时会返回 404 错误，NoSuchUpload。
//
// 当 r 是个 io.ReadCloser 时 UploadPart 方法不会自动调用 r.Close()，用户需要自行选择合适的时机去调用 r.Close() 方法对 r 进行资源回收
//
// https://cloud.tencent.com/document/product/436/7750
func (s *ObjectService) UploadPart(ctx context.Context, name, uploadID string, partNumber int, r io.Reader, opt *ObjectUploadPartOptions) (*Response, error) {
	u := fmt.Sprintf("/%s?partNumber=%d&uploadId=%s", encodeURIComponent(name), partNumber, uploadID)
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       u,
		method:    http.MethodPut,
		optHeader: opt,
		body:      r,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectListPartsOptions ...
type ObjectListPartsOptions struct {
	EncodingType     string `url:"Encoding-type,omitempty"`
	MaxParts         int    `url:"max-parts,omitempty"`
	PartNumberMarker int    `url:"part-number-marker,omitempty"`
}

// ObjectListPartsResult ...
//
// https://cloud.tencent.com/document/product/436/7747
type ObjectListPartsResult struct {
	XMLName xml.Name `xml:"ListPartsResult"`
	// 分块上传的目标 Bucket，存储桶的名字，由用户自定义字符串和系统生成 appid 数字串由中划线连接而成，
	// 如：mybucket-1250000000
	Bucket string
	// 编码格式
	EncodingType string `xml:"Encoding-type,omitempty"`
	// Object 的名字
	Key string
	// 标识本次分块上传的 ID
	UploadID string `xml:"UploadId"`
	// 用来表示这些分块所有者的信息
	Initiator *Initiator `xml:"Initiator,omitempty"`
	// 用来表示这些分块所有者的信息
	Owner *Owner `xml:"Owner,omitempty"`
	// 用来表示这些分块的存储级别，枚举值：STANDARD，STANDARD_IA，ARCHIVE
	StorageClass string
	// 默认以 UTF-8 二进制顺序列出条目，所有列出条目从 marker 开始
	PartNumberMarker int
	// 假如返回条目被截断，则返回 NextMarker 就是下一个条目的起点
	NextPartNumberMarker int `xml:"NextPartNumberMarker,omitempty"`
	// 单次返回最大的条目数量
	MaxParts int
	// 响应请求条目是否被截断，布尔值：true，false
	IsTruncated bool
	// 元数据信息
	Parts []Object `xml:"Part,omitempty"`
}

// ListParts ...
//
// List Parts 用来查询特定分块上传中的已上传的块，即罗列出指定 UploadId 所属的所有已上传成功的分块。
//
// https://cloud.tencent.com/document/product/436/7747
func (s *ObjectService) ListParts(ctx context.Context, name, uploadID string) (*ObjectListPartsResult, *Response, error) {
	u := fmt.Sprintf("/%s?uploadId=%s", encodeURIComponent(name), uploadID)
	var res ObjectListPartsResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     u,
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// CompleteMultipartUploadOptions ...
//
// https://cloud.tencent.com/document/product/436/7742
type CompleteMultipartUploadOptions struct {
	XMLName xml.Name `xml:"CompleteMultipartUpload"`
	Parts   []Object `xml:"Part"`
}

// CompleteMultipartUploadResult ...
//
// https://cloud.tencent.com/document/product/436/7742
type CompleteMultipartUploadResult struct {
	XMLName xml.Name `xml:"CompleteMultipartUploadResult"`
	// 创建的Object的外网访问域名
	Location string
	// 分块上传的目标Bucket，由用户自定义字符串和系统生成appid数字串由中划线连接而成，
	// 如：mybucket-1250000000
	Bucket string
	// Object的名称
	Key string
	// 合并后对象的唯一标签值，该值不是对象内容的 MD5 校验值，仅能用于检查对象唯一性
	ETag string
}

// CompleteMultipartUpload ...
//
// Complete Multipart Upload用来实现完成整个分块上传。当您已经使用Upload Parts上传所有块以后，你可以用该API完成上传。
// 在使用该API时，您必须在Body中给出每一个块的PartNumber和ETag，用来校验块的准确性。
//
// 由于分块上传的合并需要数分钟时间，因而当合并分块开始的时候，COS就立即返回200的状态码，在合并的过程中，
// COS会周期性的返回空格信息来保持连接活跃，直到合并完成，COS会在Body中返回合并后块的内容。
//
// 当上传块小于1 MB的时候，在调用该请求时，会返回400 EntityTooSmall；
// 当上传块编号不连续的时候，在调用该请求时，会返回400 InvalidPart；
// 当请求Body中的块信息没有按序号从小到大排列的时候，在调用该请求时，会返回400 InvalidPartOrder；
// 当UploadId不存在的时候，在调用该请求时，会返回404 NoSuchUpload。
//
// 建议您及时完成分块上传或者舍弃分块上传，因为已上传但是未终止的块会占用存储空间进而产生存储费用。
//
// https://cloud.tencent.com/document/product/436/7742
func (s *ObjectService) CompleteMultipartUpload(ctx context.Context, name, uploadID string, opt *CompleteMultipartUploadOptions) (*CompleteMultipartUploadResult, *Response, error) {
	u := fmt.Sprintf("/%s?uploadId=%s", encodeURIComponent(name), uploadID)
	var res CompleteMultipartUploadResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     u,
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// AbortMultipartUpload ...
//
// Abort Multipart Upload 用来实现舍弃一个分块上传并删除已上传的块。当您调用 Abort Multipart Upload 时，
// 如果有正在使用这个Upload Parts上传块的请求，则Upload Parts会返回失败。当该UploadID不存在时，会返回404 NoSuchUpload。
//
// 建议您及时完成分块上传或者舍弃分块上传，因为已上传但是未终止的块会占用存储空间进而产生存储费用。
//
// https://cloud.tencent.com/document/product/436/7740
func (s *ObjectService) AbortMultipartUpload(ctx context.Context, name, uploadID string) (*Response, error) {
	u := fmt.Sprintf("/%s?uploadId=%s", encodeURIComponent(name), uploadID)
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     u,
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
