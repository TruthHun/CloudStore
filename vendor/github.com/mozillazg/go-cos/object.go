package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ObjectService ...
//
// Object 相关 API
type ObjectService service

// ObjectGetOptions ...
//
// https://cloud.tencent.com/document/product/436/7753
type ObjectGetOptions struct {
	// 设置响应头部中的 Content-Type 参数
	ResponseContentType string `url:"response-content-type,omitempty" header:"-"`
	// 设置响应头部中的 Content-Language 参数
	ResponseContentLanguage string `url:"response-content-language,omitempty" header:"-"`
	// 设置响应头部中的 Content-Expires 参数
	ResponseExpires string `url:"response-expires,omitempty" header:"-"`
	// 设置响应头部中的 Cache-Control 参数
	ResponseCacheControl string `url:"response-cache-control,omitempty" header:"-"`
	// 设置响应头部中的 Content-Disposition 参数
	ResponseContentDisposition string `url:"response-content-disposition,omitempty" header:"-"`
	// 设置响应头部中的 Content-Encoding 参数
	ResponseContentEncoding string `url:"response-content-encoding,omitempty" header:"-"`
	// RFC 2616 中定义的指定文件下载范围，以字节（bytes）为单位
	Range string `url:"-" header:"Range,omitempty"`
	// 如果文件修改时间早于或等于指定时间，才返回文件内容。否则返回 412 (precondition failed)
	IfUnmodifiedSince string `url:"-" header:"If-Unmodified-Since,omitempty"`
	// 当 Object 在指定时间后被修改，则返回对应 Object meta 信息，否则返回 304(not modified)
	IfModifiedSince string `url:"-" header:"If-Modified-Since,omitempty"`
	// 当 ETag 与指定的内容一致，才返回文件。否则返回 412 (precondition failed)
	IfMatch string `url:"-" header:"If-Match,omitempty"`
	// 当 ETag 与指定的内容不一致，才返回文件。否则返回 304 (not modified)
	IfNoneMatch string `url:"-" header:"If-None-Match,omitempty"`

	// 预签名授权 URL
	PresignedURL *url.URL `header:"-" url:"-" xml:"-"`
}

