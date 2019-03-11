package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketService ...
//
// Bucket 相关 API
type BucketService service

// BucketGetResult 响应结果
//
// https://cloud.tencent.com/document/product/436/7734
type BucketGetResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	// 说明 Bucket 的信息
	Name string
	// 前缀匹配，用来规定响应请求返回的文件前缀地址
	Prefix string `xml:"Prefix,omitempty"`
	// 默认以 UTF-8 二进制顺序列出条目，所有列出条目从 marker 开始
	Marker string `xml:"Marker,omitempty"`
	// 假如返回条目被截断，则返回 NextMarker 就是下一个条目的起点
	NextMarker string `xml:"NextMarker,omitempty"`
	// 定界符，见 BucketGetOptions.Delimiter
	Delimiter string `xml:"Delimiter,omitempty"`
	// 单次响应请求内返回结果的最大的条目数量
	MaxKeys int
	// 响应请求条目是否被截断，布尔值：true，false
	IsTruncated bool
	// 元数据信息
	Contents []Object `xml:"Contents,omitempty"`
	// 将 Prefix 到 delimiter 之间的相同路径归为一类，定义为 Common Prefix
	CommonPrefixes []string `xml:"CommonPrefixes>Prefix,omitempty"`
	// 编码格式
	EncodingType string `xml:"Encoding-Type,omitempty"`
}

// BucketGetOptions 请求参数
//
// https://cloud.tencent.com/document/product/436/7734
type BucketGetOptions struct {
	// 前缀匹配，用来规定返回的文件前缀地址
	Prefix string `url:"prefix,omitempty"`
	// 定界符为一个符号，如果有 Prefix，则将 Prefix 到 delimiter 之间的相同路径归为一类，
	// 定义为 Common Prefix，然后列出所有 Common Prefix。如果没有 Prefix，则从路径起点开始
	Delimiter string `url:"delimiter,omitempty"`
	// 规定返回值的编码方式，可选值：url
	EncodingType string `url:"encoding-type,omitempty"`
	// 默认以 UTF-8 二进制顺序列出条目，所有列出条目从 marker 开始
	Marker string `url:"marker,omitempty"`
	// 单次返回最大的条目数量，默认 1000
	MaxKeys int `url:"max-keys,omitempty"`
}

// Get Bucket 请求等同于 List Object请求，可以列出该 Bucket 下的部分或者全部 Object。
// 此 API 调用者需要对 Bucket 有 Read 权限。
//
// https://cloud.tencent.com/document/product/436/7734
func (s *BucketService) Get(ctx context.Context, opt *BucketGetOptions) (*BucketGetResult, *Response, error) {
	var res BucketGetResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.BucketURL,
		uri:      "/",
		method:   http.MethodGet,
		optQuery: opt,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// BucketPutOptions ...
type BucketPutOptions ACLHeaderOptions

// Put Bucket 接口请求可以在指定账号下创建一个 Bucket。该 API 接口不支持匿名请求，
// 您需要使用帯 Authorization 签名认证的请求才能创建新的 Bucket 。
// 创建 Bucket 的用户默认成为 Bucket 的持有者。
//
// 细节分析
//
// 创建 Bucket 时，如果没有指定访问权限，则默认使用私有读写（private）权限。
//
// https://cloud.tencent.com/document/product/436/7738
func (s *BucketService) Put(ctx context.Context, opt *BucketPutOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/",
		method:    http.MethodPut,
		optHeader: opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// Delete Bucket 请求可以确认该 Bucket 是否存在，是否有权限访问。HEAD 的权限与 Read 一致。
// 当该 Bucket 存在时，返回 HTTP 状态码 200；当该 Bucket 无访问权限时，返回 HTTP 状态码 403；
// 当该 Bucket 不存在时，返回 HTTP 状态码 404。
//
// 注意： 目前还没有公开获取 Bucket 属性的接口（即可以返回 acl 等信息）。
//
// https://cloud.tencent.com/document/product/436/7735
func (s *BucketService) Delete(ctx context.Context) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// Head Bucket请求可以确认是否存在该Bucket，是否有权限访问，Head的权限与Read一致。
//
//   当其存在时，返回 HTTP 状态码200；
//   当无权限时，返回 HTTP 状态码403；
//   当不存在时，返回 HTTP 状态码404。
//
// https://www.qcloud.com/document/product/436/7735
func (s *BucketService) Head(ctx context.Context) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/",
		method:  http.MethodHead,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// Bucket ...
type Bucket struct {
	Name       string
	AppID      string `xml:",omitempty"`
	Region     string `xml:"Location,omitempty"`
	CreateDate string `xml:"CreationDate,omitempty"`
}
