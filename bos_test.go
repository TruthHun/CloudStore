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
	objectBOS = "bos.go"
	Bos       *BOS
	err       error
)

func init() {

	key := beego.AppConfig.String("bos::accessKey")
	secret := beego.AppConfig.String("bos::secretKey")
	endpoint := beego.AppConfig.String("bos::endpoint")
	bucket := beego.AppConfig.String("bos::bucket")
	domain := strings.ToLower(beego.AppConfig.String("bos::domain"))

	Bos, err = NewBOS(key, secret, bucket, endpoint, domain)
	if err != nil {
		panic(err)
	}
}

func TestBOS(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = Bos.Upload(objectSVG, objectSVG, headerSVG)
	if err != nil {
		t.Error(err)
	}
	err = Bos.Upload(objectSVGGzip, objectSVGGzip, headerGzip, headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", Bos.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", Bos.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := Bos.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		if len(files) == 0 {
			t.Error("获取列表数据失败")
		}
		t.Log(fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := Bos.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := Bos.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(Bos.GetSignURL(objectSVG, 120))
	t.Log(Bos.GetSignURL(objectSVGGzip, 120))
	t.Log("========Finished========")
}

func TestBOS_GetSignURL(t *testing.T) {
	t.Log(Bos.GetSignURL(objectSVG, 3600))
}

func TestBOS_Delete(t *testing.T) {
	if err := Bos.Delete(objectSVG, objectSVGGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := Bos.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}
