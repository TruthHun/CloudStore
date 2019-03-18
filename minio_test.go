package CloudStore

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/astaxie/beego"
)

var (
	Minio       *MinIO
	objectMinio = "minio.go"
)

func init() {
	key := beego.AppConfig.String("minio::accessKey")
	secret := beego.AppConfig.String("minio::secretKey")
	bucket := beego.AppConfig.String("minio::bucket")
	domain := strings.ToLower(beego.AppConfig.String("minio::domain"))
	endpoint := strings.ToLower(beego.AppConfig.String("minio::endpoint"))
	Minio, err = NewMinIO(key, secret, bucket, endpoint, domain)
	if err != nil {
		panic(err)
	}
}
func TestMinIO(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = Minio.Upload(objectSVG, objectSVG,headerSVG)
	if err != nil {
		t.Error(err)
	}
	err = Minio.Upload(objectSVGGzip, objectSVGGzip, headerGzip,headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", Minio.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", Minio.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := Minio.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := Minio.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := Minio.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(Minio.GetSignURL(objectSVG, 120))
	t.Log(Minio.GetSignURL(objectSVGGzip, 120))
	t.Log("========Finished========")
}

func TestMinIO_Delete(t *testing.T) {
	if err := Minio.Delete(objectSVG, objectSVGGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := Minio.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}