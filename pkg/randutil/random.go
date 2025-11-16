package randutil

import (
	"math/rand"
	"time"
)


func NewSourceRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func PickRandom(r *rand.Rand, items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[r.Intn(len(items))]
}
