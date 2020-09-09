package bw6

import (
	"errors"
	"math/big"
)

type fp6Temp struct {
	t [4]*fe3
}

type fp6 struct {
	fp3 *fp3
	fp6Temp
}

func newFp6Temp() fp6Temp {
	t := [4]*fe3{}
	for i := 0; i < len(t); i++ {
		t[i] = &fe3{}
	}
	return fp6Temp{t}
}

func newFp6(f *fp3) *fp6 {
	t := newFp6Temp()
	if f == nil {
		return &fp6{newFp3(), t}
	}
	return &fp6{f, t}
}

func (e *fp6) fromBytes(in []byte) (*fe6, error) {
	if len(in) != 6*FE_BYTE_SIZE {
		return nil, errors.New("input string should be larger than 96 bytes")
	}
	fp3 := e.fp3
	c0, err := fp3.fromBytes(in[:3*FE_BYTE_SIZE])
	if err != nil {
		return nil, err
	}
	c1, err := fp3.fromBytes(in[3*FE_BYTE_SIZE:])
	if err != nil {
		return nil, err
	}
	return &fe6{*c0, *c1}, nil
}

func (e *fp6) toBytes(a *fe6) []byte {
	out := make([]byte, 6*FE_BYTE_SIZE)
	fp3 := e.fp3
	copy(out[:3*FE_BYTE_SIZE], fp3.toBytes(&a[0]))
	copy(out[3*FE_BYTE_SIZE:], fp3.toBytes(&a[1]))
	return out
}

func (e *fp6) new() *fe6 {
	return new(fe6).zero()
}

func (e *fp6) zero() *fe6 {
	return new(fe6).zero()
}

func (e *fp6) one() *fe6 {
	return new(fe6).one()
}

func (e *fp6) add(c, a, b *fe6) {
	// c0 = a0 + b0
	// c1 = a1 + b1
	fp3 := e.fp3
	fp3.add(&c[0], &a[0], &b[0])
	fp3.add(&c[1], &a[1], &b[1])
}

func (e *fp6) ladd(c, a, b *fe6) {
	// c0 = a0 + b0
	// c1 = a1 + b1
	fp3 := e.fp3
	fp3.ladd(&c[0], &a[0], &b[0])
	fp3.ladd(&c[1], &a[1], &b[1])
}

func (e *fp6) double(c, a *fe6) {
	// c0 = 2a0
	// c1 = 2a1
	fp3 := e.fp3
	fp3.double(&c[0], &a[0])
	fp3.double(&c[1], &a[1])
}

func (e *fp6) ldouble(c, a *fe6) {
	// c0 = 2a0
	// c1 = 2a1
	fp3 := e.fp3
	fp3.ldouble(&c[0], &a[0])
	fp3.ldouble(&c[1], &a[1])
}

func (e *fp6) sub(c, a, b *fe6) {
	// c0 = a0 - b0
	// c1 = a1 - b1
	fp3 := e.fp3
	fp3.sub(&c[0], &a[0], &b[0])
	fp3.sub(&c[1], &a[1], &b[1])
}

func (e *fp6) neg(c, a *fe6) {
	// c0 = -a0
	// c1 = -a1
	fp3 := e.fp3
	fp3.neg(&c[0], &a[0])
	fp3.neg(&c[1], &a[1])
}

func (e *fp6) conjugate(c, a *fe6) {
	// c0 = a0
	// c1 = -a1
	fp3 := e.fp3
	c.set(a)
	c[0].set(&a[0])
	fp3.neg(&c[1], &a[1])
}

func (e *fp6) mul(c, a, b *fe6) {
	// Multiplication and Squaring on Pairing-Friendly Fields
	// Karatsuba multiplication algorithm
	// https://eprint.iacr.org/2006/471
	//
	// v0 = a0b0
	// c0 = v0 + αv1 = v0 - 4v1
	// c1 = (a0 + a1)(b0 + b1) - v0 - v1
	fp3, t := e.fp3, e.t

	fp3.mul(t[1], &a[0], &b[0]) // v0 = a0b0
	fp3.mul(t[2], &a[1], &b[1]) // v1 = a1b1

	fp3.ladd(t[0], &a[0], &a[1]) // a0 + a1
	fp3.ladd(t[3], &b[0], &b[1]) // b0 + b1
	fp3.mul(t[0], t[0], t[3])    // (a0 + a1)(b0 + b1)
	fp3.sub(t[0], t[0], t[1])    // (a0 + a1)(b0 + b1) - v0
	fp3.sub(&c[1], t[0], t[2])   // c1 = (a0 + a1)(b0 + b1) - v0 - v1

	fp3.double(t[2], t[2])     //
	fp3.double(t[2], t[2])     // -4v1
	fp3.sub(&c[0], t[1], t[2]) // c0 = v0 - 4v1
}

