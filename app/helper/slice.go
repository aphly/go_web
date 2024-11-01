package helper

import (
	"fmt"
	"math/rand"
	"time"
)

func SliceRand(slice []any) any {
	if len(slice) == 0 {
		panic("slice is empty")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randIndex := r.Intn(len(slice))
	return slice[randIndex]
}

func InArray(needle any, haystack []any) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

func SliceJoin(anySlice []any, sep string) string {
	var str = ""
	var count = len(anySlice)
	for i, v := range anySlice {
		if s, ok := v.(string); ok {
			if i == count-1 {
				str += s
			} else {
				str += s + sep
			}
		} else {
			if i == count-1 {
				str += fmt.Sprintf("%v", v)
			} else {
				str += fmt.Sprintf("%v", v) + sep
			}
		}
	}
	return str
}