// Get Object 请接口请求可以在 COS 的存储桶中将一个文件（对象）下载至本地。
// 该操作需要请求者对目标对象具有读权限或目标对象对所有人都开放了读权限（公有读）。
//
// 版本
//
// 当启用多版本，该 GET 操作返回对象的当前版本。要返回不同的版本，请使用 versionId 参数。
//
// 注意
//
// 如果该对象的当前版本是删除标记，则 COS 的行为表现为该对象不存在，并返回响应 x-cos-delete-marker: true。
//
// https://cloud.tencent.com/document/product/436/7753
func (s *ObjectService) Get(ctx context.Context, name string, opt *ObjectGetOptions) (*Response, error) {
	baseURL := s.client.BaseURL.BucketURL
	uri := "/" + encodeURIComponent(name)
	if opt != nil && opt.PresignedURL != nil {
		baseURL = opt.PresignedURL
		uri = ""
	}
	sendOpt := sendOptions{
		baseURL:          baseURL,
		uri:              uri,
		method:           http.MethodGet,
		optQuery:         opt,
		optHeader:        opt,
		disableCloseBody: true,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectPutHeaderOptions ...
type ObjectPutHeaderOptions struct {
	// RFC 2616 中定义的缓存策略，将作为 Object 元数据保存。
	CacheControl string `header:"Cache-Control,omitempty" url:"-"`
	// RFC 2616 中定义的文件名称，将作为 Object 元数据保存。
	ContentDisposition string `header:"Content-Disposition,omitempty" url:"-"`
	// RFC 2616 中定义的编码格式，将作为 Object 元数据保存。
	ContentEncoding string `header:"Content-Encoding,omitempty" url:"-"`
	// RFC 2616 中定义的内容类型（MIME），将作为 Object 元数据保存。
	ContentType string `header:"Content-Type,omitempty" url:"-"`
	//
	ContentLength int `header:"Content-Length,omitempty" url:"-"`
	// RFC 2616 中定义的文件日期和时间，将作为 Object 元数据保存。
	Expect          string `header:"Expect,omitempty" url:"-"`
	Expires         string `header:"Expires,omitempty" url:"-"`
	XCosContentSHA1 string `header:"x-cos-content-sha1,omitempty" url:"-"`
	// 自定义的 x-cos-meta-* header
	// 包括用户自定义头部后缀和用户自定义头部信息，将作为 Object 元数据返回，大小限制为 2KB。
	// 注意：用户自定义头部信息支持下划线，但用户自定义头部后缀不支持下划线。
	XCosMetaXXX *http.Header `header:"x-cos-meta-*,omitempty" url:"-"`
	// 设置 Object 的存储级别，枚举值：STANDARD, STANDARD_IA，默认值：STANDARD
	XCosStorageClass string `header:"x-cos-storage-class,omitempty" url:"-"`

	// XCosServerSideEncryption 用于指定腾讯云 COS 在数据存储时，应用数据加密的保护策略。
	// 腾讯云 COS 会帮助您在数据写入数据中心时自动加密，并在您取用该数据时自动解密。
	// 目前支持使用腾讯云 COS 主密钥对数据进行 AES-256 加密。
	// 如果您需要对数据启用服务端加密，则需指定 XCosServerSideEncryption。
	//
	// 指定将对象启用服务端加密的方式。使用 COS 主密钥加密填写：AES256
	XCosServerSideEncryption string `header:"x-cos-server-side-encryption,omitempty" url:"-"`
	// 可选值: Normal, Appendable
	//XCosObjectType string `header:"x-cos-object-type,omitempty" url:"-"`
}

// ObjectPutOptions ...
type ObjectPutOptions struct {
	*ACLHeaderOptions       `header:",omitempty" url:"-" xml:"-"`
	*ObjectPutHeaderOptions `header:",omitempty" url:"-" xml:"-"`

	// 预签名授权 URL
	PresignedURL *url.URL `header:"-" url:"-" xml:"-"`
}

// Put Object请求可以将一个文件（Object）上传至指定Bucket。
//
// 版本
//
// * 如果对存储桶启用版本控制，对象存储将自动为要添加的对象生成唯一的版本 ID。
//   对象存储使用 x-cos-version-id 响应头部在响应中返回此标识。
// * 如果需要暂停存储桶的版本控制，则对象存储始终将其 null 用作存储在存储桶中的对象的版本 ID。
//
// 细节分析
//
// * 需要有 Bucket 的写权限；
// * 如果请求头的 Content-Length 值小于实际请求体（body）中传输的数据长度，COS 仍将成功创建文件，
//   但 Object 大小只等于 Content-Length 中定义的大小，其他数据将被丢弃；
// * 如果试图添加的 Object 的同名文件已经存在，那么新上传的文件，将覆盖原来的文件，成功时返回 200 OK。
//
// 当 r 是个 io.ReadCloser 时 Put 方法不会自动调用 r.Close()，用户需要自行选择合适的时机去调用 r.Close() 方法对 r 进行资源回收
//
// https://cloud.tencent.com/document/product/436/7749
func (s *ObjectService) Put(ctx context.Context, name string, r io.Reader, opt *ObjectPutOptions) (*Response, error) {
	baseURL := s.client.BaseURL.BucketURL
	uri := "/" + encodeURIComponent(name)
	if opt != nil && opt.PresignedURL != nil {
		baseURL = opt.PresignedURL
		uri = ""
	}
	sendOpt := sendOptions{
		baseURL:   baseURL,
		uri:       uri,
		method:    http.MethodPut,
		body:      r,
		optHeader: opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectCopyHeaderOptions ...
// https://cloud.tencent.com/document/product/436/10881
type ObjectCopyHeaderOptions struct {
	// 是否拷贝源文件的元数据，枚举值：Copy, Replaced，默认值 Copy。假如标记为 Copy，
	// 则拷贝源文件的元数据；假如标记为 Replaced，则按本次请求的 Header 信息修改元数据。
	// 当目标路径和源路径一致，即用户试图修改元数据时，则标记必须为 Replaced。
	XCosMetadataDirective string `header:"x-cos-metadata-directive,omitempty" url:"-" xml:"-"`
	// 当 Object 在指定时间后被修改，则执行操作，否则返回 412。
	// 可与 XCosCopySourceIfNoneMatch 一起使用，与其他条件联合使用返回冲突。
	XCosCopySourceIfModifiedSince string `header:"x-cos-copy-source-If-Modified-Since,omitempty" url:"-" xml:"-"`
	// 当 Object 在指定时间后未被修改，则执行操作，否则返回 412。
	// 可与 XCosCopySourceIfMatch 一起使用，与其他条件联合使用返回冲突。
	XCosCopySourceIfUnmodifiedSince string `header:"x-cos-copy-source-If-Unmodified-Since,omitempty" url:"-" xml:"-"`
	// 当 Object 的 Etag 和给定一致时，则执行操作，否则返回 412。
	// 可与 XCosCopySourceIfUnmodifiedSince 一起使用，与其他条件联合使用返回冲突。
	XCosCopySourceIfMatch string `header:"x-cos-copy-source-If-Match,omitempty" url:"-" xml:"-"`
	// 当 Object 的 Etag 和给定不一致时，则执行操作，否则返回 412。
	// 可与 XCosCopySourceIfModifiedSince 一起使用，与其他条件联合使用返回冲突。
	XCosCopySourceIfNoneMatch string `header:"x-cos-copy-source-If-None-Match,omitempty" url:"-" xml:"-"`
	// 设置 Object 的存储级别，枚举值：STANDARD，STANDARD_IA。默认值：STANDARD
	XCosStorageClass string `header:"x-cos-storage-class,omitempty" url:"-" xml:"-"`
	// 自定义的 x-cos-meta-* header
	XCosMetaXXX *http.Header `header:"x-cos-meta-*,omitempty" url:"-"`
	// 源文件 URL 路径，可以通过 versionid 子资源指定历史版本
	XCosCopySource string `header:"x-cos-copy-source" url:"-" xml:"-"`

	// XCosServerSideEncryption 用于指定腾讯云 COS 在数据存储时，应用数据加密的保护策略。
	// 腾讯云 COS 会帮助您在数据写入数据中心时自动加密，并在您取用该数据时自动解密。
	// 目前支持使用腾讯云 COS 主密钥对数据进行 AES-256 加密。
	// 如果您需要对数据启用服务端加密，则需指定 XCosServerSideEncryption。
	//
	// 指定将对象启用服务端加密的方式。使用 COS 主密钥加密填写：AES256
	XCosServerSideEncryption string `header:"x-cos-server-side-encryption,omitempty" url:"-"`
}

// ObjectCopyOptions ...
//
// https://cloud.tencent.com/document/product/436/10881
type ObjectCopyOptions struct {
	*ObjectCopyHeaderOptions `header:",omitempty" url:"-" xml:"-"`
	*ACLHeaderOptions        `header:",omitempty" url:"-" xml:"-"`
}

// ObjectCopyResult ...
type ObjectCopyResult struct {
	XMLName xml.Name `xml:"CopyObjectResult"`
	// 返回文件的 MD5 算法校验值。ETag 的值可以用于检查 Object 的内容是否发生变化。
	ETag string `xml:"ETag,omitempty"`
	// 返回文件最后修改时间，GMT 格式
	LastModified string `xml:"LastModified,omitempty"`
}

// Copy ...
//
// Put Object Copy 请求实现将一个文件从源路径复制到目标路径。建议文件大小 1M 到 5G，
// 超过 5G 的文件请使用分块上传 Upload - Copy。在拷贝的过程中，文件元属性和 ACL 可以被修改。
//
// 用户可以通过该接口实现文件移动，文件重命名，修改文件属性和创建副本。
//
// 版本
//
// * 默认情况下，在目标存储桶上启用版本控制，对象存储会为正在复制的对象生成唯一的版本 ID。此版本 ID 与源对象的版本 ID 不同。
//   对象存储会在 x-cos-version-id 响应中的响应标头中返回复制对象的版本 ID。
// * 如果您在目标存储桶没有启用版本控制或暂停版本控制，则对象存储生成的版本 ID 始终为 null。
//
// 注意：在跨帐号复制的时候，需要先设置被复制文件的权限为公有读，或者对目标帐号赋权，同帐号则不需要。
//
// https://cloud.tencent.com/document/product/436/10881
func (s *ObjectService) Copy(ctx context.Context, name, sourceURL string, opt *ObjectCopyOptions) (*ObjectCopyResult, *Response, error) {
	var res ObjectCopyResult
	if opt == nil {
		opt = new(ObjectCopyOptions)
	}
	if opt.ObjectCopyHeaderOptions == nil {
		opt.ObjectCopyHeaderOptions = new(ObjectCopyHeaderOptions)
	}
	opt.XCosCopySource = sourceURL

	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name),
		method:    http.MethodPut,
		body:      nil,
		optHeader: opt,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// Delete Object 接口请求可以在 COS 的 Bucket 中将一个文件（Object）删除。该操作需要请求者对 Bucket 有 WRITE 权限。
//
// 细节分析
//
// * 在 DELETE Object 请求中删除一个不存在的 Object，仍然认为是成功的，返回 204 No Content。
// * DELETE Object 要求用户对该 Object 要有写权限。
//
// https://cloud.tencent.com/document/product/436/7743
func (s *ObjectService) Delete(ctx context.Context, name string) (*Response, error) {
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/" + encodeURIComponent(name),
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectHeadOptions ...
type ObjectHeadOptions struct {
	// 当 Object 在指定时间后被修改，则返回对应 Object 的 meta 信息，否则返回 304
	IfModifiedSince string `url:"-" header:"If-Modified-Since,omitempty"`
}

// Head Object请求可以取回对应Object的元数据，Head的权限与Get的权限一致
//
// 默认情况下，HEAD 操作从当前版本的对象中检索元数据。如要从不同版本检索元数据，请使用 versionId 子资源。
//
// https://cloud.tencent.com/document/product/436/7745
func (s *ObjectService) Head(ctx context.Context, name string, opt *ObjectHeadOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name),
		method:    http.MethodHead,
		optHeader: opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectOptionsOptions ...
//
// https://cloud.tencent.com/document/product/436/8288
type ObjectOptionsOptions struct {
	// 模拟跨域访问的请求来源域名，必选
	Origin string `url:"-" header:"Origin"`
	// 模拟跨域访问的请求 HTTP 方法，必选
	AccessControlRequestMethod string `url:"-" header:"Access-Control-Request-Method"`
	// 模拟跨域访问的请求头部
	AccessControlRequestHeaders string `url:"-" header:"Access-Control-Request-Headers,omitempty"`
}

// Options Object 接口实现 Object 跨域访问配置的预请求。
// 即在发送跨域请求之前会发送一个 OPTIONS 请求并带上特定的来源域，HTTP 方法和 Header 信息等给 COS，
// 以决定是否可以发送真正的跨域请求。当 CORS 配置不存在时，请求返回 403 Forbidden。
// 可以通过 PUT Bucket cors 接口来开启 Bucket 的 CORS 支持。
//
// https://cloud.tencent.com/document/product/436/8288
func (s *ObjectService) Options(ctx context.Context, name string, opt *ObjectOptionsOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name),
		method:    http.MethodOptions,
		optHeader: opt,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// Append ...
//
// Append请求可以将一个文件（Object）以分块追加的方式上传至 Bucket 中。使用Append Upload的文件必须事前被设定为Appendable。
// 当Appendable的文件被执行Put Object的操作以后，文件被覆盖，属性改变为Normal。
//
// 文件属性可以在Head Object操作中被查询到，当您发起Head Object请求时，会返回自定义Header『x-cos-object-type』，该Header只有两个枚举值：Normal或者Appendable。
//
// 追加上传建议文件大小1M - 5G。如果position的值和当前Object的长度不致，COS会返回409错误。
// 如果Append一个Normal的Object，COS会返回409 ObjectNotAppendable。
//
// Appendable的文件不可以被复制，不参与版本管理，不参与生命周期管理，不可跨区域复制。
//
// 当 r 不是 bytes.Buffer/bytes.Reader/strings.Reader 时，必须指定 opt.ObjectPutHeaderOptions.ContentLength
// 当 r 是个 io.ReadCloser 时 Append 方法不会自动调用 r.Close()，用户需要自行选择合适的时机去调用 r.Close() 方法对 r 进行资源回收
//
// https://www.qcloud.com/document/product/436/7741
func (s *ObjectService) Append(ctx context.Context, name string, position int, r io.Reader, opt *ObjectPutOptions) (*Response, error) {
	u := fmt.Sprintf("/%s?append&position=%d", encodeURIComponent(name), position)
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       u,
		method:    http.MethodPost,
		optHeader: opt,
		body:      r,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectDeleteMultiOptions ...
//
// https://cloud.tencent.com/document/product/436/8289
type ObjectDeleteMultiOptions struct {
	XMLName xml.Name `xml:"Delete" header:"-"`
	// 布尔值，这个值决定了是否启动 Quiet 模式。
	// 值为 true 启动 Quiet 模式，值为 false 则启动 Verbose 模式，默认值为 False
	Quiet bool `xml:"Quiet" header:"-"`
	// 说明每个将要删除的目标 Object 信息，只需指定 Key 即可
	Objects []Object `xml:"Object" header:"-"`
	//XCosSha1 string `xml:"-" header:"x-cos-sha1"`
}

// ObjectDeleteMultiResult ...
//
// https://cloud.tencent.com/document/product/436/8289
type ObjectDeleteMultiResult struct {
	XMLName xml.Name `xml:"DeleteResult"`
	// 说明本次删除的成功 Object 信息
	DeletedObjects []Object `xml:"Deleted,omitempty"`
	// 说明本次删除的失败 Object 信息
	Errors []struct {
		// 删除失败的 Object 的名称
		Key string
		// 删除失败的错误代码
		Code string
		// 删除失败的错误信息
		Message string
	} `xml:"Error,omitempty"`
}

// DeleteMulti ...
//
// Delete Multiple Object请求实现批量删除文件，最大支持单次删除1000个文件。
// 对于返回结果，COS提供Verbose和Quiet两种结果模式。Verbose模式将返回每个Object的删除结果；
// Quiet模式只返回报错的Object信息。
//
// 细节分析
//
// * 每一个批量删除请求，最多只能包含 1000个 需要删除的对象；
// * 批量删除支持二种模式的放回，verbose 模式和 quiet 模式，默认为 verbose 模式。
//   verbose 模式返回每个 key 的删除情况，quiet 模式只返回删除失败的 key 的情况；
// * 批量删除需要携带 Content-MD5 头部，用以校验请求 body 没有被修改；
// * 批量删除请求允许删除一个不存在的 key，仍然认为成功。
//
// https://cloud.tencent.com/document/product/436/8289
func (s *ObjectService) DeleteMulti(ctx context.Context, opt *ObjectDeleteMultiOptions) (*ObjectDeleteMultiResult, *Response, error) {
	var res ObjectDeleteMultiResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     "/?delete",
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

type objectPresignedURLTestingOptions struct {
	authTime *AuthTime
}

// PresignedURL 生成预签名授权 URL，可用于无需知道 SecretID 和 SecretKey 就可以上传和下载文件 。
//
// httpMethod:
//   * 下载文件：http.MethodGet
//   * 上传文件: http.MethodPut
//
// 下载文件 时 opt 可以是 *ObjectGetOptions ，上传文件时 opt 可以是 *ObjectPutOptions
//
// https://cloud.tencent.com/document/product/436/14116
// https://cloud.tencent.com/document/product/436/14114
func (s *ObjectService) PresignedURL(ctx context.Context, httpMethod, name string, auth Auth, opt interface{}) (*url.URL, error) {
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name),
		method:    httpMethod,
		optQuery:  opt,
		optHeader: opt,
	}
	req, err := s.client.newRequest(ctx, &sendOpt)
	if err != nil {
		return nil, err
	}

	var authTime *AuthTime
	if opt != nil {
		if opt, ok := opt.(*objectPresignedURLTestingOptions); ok {
			authTime = opt.authTime
		}
	}
	if authTime == nil {
		authTime = NewAuthTime(auth.Expire)
	}
	authorization := newAuthorization(auth, req, *authTime)
	sign := encodeURIComponent(authorization)

	if req.URL.RawQuery == "" {
		req.URL.RawQuery = fmt.Sprintf("sign=%s", sign)
	} else {
		req.URL.RawQuery = fmt.Sprintf("%s&sign=%s", req.URL.RawQuery, sign)
	}
	return req.URL, nil
}

// Object ...
type Object struct {
	// Object 的 Key
	Key string `xml:",omitempty"`
	// 文件的 MD-5 算法校验值
	ETag string `xml:",omitempty"`
	// 说明文件大小，单位是 Byte
	Size int `xml:",omitempty"`
	// 块编号
	PartNumber int `xml:",omitempty"`
	// 说明 Object 最后被修改时间
	LastModified string `xml:",omitempty"`
	// Object 的存储级别，枚举值：STANDARD，STANDARD_IA，ARCHIVE
	StorageClass string `xml:",omitempty"`
	// Bucket 持有者信息
	Owner *Owner `xml:",omitempty"`
}
