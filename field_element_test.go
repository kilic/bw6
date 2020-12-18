package bw6

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestFieldElementValidation(t *testing.T) {
	// fe
	zero := new(fe).zero()
	if !zero.isValid() {
		t.Fatal("zero must be valid")
	}
	one := new(fe).one()
	if !one.isValid() {
		t.Fatal("one must be valid")
	}
	if modulus.isValid() {
		t.Fatal("modulus must be invalid")
	}
	n := modulus.big()
	n.Add(n, big.NewInt(1))
	if new(fe).setBig(n).isValid() {
		t.Fatal("number greater than modulus must be invalid")
	}
}

func TestFieldElementEquality(t *testing.T) {
	// fe
	zero := new(fe).zero()
	if !zero.equal(zero) {
		t.Fatal("0 == 0")
	}
	one := new(fe).one()
	if !one.equal(one) {
		t.Fatal("1 == 1")
	}
	a, _ := new(fe).rand(rand.Reader)
	if !a.equal(a) {
		t.Fatal("a == a")
	}
	b := new(fe)
	add(b, a, one)
	if a.equal(b) {
		t.Fatal("a != a + 1")
	}
}

func TestFieldElementHelpers(t *testing.T) {
	// fe
	zero := new(fe).zero()
	if !zero.isZero() {
		t.Fatal("'zero' is not zero")
	}
	one := new(fe).one()
	if !one.isOne() {
		t.Fatal("'one' is not one")
	}
	odd := new(fe).setBig(big.NewInt(1))
	if !odd.isOdd() {
		t.Fatal("1 must be odd")
	}
	if odd.isEven() {
		t.Fatal("1 must not be even")
	}
	even := new(fe).setBig(big.NewInt(2))
	if !even.isEven() {
		t.Fatal("2 must be even")
	}
	if even.isOdd() {
		t.Fatal("2 must not be odd")
	}
}

func TestFieldElementSerialization(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		in := make([]byte, fpByteSize)
		fe := new(fe).setBytes(in)
		if !fe.isZero() {
			t.Fatal("bad serialization")
		}
		if !bytes.Equal(in, fe.bytes()) {
			t.Fatal("bad serialization")
		}
	})
	t.Run("bytes", func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			a, _ := new(fe).rand(rand.Reader)
			b := new(fe).setBytes(a.bytes())
			if !a.equal(b) {
				t.Fatal("bad serialization")
			}
		}
	})
	t.Run("big", func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			a, _ := new(fe).rand(rand.Reader)
			b := new(fe).setBig(a.big())
			if !a.equal(b) {
				t.Fatal("bad encoding or decoding")
			}
		}
	})
	t.Run("string", func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			a, _ := new(fe).rand(rand.Reader)
			b, err := new(fe).setString(a.string())
			if err != nil {
				t.Fatal(err)
			}
			if !a.equal(b) {
				t.Fatal("bad encoding or decoding")
			}
		}
	})
}

func TestFieldElementByteInputs(t *testing.T) {
	zero := new(fe).zero()
	in := make([]byte, 0)
	a := new(fe).setBytes(in)
	if !a.equal(zero) {
		t.Fatal("bad serialization")
	}
	in = make([]byte, fpByteSize)
	a = new(fe).setBytes(in)
	if !a.equal(zero) {
		t.Fatal("bad serialization")
	}
	in = make([]byte, fpByteSize+200)
	a = new(fe).setBytes(in)
	if !a.equal(zero) {
		t.Fatal("bad serialization")
	}
	in = make([]byte, fpByteSize+1)
	in[fpByteSize-1] = 1
	normalOne := &fe{1}
	a = new(fe).setBytes(in)
	if !a.equal(normalOne) {
		t.Fatal("bad serialization")
	}
}

func TestFieldElementCopy(t *testing.T) {
	a, _ := new(fe).rand(rand.Reader)
	b := new(fe).set(a)
	if !a.equal(b) {
		t.Fatal("bad copy, 1")
	}
}
