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

// TODO: zexe uses different name which is mul by 034
// lets say {a,b} \in Fp6 where
// in sparse multiplication of a*b
// a = (a00 + a01*u + a02*u^2) + (a10 + a11*u + a12* u^2)
// b = (b00 + 0 + b02*u^2) + (0 + b11*u + 0)
func (e *fp6) mulBy024Assign(a *fe6, c0, c1, c2 *fe) {
	z0, z1, z2, z3, z4, z5 := &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}
	z0.set(&a[0][0])
	z1.set(&a[0][1])
	z2.set(&a[0][2])
	z3.set(&a[1][0])
	z4.set(&a[1][1])
	z5.set(&a[1][2])

	x0, x2, x4 := &fe{}, &fe{}, &fe{}
	x0.set(c0)
	x2.set(c1)
	x4.set(c2)

	d0, d2, d4 := &fe{}, &fe{}, &fe{}
	s0, s1 := &fe{}, &fe{}
	t0, t1, t2, t3, t4 := &fe{}, &fe{}, &fe{}, &fe{}, &fe{}

	mul(d0, z0, x0)
	mul(d2, z2, x2)
	mul(d4, z4, x4)
	add(t2, z0, z4)
	add(t1, z0, z2)
	add(s0, z1, z3)
	addAssign(s0, z5)

	// z0
	mul(s1, z1, x2)
	add(t3, s1, d4)
	mul(t4, t3, nonResidue1) // TODO: check
	addAssign(t4, d0)
	z0.set(t4)

	// z1
	mul(t3, z5, x4)
	add(s1, s1, t3)
	add(t3, t3, d2)
	mul(t4, t3, nonResidue1)
	mul(t3, z1, x0)
	addAssign(s1, t3)
	addAssign(t4, t3)
	z1.set(t4)

	// z2
	add(t0, x0, x2)
	mul(t3, t1, t0)
	subAssign(t3, d0)
	subAssign(t3, d2)
	mul(t4, z3, x4)
	addAssign(s1, t4)
	addAssign(t3, t4)

	// z3
	add(t0, z2, z4)
	z2.set(t3)
	add(t1, x2, x4)
	mul(t3, t0, t1)
	subAssign(t3, d2)
	subAssign(t3, d4)
	mul(t4, t3, nonResidue1)
	mul(t3, z3, x0)
	addAssign(s1, t3)
	addAssign(t4, t3)
	z3.set(t4)

	// z4
	mul(t3, z5, x2)
	addAssign(s1, t3)
	mul(t4, t3, nonResidue1)
	add(t0, x0, x4)
	mul(t3, t2, t0)
	subAssign(t3, d0)
	subAssign(t3, d4)
	addAssign(t4, t3)
	z4.set(t4)

	// z5
	add(t0, x0, x2)
	addAssign(t0, x4)
	add(t3, s0, t0)
	subAssign(t3, s1)
	z5.set(t3)

	// result
	a[0][0].set(z0)
	a[0][1].set(z1)
	a[0][2].set(z2)
	a[1][0].set(z3)
	a[1][1].set(z4)
	a[1][2].set(z5)
}

// TODO: zexe uses different name which is mul by 034
// lets say {a,b} \in Fp6 where
// in sparse multiplication of a*b
// a = (a00 + a01*u + a02*u^2) + (a10 + a11*u + a12*u^2)
// b = (b00 + 0 + 0) + (0 + b11*u + b12*u^2)
func (e *fp6) mulBy045Assign(a *fe6, c0, c1, c2 *fe) {
	z0, z1, z2, z3, z4, z5 := &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}
	z0.set(&a[0][0])
	z1.set(&a[0][1])
	z2.set(&a[0][2])
	z3.set(&a[1][0])
	z4.set(&a[1][1])
	z5.set(&a[1][2])

	x0, x4, x5 := &fe{}, &fe{}, &fe{}
	x0.set(c0)
	x4.set(c1)
	x5.set(c2)

	t0, t1, t2, t3, t4, t5, t6, t7, t8 := &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}

	mul(t6, x4, nonResidue1)
	mul(t7, x5, nonResidue1)

	mul(t0, x0, z0)
	mul(t8, t6, z4)
	addAssign(t0, t8)
	mul(t8, t7, z3)
	addAssign(t0, t8)

	mul(t1, x0, z1)
	mul(t8, t6, z5)
	addAssign(t1, t8)
	mul(t8, t7, z4)
	addAssign(t1, t8)

	mul(t2, x0, z2)
	mul(t8, x4, z3)
	addAssign(t2, t8)
	mul(t8, t7, z1)

	mul(t3, x0, z3)
	mul(t8, t6, z2)
	addAssign(t3, t8)
	mul(t8, t7, z1)
	addAssign(t3, t8)

	mul(t4, x0, z4)
	mul(t8, x4, z0)
	addAssign(t4, t8)
	mul(t8, t7, z2)
	addAssign(t4, t8)

	mul(t5, x0, z5)
	mul(t8, x4, z1)
	addAssign(t5, t8)
	mul(t8, x5, z0)
	addAssign(t5, t8)

	a[0][0].set(z0)
	a[0][1].set(z1)
	a[0][2].set(z2)
	a[1][0].set(z3)
	a[1][1].set(z4)
	a[1][2].set(z5)
}
