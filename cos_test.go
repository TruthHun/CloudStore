package CloudStore

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/astaxie/beego"
)

var (
	Cos       *COS
	objectCOS = "cos.go"
)

func init() {
	var err error
	secretId := beego.AppConfig.String("cos::accessKey")
	secretKey := beego.AppConfig.String("cos::secretKey")
	bucket := beego.AppConfig.String("cos::bucket")
	appId := beego.AppConfig.String("cos::appId")
	region := beego.AppConfig.String("cos::region")
	domain := beego.AppConfig.String("cos::domain")
	Cos, err = NewCOS(secretId, secretKey, bucket, appId, region, domain)
	if err != nil {
		panic(err)
	}
}

func TestCOS(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = Cos.Upload(objectSVG, objectSVG,headerSVG)
	if err != nil {
		t.Error(err)
	}
	err = Cos.Upload(objectSVGGzip, objectSVGGzip, headerGzip,headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", Cos.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", Cos.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := Cos.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := Cos.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := Cos.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(Cos.GetSignURL(objectSVG, 120))
	t.Log(Cos.GetSignURL(objectSVGGzip, 120))
	t.Log("========Finished========")
}

func TestCOS_Delete(t *testing.T) {
	if err := Cos.Delete(objectSVG, objectSVGGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := Cos.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}
