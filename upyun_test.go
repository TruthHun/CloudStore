package CloudStore

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/astaxie/beego"
)

var Up *UpYun

func init() {
	bucket := beego.AppConfig.String("upyun::bucket")
	operator := beego.AppConfig.String("upyun::operator")
	password := beego.AppConfig.String("upyun::password")
	domain := strings.ToLower(beego.AppConfig.String("upyun::domain"))
	secret := strings.ToLower(beego.AppConfig.String("upyun::secret"))
	Up = NewUpYun(bucket, operator, password, domain, secret)
}

func TestUpYun(t *testing.T) {
	// upload
	t.Log("=====Upload=====", objectSVG, objectSVGGzip)
	err = Up.Upload(objectSVG, objectSVG, headerSVG)
	if err != nil {
		t.Error(err)
	}
	err = Up.Upload(objectSVGGzip, objectSVGGzip, headerGzip, headerSVG)
	if err != nil {
		t.Error(err)
	}
	t.Log("=====IsExist=====")
	t.Log(objectSVG, "is exist?(Y):", Up.IsExist(objectSVG) == nil)
	t.Log(objectNotExist, "is exist?(N):", Up.IsExist(objectNotExist) == nil)
	t.Log("=====Lists=====")
	if files, err := Up.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(objectPrefix, ":", fmt.Sprintf("%+v", files))
	}
	t.Log("=====GetInfo=====")
	if info, err := Up.GetInfo(objectSVG); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(fmt.Sprintf("%+v", info))
	}
	t.Log("=====Download=====")
	if err := Up.Download(objectSVG, objectDownload); err != nil {
		t.Error(err)
	} else {
		t.Log("download success")
		b, _ := ioutil.ReadFile(objectDownload)
		t.Log("Content:", string(b))
		os.Remove(objectDownload)
	}
	t.Log("====GetSignURL====")
	t.Log(Up.GetSignURL(objectSVG, 1200))
	t.Log(Up.GetSignURL(objectSVGGzip, 1200))
	t.Log("========Finished========")
}

func TestUpYun_Delete(t *testing.T) {
	if err := Up.Delete(objectSVG, objectSVGGzip); err != nil {
		t.Error(err)
	} else {
		t.Log("delete success")
	}

	if files, err := Up.Lists(objectPrefix); err != nil {
		t.Error(err)
	} else {
		t.Log(fmt.Sprintf("%+v", files))
	}
}
