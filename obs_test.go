package CloudStore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"obs"
	"os"
	"testing"

	"github.com/astaxie/beego"
)

var Obs *OBS

func init() {
	accessKey := beego.AppConfig.String("obs::accessKey")
	secretKey := beego.AppConfig.String("obs::secretKey")
	bucket := beego.AppConfig.String("obs::bucket")
	endpoint := beego.AppConfig.String("obs::endpoint")
	domain := beego.AppConfig.String("obs::domain")
	Obs, err = NewOBS(accessKey, secretKey, bucket, endpoint, domain)
	if err != nil {
		panic(err)
	}
}

func TestOBS(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = Obs.Upload(objectSVG, objectSVG)
	if err != nil {
		t.Error(err)
	}
	err = Obs.Upload(objectSVGGzip, objectSVGGzip, headerGzip, headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", Obs.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", Obs.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := Obs.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := Obs.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := Obs.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(Obs.GetSignURL(objectSVG, 1200))
	t.Log(Obs.GetSignURL(objectSVGGzip, 1200))
	t.Log("========Finished========")
}

func TestOBS_Upload(t *testing.T) {
	input := &obs.PutObjectInput{}
	input.Bucket = Obs.Bucket
	input.Key = "official.html"
	b, _ := ioutil.ReadFile(objectHtml)
	input.Body = bytes.NewBuffer(b)

	_, err := Obs.Client.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create object:%s successfully!\n", input.Key)
	fmt.Println()
	fmt.Println(Obs.GetSignURL(input.Key, 3600))
}

func TestOBS_Delete(t *testing.T) {
	if err := Obs.Delete(objectSVG, objectSVGGzip, objectHtml, objectHtmlGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := Obs.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}
