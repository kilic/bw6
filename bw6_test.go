package bw6

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func padBytes(in []byte, size int) []byte {
	out := make([]byte, size)
	if len(in) > size {
		panic("bad input for padding")
	}
	copy(out[size-len(in):], in)
	return out
}

func randScalar(max *big.Int) *big.Int {
	a, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(errors.New(""))
	}
	return a
}
