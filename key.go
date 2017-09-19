package main

import (
	"math/rand"
	"time"
)

func keyGen(n int) string {
	// Entropy of 26**n...
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := make([]byte, n);

	for i := 0; i < n; i++ {
		// Generate a random letter in a-z
		var offset int
		if (r.Intn(2) == 1) {
			offset = 97
		} else {
			offset = 65
		}
		ret[i] = (byte)(r.Intn(26) + offset)
	}
	realret := string(ret[:])
	
	return realret
}
