package bw6

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"testing"
)

var fuz int

func TestMain(m *testing.M) {
	_fuz := flag.Int("fuzz", 10, "# of iterations")
	flag.Parse()
	fuz = *_fuz
	os.Exit(m.Run())
}

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

func TestSome(t *testing.T) {

	fmt.Println(q.BitLen())
}
