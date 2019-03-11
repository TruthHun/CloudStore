package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketGetLocationResult ...
type BucketGetLocationResult struct {
	XMLName xml.Name `xml:"LocationConstraint"`
	// 说明 Bucket 所在地域，枚举值参见 可用地域[1] 文档，如：ap-beijing、ap-hongkong、eu-frankfurt 等
	// [1]: https://cloud.tencent.com/document/product/436/6224
	Location string `xml:",chardata"`
}

// GetLocation ...
//
// Get Bucket Location 接口用于获取 Bucket 所在的地域信息，该 GET 操作使用 location 参数返回 Bucket 所在的区域，
// 只有 Bucket 持有者才有该 API 接口的操作权限。
//
// https://cloud.tencent.com/document/product/436/8275
func (s *BucketService) GetLocation(ctx context.Context) (*BucketGetLocationResult, *Response, error) {
	var res BucketGetLocationResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?location",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}
