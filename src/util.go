package soccer

import (
	"math/rand"
	"os"
	"strconv"
)

/*
func stringListContains(list []string, key string) bool {
	for _, a := range list {
		if a == key {
			return true
		}
	}

	return false
}*/

func createTempVideoDir(path string) (v string) {
	v = "./tempmedia/" + path
	os.MkdirAll(v, os.ModePerm)
	return v
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func IntToFfmpegInt(s int) string {
	if s <= 0 {
		return "00"
	} else if s < 10 {
		return "0" + strconv.Itoa(s)
	}
	return strconv.Itoa(s)
}
