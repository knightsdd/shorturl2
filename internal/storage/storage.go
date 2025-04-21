package storage

type UrlStorage map[string]string

var storage UrlStorage

func init() {
	storage = make(UrlStorage, 10)
}

func GetStorage() UrlStorage {
	return storage
}

func (s UrlStorage) SaveValue(key, value string) {
	s[key] = value
}

func (s UrlStorage) GetValue(key string) (value string, ok bool) {
	value, ok = s[key]
	return
}
