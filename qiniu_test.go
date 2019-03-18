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
	objectQiniu = "qiniu.go"
	Qiniu       *QINIU
)

func init() {
	key := beego.AppConfig.String("qiniu::accessKey")
	secret := beego.AppConfig.String("qiniu::secretKey")
	bucket := beego.AppConfig.String("qiniu::bucket")
	domain := strings.ToLower(beego.AppConfig.String("qiniu::domain"))
	Qiniu, err = NewQINIU(key, secret, bucket, domain)
	if err != nil {
		panic(err)
	}
}


func TestQINIU(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = Qiniu.Upload(objectSVG, objectSVG,headerSVG)
	if err != nil {
		t.Error(err)
	}
	err = Qiniu.Upload(objectSVGGzip, objectSVGGzip, headerGzip, headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", Qiniu.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", Qiniu.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := Qiniu.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := Qiniu.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := Qiniu.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(Qiniu.GetSignURL(objectSVG, 1200))
	t.Log(Qiniu.GetSignURL(objectSVGGzip, 1200))
	t.Log("========Finished========")
}

func TestQINIU_Delete(t *testing.T) {
	if err := Qiniu.Delete(objectSVG, objectSVGGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := Qiniu.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}