package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketLifecycleFilter ...
type BucketLifecycleFilter struct {
	// 指定规则所适用的前缀。匹配前缀的对象受该规则影响，Prefix 最多只能有一个
	Prefix string                 `xml:"Prefix,omitempty"`
	And    *BucketLifecycleFilter `xml:"And,omitempty"`
}

// BucketLifecycleExpiration ...
type BucketLifecycleExpiration struct {
	// 指明规则对应的动作在何时操作
	Date string `xml:"Date,omitempty"`
	// 指明规则对应的动作在对象最后的修改日期过后多少天操作，该字段有效值为正整数
	Days int `xml:"Days,omitempty"`
	// 删除过期对象删除标记，枚举值 true，false
	ExpiredObjectDeleteMarker bool `xml:"ExpiredObjectDeleteMarker,omitempty"`
	// 指明规则对应的动作在对象变成非当前版本多少天后执行，该字段有效值是正整数
	// 只在作为 NoncurrentVersionExpiration 字段的值时有效
	NoncurrentDays int
}

// BucketLifecycleTransition ...
type BucketLifecycleTransition struct {
	// 指明规则对应的动作在何时操作
	Date string `xml:"Date,omitempty"`
	// 指明规则对应的动作在对象最后的修改日期过后多少天操作，该字段有效值是非负整数
	Days int `xml:"Days,omitempty"`
	// 指定 Object 转储到的目标存储类型，枚举值： STANDARD_IA, ARCHIVE
	StorageClass string
	// 指明规则对应的动作在对象变成非当前版本多少天后执行，该字段有效值是非负整数
	// 只在作为 NoncurrentVersionTransition 字段的值时有效
	NoncurrentDays int
}

// BucketLifecycleAbortIncompleteMultipartUpload ...
type BucketLifecycleAbortIncompleteMultipartUpload struct {
	// 指明分片上传开始后多少天内必须完成上传
	DaysAfterInitiation int `xml:"DaysAfterInitiation,omitempty"`
}

// BucketLifecycleRule ...
//
// https://cloud.tencent.com/document/product/436/8280
type BucketLifecycleRule struct {
	// 用于唯一地标识规则，长度不能超过 255 个字符
	ID string `xml:"ID,omitempty"`
	// Filter 用于描述规则影响的 Object 集合
	Filter *BucketLifecycleFilter
	// 已废弃，改为使用 Filter
	Prefix string
	// 指明规则是否启用，枚举值：Enabled，Disabled
	Status string
	// 规则转换属性，对象何时转换为 Standard_IA 或 Archive
	Transition *BucketLifecycleTransition `xml:"Transition,omitempty"`
	// 规则过期属性
	Expiration *BucketLifecycleExpiration `xml:"Expiration,omitempty"`
	// 设置允许分片上传保持运行的最长时间
	AbortIncompleteMultipartUpload *BucketLifecycleAbortIncompleteMultipartUpload `xml:"AbortIncompleteMultipartUpload,omitempty"`
	// 指明非当前版本对象何时过期
	NoncurrentVersionExpiration *BucketLifecycleExpiration `xml:"NoncurrentVersionExpiration,omitempty"`
	// 指明非当前版本对象何时转换为 STANDARD_IA 或 ARCHIVE
	NoncurrentVersionTransition *BucketLifecycleTransition `xml:"NoncurrentVersionTransition,omitempty"`
}

// BucketGetLifecycleResult ...
//
// https://cloud.tencent.com/document/product/436/8278
type BucketGetLifecycleResult struct {
	XMLName xml.Name              `xml:"LifecycleConfiguration"`
	Rules   []BucketLifecycleRule `xml:"Rule,omitempty"`
}

// GetLifecycle ...
//
// Get Bucket Lifecycle 用来查询 Bucket 的生命周期配置。
//
// https://cloud.tencent.com/document/product/436/8278
func (s *BucketService) GetLifecycle(ctx context.Context) (*BucketGetLifecycleResult, *Response, error) {
	var res BucketGetLifecycleResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?lifecycle",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// BucketPutLifecycleOptions ...
type BucketPutLifecycleOptions struct {
	XMLName xml.Name              `xml:"LifecycleConfiguration"`
	Rules   []BucketLifecycleRule `xml:"Rule,omitempty"`
}

// PutLifecycle ...
//
// COS 支持用户以生命周期配置的方式来管理 Bucket 中 Object 的生命周期。
// 生命周期配置包含一个或多个将应用于一组对象规则的规则集 (其中每个规则为 COS 定义一个操作)。
//
// 这些操作分为以下两种：
//
// 转换操作：定义对象转换为另一个存储类的时间。例如，您可以选择在对象创建 30 天后将其转换为低频存储（STANDARD_IA，适用于不常访问)
// 		   存储类别。同时也支持将数据沉降到归档存储（Archive，成本更低，目前支持国内园区）。具体参数参见请求示例说明中 Transition 项。
// 过期操作：指定 Object 的过期时间。COS 将会自动为用户删除过期的 Object。
//
// 细节分析
//
// PUT Bucket lifecycle 用于为 Bucket 创建一个新的生命周期配置。如果该 Bucket 已配置生命周期，
// 使用该接口创建新的配置的同时则会覆盖原有的配置。
//
// https://cloud.tencent.com/document/product/436/8280
func (s *BucketService) PutLifecycle(ctx context.Context, opt *BucketPutLifecycleOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?lifecycle",
		method:  http.MethodPut,
		body:    opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// DeleteLifecycle ...
//
// Delete Bucket Lifecycle 用来删除 Bucket 的生命周期配置。
//
// https://cloud.tencent.com/document/product/436/8284
func (s *BucketService) DeleteLifecycle(ctx context.Context) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?lifecycle",
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
