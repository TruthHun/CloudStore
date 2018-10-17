package CloudStore

type CloudStore interface {
	PutObject(string, string, map[string]string) error
	DeleteObjects([]string) error
	GetObjectURL(string, int64) (string, error)
}
