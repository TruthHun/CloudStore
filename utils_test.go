package CloudStore

import "testing"

var (
	objectSVG      = "test_data/test.svg"             //未经过gzip压缩的svg图片
	objectSVGGzip  = "test_data/test.gzip.svg"        //gzip压缩后的svg图片
	objectHtml     = "test_data/helloworld.html"      //未经gzip压缩的HTML
	objectHtmlGzip = "test_data/helloworld.gzip.html" //gzip压缩后的HTML
	objectNotExist = "not exist object"
	objectPrefix   = "test_data"
	objectDownload = "test_data/download.svg"
	headerGzip     = map[string]string{"Content-Encoding": "gzip"}
	headerSVG      = map[string]string{"Content-Type": "image/svg+xml"}
	headerHtml     = map[string]string{"Content-Type": "text/html; charset=UTF-8"}
)

func TestCompressByGzip(t *testing.T) {
	err = CompressByGzip(objectSVG, objectSVGGzip)
	if err != nil {
		t.Error(err)
	}
	err = CompressByGzip(objectHtml, objectHtmlGzip)
	if err != nil {
		t.Error(err)
	}
}
