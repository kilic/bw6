package bw6

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

type fe /***			***/ [N_LIMBS]uint64

func (e *fe) setBytes(in []byte) *fe {
	l := len(in)
	if l >= FE_BYTE_SIZE {
		l = FE_BYTE_SIZE
	}
	padded := make([]byte, FE_BYTE_SIZE)
	copy(padded[FE_BYTE_SIZE-l:], in[:])
	var a int
	for i := 0; i < N_LIMBS; i++ {
		a = FE_BYTE_SIZE - i*8
		e[i] = uint64(padded[a-1]) | uint64(padded[a-2])<<8 |
			uint64(padded[a-3])<<16 | uint64(padded[a-4])<<24 |
			uint64(padded[a-5])<<32 | uint64(padded[a-6])<<40 |
			uint64(padded[a-7])<<48 | uint64(padded[a-8])<<56
	}
	return e
}

func (e *fe) setBig(a *big.Int) *fe {
	return e.setBytes(a.Bytes())
}

func (e *fe) setString(s string) (*fe, error) {
	if s[:2] == "0x" {
		s = s[2:]
	}
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return e.setBytes(bytes), nil
}

func (e *fe) set(e2 *fe) *fe {
	e[0] = e2[0]
	e[1] = e2[1]
	e[2] = e2[2]
	e[3] = e2[3]
	e[4] = e2[4]
	e[5] = e2[5]
	e[6] = e2[6]
	e[7] = e2[7]
	e[8] = e2[8]
	e[9] = e2[9]
	e[10] = e2[10]
	e[11] = e2[11]
	return e
}

func (e *fe) bytes() []byte {
	out := make([]byte, FE_BYTE_SIZE)

	var a int
	for i := 0; i < N_LIMBS; i++ {
		a = FE_BYTE_SIZE - i*8
		out[a-1] = byte(e[i])
		out[a-2] = byte(e[i] >> 8)
		out[a-3] = byte(e[i] >> 16)
		out[a-4] = byte(e[i] >> 24)
		out[a-5] = byte(e[i] >> 32)
		out[a-6] = byte(e[i] >> 40)
		out[a-7] = byte(e[i] >> 48)
		out[a-8] = byte(e[i] >> 56)
	}
	return out
}

func (e *fe) big() *big.Int {
	return new(big.Int).SetBytes(e.bytes())
}

func (e *fe) string() (s string) {
	for i := N_LIMBS - 1; i >= 0; i-- {
		s = fmt.Sprintf("%s%16.16x", s, e[i])
	}
	return "0x" + s
}

func (e *fe) zero() *fe {
	e[0] = 0
	e[1] = 0
	e[2] = 0
	e[3] = 0
	e[4] = 0
	e[5] = 0
	e[6] = 0
	e[7] = 0
	e[8] = 0
	e[9] = 0
	e[10] = 0
	e[11] = 0
	return e
}

func (e *fe) one() *fe {
	return e.set(r1)
}

func (e *fe) rand(r io.Reader) (*fe, error) {
	bi, err := rand.Int(r, modulus.big())
	if err != nil {
		return nil, err
	}
	return e.setBig(bi), nil
}

func (e *fe) isValid() bool {
	return e.cmp(&modulus) == -1
}

func (e *fe) isOdd() bool {
	var mask uint64 = 1
	return e[0]&mask != 0
}

func (e *fe) isEven() bool {
	var mask uint64 = 1
	return e[0]&mask == 0
}

func (e *fe) isZero() bool {
	return (e[11] | e[10] | e[9] | e[8] | e[7] | e[6] | e[5] | e[4] | e[3] | e[2] | e[1] | e[0]) == 0
}

func (e *fe) isOne() bool {
	return e.equal(r1)
}

func (e *fe) cmp(e2 *fe) int {
	for i := N_LIMBS - 1; i >= 0; i-- {
		if e[i] > e2[i] {
			return 1
		} else if e[i] < e2[i] {
			return -1
		}
	}
	return 0
}

func (e *fe) equal(e2 *fe) bool {
	return e2[0] == e[0] && e2[1] == e[1] && e2[2] == e[2] && e2[3] == e[3] && e2[4] == e[4] && e2[5] == e[5] && e2[6] == e[6] && e2[7] == e[7] && e2[8] == e[8] && e2[9] == e[9] && e2[10] == e[10] && e2[11] == e[11]
}

func (e *fe) signBE() bool {
	negZ, z := new(fe), new(fe)
	fromMont(z, e)
	neg(negZ, z)
	return negZ.cmp(z) > -1
}

func (e *fe) sign() bool {
	r := new(fe)
	fromMont(r, e)
	return r[0]&1 == 0
}

func (e *fe) div2(u uint64) {
	e[0] = e[0]>>1 | e[1]<<63
	e[1] = e[1]>>1 | e[2]<<63
	e[2] = e[2]>>1 | e[3]<<63
	e[3] = e[3]>>1 | e[4]<<63
	e[4] = e[4]>>1 | e[5]<<63
	e[5] = e[5]>>1 | e[6]<<63
	e[6] = e[6]>>1 | e[7]<<63
	e[7] = e[7]>>1 | e[8]<<63
	e[8] = e[8]>>1 | e[9]<<63
	e[9] = e[9]>>1 | e[10]<<63
	e[10] = e[10]>>1 | e[11]<<63
	e[11] = e[11]>>1 | u<<63
}

func (e *fe) mul2() uint64 {
	u := e[11] >> 63
	e[11] = e[11]<<1 | e[10]>>63
	e[10] = e[10]<<1 | e[9]>>63
	e[9] = e[9]<<1 | e[8]>>63
	e[8] = e[8]<<1 | e[7]>>63
	e[7] = e[7]<<1 | e[6]>>63
	e[6] = e[6]<<1 | e[5]>>63
	e[5] = e[5]<<1 | e[4]>>63
	e[4] = e[4]<<1 | e[3]>>63
	e[3] = e[3]<<1 | e[2]>>63
	e[2] = e[2]<<1 | e[1]>>63
	e[1] = e[1]<<1 | e[0]>>63
	e[0] = e[0] << 1
	return u
}
