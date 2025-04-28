package storage

import (
	"math/rand"
)

func randStr(length int) string {
	siqBytes := "abcdifghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = siqBytes[rand.Intn(len(siqBytes))]
	}
	url := string(b)
	return url
}

func getPostfix(storage UrlStorage, length int) string {
	for {
		postfix := randStr(length)
		if _, ok := storage[postfix]; !ok {
			return postfix
		}
	}
}
