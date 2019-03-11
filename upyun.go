package CloudStore

type UpYun struct {
	Bucket   string
	Operator string
	Password string
	Domain   string
}

func NewUpYun(bucket, operator, password, domain string) *UpYun {
	return &UpYun{
		Bucket:   bucket,
		Operator: operator,
		Password: password,
		Domain:   domain,
	}
}

func (u *UpYun) IsExist(object string) (err error) {

	return
}

func (u *UpYun) Put(tmpFile, saveFile string, header ...map[string]string) (err error) {
	return
}

func (u *UpYun) Delete(object ...string) (err error) {
	return
}

func (u *UpYun) GetSignURL(object string, expire int64) (link string, err error) {
	return
}
