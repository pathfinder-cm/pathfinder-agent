package util

import (
	"math/rand"
	"time"
)

func RandomIntRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn((max-min)+1) + min
}

func RandomString(len int) string {
	bytes := make([]byte, len)

	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}

	return string(bytes)
}
