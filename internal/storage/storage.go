package storage

type UrlStorage map[string]string

const (
	postfixLength = 8
)

var storage UrlStorage

func init() {
	storage = make(UrlStorage, 10)
}

func GetStorage() UrlStorage {
	return storage
}

func (s UrlStorage) SaveValue(value string) (key string) {
	key = getPostfix(s, postfixLength)
	s[key] = value
	return
}

func (s UrlStorage) GetValue(key string) (value string, ok bool) {
	value, ok = s[key]
	return
}
