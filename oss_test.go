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
	O         *OSS
	objectOSS = "oss.go"
)

func init() {
	var err error

	key := beego.AppConfig.String("oss::accessKey")
	secret := beego.AppConfig.String("oss::secretKey")
	endpoint := beego.AppConfig.String("oss::endpoint")
	bucket := beego.AppConfig.String("oss::bucket")
	domain := strings.ToLower(beego.AppConfig.String("oss::domain"))

	O, err = NewOSS(key, secret, endpoint, bucket, domain)
	if err != nil {
		panic(err)
	}
}



func TestOSS(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = O.Upload(objectSVG, objectSVG,headerSVG)
	if err != nil {
		t.Error(err)
	}
	err = O.Upload(objectSVGGzip, objectSVGGzip, headerGzip, headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", O.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", O.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := O.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := O.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := O.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(O.GetSignURL(objectSVG, 1200))
	t.Log(O.GetSignURL(objectSVGGzip, 1200))
	t.Log("========Finished========")
}

func TestOSS_Delete(t *testing.T) {
	if err := O.Delete(objectSVG, objectSVGGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := O.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}
