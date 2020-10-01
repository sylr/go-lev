package rand

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func GetRandomHashSlice(count int) []string {
	hashes := make([]string, count)
	token := make([]byte, 4)

	for i := 0; i < count; i++ {
		rand.Read(token)
		hashes[i] = fmt.Sprintf("%x", md5.Sum(token))
	}

	return hashes
}
