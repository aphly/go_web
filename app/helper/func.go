package helper

import (
	"math/rand"
	"time"
)

func Rand(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func RandStr(lenNum int, diff ...int64) string {
	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o",
		"p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y",
		"Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	str := ""
	length := len(chars)
	var r *rand.Rand
	if len(diff) == 0 {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano() + diff[0]))
	}
	for i := 0; i < lenNum; i++ {
		str += chars[r.Intn(length)]
	}
	return str
}