func (e *fp6) square(c, a *fe6) {
	e.squareComplex(c, a)
}

func (e *fp6) squareKaratsuba(c, a *fe6) {
	// Multiplication and Squaring on Pairing-Friendly Fields
	// Karatsuba squaring algorithm
	// https://eprint.iacr.org/2006/471
	//
	// v0 = a0^2
	// v1 = a1^2
	// c0 = v0 + αv1 = v0 - 4v1
	// c1 = (a0 + a1)^2 - v0 - v1

	fp3, t := e.fp3, e.t
	fp3.square(t[0], &a[0]) // v0 = a0^2
	fp3.square(t[1], &a[1]) // v1 = a1^2

	fp3.double(t[2], t[1]) //
	fp3.double(t[2], t[2]) // 4v1

	fp3.sub(t[3], t[0], t[2]) // c0 = v0 - 4v1

	fp3.ladd(t[2], &a[0], &a[1]) // a0 + a1
	fp3.square(t[2], t[2])       // (a0 + a1)^2

	fp3.sub(t[2], t[2], t[0])  // (a0 + a1)^2 - v0
	fp3.sub(&c[1], t[2], t[1]) // c1 = (a0 + a1)^2 - v0 - v1

	c[0].set(t[3])

}

func (e *fp6) squareComplex(c, a *fe6) {
	// Multiplication and Squaring on Pairing-Friendly Fields
	// Complex squaring algorithm
	// https://eprint.iacr.org/2006/471
	//
	// v0 = a0a1
	// c0 = (a0 + a1)(a0 + αa1) - v0 - αv0 = (a0 + a1)(a0 - 4a1) + 3a1a0
	// c1 = 2v0

	fp3, t := e.fp3, e.t
	fp3.double(t[0], &a[1])      // 2a1
	fp3.double(t[0], t[0])       // 4a1
	fp3.mul(t[1], &a[0], &a[1])  // a0a1
	fp3.double(t[2], t[1])       // 2a0a1
	fp3.ladd(t[3], &a[0], &a[1]) // a0 + a1
	c[1].set(t[2])               // c1 = 2a0a1
	fp3.add(t[2], t[2], t[1])    // 3a0a1
	fp3.sub(t[0], &a[0], t[0])   // (a0 - 4a1)
	fp3.mul(t[0], t[0], t[3])    // (a0 + a1)(a0 - 4a1)
	fp3.add(&c[0], t[2], t[0])   // (a0 + a1)(a0 - 4a1) + 3a0a1
}

func (e *fp6) inverse(c, a *fe6) {
	// Guide to Pairing Based Cryptography
	// Algorithm 5.19
	//
	// v = (a0^2 - βa1^2)^-1 =  (a0^2 + 4a1^2)^-1
	// c0 = a0v
	// c1 = a1v
	fp3, t := e.fp3, e.t

	fp3.square(t[0], &a[0]) // a0^2
	fp3.square(t[1], &a[1]) // a1^2

	fp3.double(t[2], t[1])
	fp3.double(t[2], t[2]) // 4a1^2

	fp3.add(t[0], t[0], t[2]) //
	fp3.inverse(t[1], t[0])   // v = a0^2 + 4a1^2

	fp3.mul(&c[0], t[1], &a[0]) // a0v1
	fp3.mul(t[1], t[1], &a[1])  // a1v1
	fp3.neg(&c[1], t[1])
}

func (e *fp6) exp(c, a *fe6, s *big.Int) {
	z := e.one()
	for i := s.BitLen() - 1; i >= 0; i-- {
		e.square(z, z)
		if s.Bit(i) == 1 {
			e.mul(z, z, a)
		}
	}
	c.set(z)
}

func (e *fp6) frobeniusMap(c, a *fe6, power int) {
	fp3 := e.fp3
	fp3.frobeniusMap(&c[0], &a[0], power)
	mul(&c[1][0], &a[1][0], &frobeniusCoeffs6[power%6][0])
	mul(&c[1][1], &a[1][1], &frobeniusCoeffs6[power%6][1])
	mul(&c[1][2], &a[1][2], &frobeniusCoeffs6[power%6][2])
}

func (e *fp6) frobeniusMap1(c, a *fe6, power int) {
	fp3 := e.fp3
	fp3.frobeniusMap(&c[0], &a[0], power)
	neg(&c[1][0], &a[1][0])
	mul(&c[1][1], &a[1][1], &frobeniusCoeffs6[1][1])
	mul(&c[1][2], &a[1][2], &frobeniusCoeffs6[1][2])
}
