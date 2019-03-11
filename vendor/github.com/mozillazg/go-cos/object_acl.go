package cos

import (
	"context"
	"net/http"
)

// ObjectGetACLResult ...
type ObjectGetACLResult ACLXml

// GetACL Get Object ACL接口实现使用API读取Object的ACL表，只有所有者有权操作。
//
// 默认情况下，该 GET 操作返回对象的当前版本。您如果需要返回不同的版本，请使用 version Id 子资源。
//
// https://cloud.tencent.com/document/product/436/7744
func (s *ObjectService) GetACL(ctx context.Context, name string) (*ObjectGetACLResult, *Response, error) {
	var res ObjectGetACLResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/" + encodeURIComponent(name) + "?acl",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// ObjectPutACLOptions ...
type ObjectPutACLOptions struct {
	// Header 和 Body 二选一
	Header *ACLHeaderOptions `url:"-" xml:"-"`
	Body   *ACLXml           `url:"-" header:"-"`
}

// PutACL 使用API写入Object的ACL表。
//
// https://cloud.tencent.com/document/product/436/7748
func (s *ObjectService) PutACL(ctx context.Context, name string, opt *ObjectPutACLOptions) (*Response, error) {
	header := opt.Header
	body := opt.Body
	if body != nil {
		header = nil
	}
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name) + "?acl",
		method:    http.MethodPut,
		optHeader: header,
		body:      body,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
