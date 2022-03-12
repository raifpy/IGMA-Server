package goalchecker

import (
	"math/rand"
	"strconv"
)

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
