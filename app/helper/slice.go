package helper

import (
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
