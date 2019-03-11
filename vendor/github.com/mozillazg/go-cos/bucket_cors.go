package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketCORSRule ...
//
// https://cloud.tencent.com/document/product/436/8279
type BucketCORSRule struct {
	// 配置规则的 ID
	ID string `xml:"ID,omitempty"`
	// 允许的 HTTP 操作，枚举值：GET，PUT，HEAD，POST，DELETE
	AllowedMethods []string `xml:"AllowedMethod"`
	// 允许的访问来源，支持通配符 * 格式为：协议://域名[:端口] 如：http://www.qq.com
	AllowedOrigins []string `xml:"AllowedOrigin"`
	// 在发送 OPTIONS 请求时告知服务端，接下来的请求可以使用哪些自定义的 HTTP 请求头部，支持通配符 *
	AllowedHeaders []string `xml:"AllowedHeader,omitempty"`
	// 设置 OPTIONS 请求得到结果的有效期
	MaxAgeSeconds int `xml:"MaxAgeSeconds,omitempty"`
	// 设置浏览器可以接收到的来自服务器端的自定义头部信息
	ExposeHeaders []string `xml:"ExposeHeader,omitempty"`
}

// BucketGetCORSResult ...
//
// https://cloud.tencent.com/document/product/436/8274
type BucketGetCORSResult struct {
	XMLName xml.Name `xml:"CORSConfiguration"`
	// 说明跨域资源共享配置的所有信息，最多可以包含100条 CORSRule
	Rules []BucketCORSRule `xml:"CORSRule,omitempty"`
}

// GetCORS ...
//
// Get Bucket CORS 接口实现 Bucket 持有者在 Bucket 上进行跨域资源共享的信息配置。
// （cors 是一个 W3C 标准，全称是"跨域资源共享"（Cross-origin resource sharing））。
// 默认情况下，Bucket 的持有者直接有权限使用该 API 接口，Bucket 持有者也可以将权限授予其他用户。
//
// https://cloud.tencent.com/document/product/436/8274
func (s *BucketService) GetCORS(ctx context.Context) (*BucketGetCORSResult, *Response, error) {
	var res BucketGetCORSResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?cors",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// BucketPutCORSOptions ...
//
// https://cloud.tencent.com/document/product/436/8279
type BucketPutCORSOptions struct {
	XMLName xml.Name `xml:"CORSConfiguration"`
	// 说明跨域资源共享配置的所有信息，最多可以包含 100 条 CORSRule
	Rules []BucketCORSRule `xml:"CORSRule,omitempty"`
}

// PutCORS ...
//
// Put Bucket CORS 接口用来请求设置 Bucket 的跨域资源共享权限，。
// 默认情况下，Bucket 的持有者直接有权限使用该 API 接口，Bucket 持有者也可以将权限授予其他用户。
//
// https://cloud.tencent.com/document/product/436/8279
func (s *BucketService) PutCORS(ctx context.Context, opt *BucketPutCORSOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?cors",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// DeleteCORS ...
//
// Delete Bucket CORS 接口请求实现删除跨域访问配置信息。
//
// https://cloud.tencent.com/document/product/436/8283
func (s *BucketService) DeleteCORS(ctx context.Context) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?cors",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
