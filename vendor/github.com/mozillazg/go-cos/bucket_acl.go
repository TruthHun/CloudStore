package cos

import (
	"context"
	"net/http"
)

// BucketGetACLResult ...
//
// https://cloud.tencent.com/document/product/436/7733
type BucketGetACLResult ACLXml

// GetACL 接口用来获取存储桶的访问权限控制列表。
//
// https://cloud.tencent.com/document/product/436/7733
func (s *BucketService) GetACL(ctx context.Context) (*BucketGetACLResult, *Response, error) {
	var res BucketGetACLResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?acl",
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// BucketPutACLOptions ...
// Header 和 Body 二选一
type BucketPutACLOptions struct {
	Header *ACLHeaderOptions `url:"-" xml:"-"`
	Body   *ACLXml           `url:"-" header:"-"`
}

// PutACL 使用API写入Bucket的ACL表
//
// Put Bucket ACL 是一个覆盖操作，传入新的ACL将覆盖原有ACL。只有所有者有权操作。
//
// 私有 Bucket 可以下可以给某个文件夹设置成公有，那么该文件夹下的文件都是公有；
// 但是把文件夹设置成私有后，在该文件夹中设置的公有属性，不会生效。
//
// https://cloud.tencent.com/document/product/436/7737
func (s *BucketService) PutACL(ctx context.Context, opt *BucketPutACLOptions) (*Response, error) {
	header := opt.Header
	body := opt.Body
	if body != nil {
		header = nil
	}
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/?acl",
		method:    http.MethodPut,
		body:      body,
		optHeader: header,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
