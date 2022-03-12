package scores

import (
	"math/rand"
	"strings"
)

type ReplaceString string

func (r ReplaceString) Replace(m map[string]string) string {
	var text = string(r)
	for key, value := range m {
		text = strings.ReplaceAll(text, "{"+key+"}", value)
	}
	return text
}

type EasyList []string

func (e EasyList) Delete(key string) EasyList {
	var n = EasyList{}
	for _, a := range e {
		if a == key {
			continue
		}
		n = append(n, a)
	}

	return n
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
