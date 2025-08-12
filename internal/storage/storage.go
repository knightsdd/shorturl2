package storage

type InMemoryStorage map[string]string

type Storage interface {
	SaveValue(value string) (key string)
	GetValue(key string) (value string, ok bool)
}

const (
	postfixLength = 8
)

var storage InMemoryStorage

func init() {
	storage = make(InMemoryStorage, 10)
}

func GetStorage() InMemoryStorage {
	return storage
}

func (s InMemoryStorage) SaveValue(value string) (key string) {
	key = getPostfix(s, postfixLength)
	s[key] = value
	return
}

func (s InMemoryStorage) GetValue(key string) (value string, ok bool) {
	value, ok = s[key]
	return
}
