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
	if len(in) != 6*fpByteSize {
		return nil, errors.New("input string should be larger than 96 bytes")
	}
	fp3 := e.fp3
	c0, err := fp3.fromBytes(in[:3*fpByteSize])
	if err != nil {
		return nil, err
	}
	c1, err := fp3.fromBytes(in[3*fpByteSize:])
	if err != nil {
		return nil, err
	}
	return &fe6{*c0, *c1}, nil
}

func (e *fp6) toBytes(a *fe6) []byte {
	out := make([]byte, 6*fpByteSize)
	fp3 := e.fp3
	copy(out[:3*fpByteSize], fp3.toBytes(&a[0]))
	copy(out[3*fpByteSize:], fp3.toBytes(&a[1]))
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
	c[0].set(&a[0])
	fp3.neg(&c[1], &a[1])
}

func (e *fp6) mul(c, a, b *fe6) {
	// Multiplication and Squaring on Pairing-Friendly Fields
	// Karatsuba multiplication algorithm
	// https://eprint.iacr.org/2006/471

	fp3, t := e.fp3, e.t

	fp3.mul(t[1], &a[0], &b[0]) // v0 = a0b0
	fp3.mul(t[2], &a[1], &b[1]) // v1 = a1b1

	fp3.add(t[0], &a[0], &a[1]) // a0 + a1
	fp3.add(t[3], &b[0], &b[1]) // b0 + b1
	fp3.mul(t[0], t[0], t[3])   // (a0 + a1)(b0 + b1)
	fp3.sub(t[0], t[0], t[1])   // (a0 + a1)(b0 + b1) - v0
	fp3.sub(&c[1], t[0], t[2])  // c1 = (a0 + a1)(b0 + b1) - v0 - v1

	fp3.mulByNonResidue(t[2], t[2])
	fp3.add(&c[0], t[1], t[2]) // c0 = v0 - ßv1
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
	// c0 = v0 + αv1 = v0 - ßv1
	// c1 = (a0 + a1)^2 - v0 - v1

	fp3, t := e.fp3, e.t
	fp3.square(t[0], &a[0]) // v0 = a0^2
	fp3.square(t[1], &a[1]) // v1 = a1^2

	fp3.mulByNonResidue(t[2], t[2])
	fp3.sub(t[3], t[0], t[2]) // c0 = v0 - ßv1

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
	// c0 = (a0 + a1)(a0 + ßa1) - v0 - ßv0
	// c1 = 2v0
	fp3, t := e.fp3, e.t
	fp3.mulByNonResidue(t[0], &a[1]) // ßa1
	fp3.mul(t[1], &a[0], &a[1])      // v0 = a0a1
	fp3.mulByNonResidue(t[2], t[1])  // ßv0

	fp3.add(t[0], t[0], &a[0]) // a0 + ßa1
	fp3.add(t[2], t[2], t[1])  // v0 + ßv0

	fp3.add(t[3], &a[0], &a[1]) // a0 + a1
	fp3.mul(t[0], t[0], t[3])   // (a0 + a1)(a0 + ßa1)

	fp3.sub(&c[0], t[0], t[2]) // (a0 + a1)(a0 + ßa1) - v0 - ßv0
	fp3.double(&c[1], t[1])    // 2v0
}

func (e *fp6) inverse(c, a *fe6) {
	// Guide to Pairing Based Cryptography
	// Algorithm 5.19

	fp3, t := e.fp3, e.t

	fp3.square(t[0], &a[0]) // a0^2
	fp3.square(t[1], &a[1]) // a1^2

	fp3.mulByNonResidue(t[2], t[1])
	fp3.sub(t[0], t[0], t[2]) // v = a0^2 + ßa1^2
	fp3.inverse(t[1], t[0])   // v = v^-1

	fp3.mul(&c[0], t[1], &a[0]) // a0v
	fp3.mul(t[1], t[1], &a[1])  // a1v
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
func (e *fp6) optimized_exp(c, a *fe6, s *big.Int) {
	naf := computeNaf(s)

	z := e.one()

	inv := new(fe6)
	e.inverse(inv, a)
	foundNonZero := false
	for i := len(naf) - 1; i >= 0; i-- {
		if foundNonZero {
			e.square(z, z)
		}
		switch naf[i] {
		case 1:
			foundNonZero = true
			e.mul(z, z, a)
		case -1:
			foundNonZero = true
			e.mul(z, z, inv)
		}

	}

	c.set(z)
}

func (e *fp6) frobeniusMap(c, a *fe6, power int) {
	fp3 := e.fp3
	fp3.frobeniusMap(&c[0], &a[0], power)
	fp3.frobeniusMap(&c[1], &a[1], power)
	fp3.mul0(&c[1], &c[1], &frobeniusCoeffs6[power%6])
}

// TODO: zexe uses different name which is mul by 034
// lets say {a,b} $\in$ Fp6 where
// in sparse multiplication of a*b
// a = (a00 + a01*u + a02*u^2) + (a10 + a11*u + a12* u^2)
// b = (b00 + 0 + b02*u^2) + (0 + b11*u + 0)
func (e *fp6) mulBy034Assign(a *fe6, c0, c3, c4 *fe) {
	// z0, z1, z2, z3, z4, z5 := &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}
	// z0.set(&a[0][0])
	// z1.set(&a[0][1])
	// z2.set(&a[0][2])
	// z3.set(&a[1][0])
	// z4.set(&a[1][1])
	// z5.set(&a[1][2])

	x0, x3, x4 := &fe{}, &fe{}, &fe{}
	x0.set(c0)
	x3.set(c3)
	x4.set(c4)

	tmp1, tmp2 := new(fe), new(fe)
	mul(tmp1, x3, nonResidue1)
	mul(tmp2, x4, nonResidue1)

	t0, t1, t2 := &fe{}, &fe{}, &fe{}

	mul(t0, &a[0][0], x0)
	mul(t1, tmp1, &a[1][2])
	mul(t2, tmp2, &a[1][1])
	add(&a[0][0], t0, t1)
	addAssign(&a[0][0], t2)

	mul(t0, x0, &a[0][1])
	mul(t1, x3, &a[1][0])
	mul(t2, tmp2, &a[1][2])
	add(&a[0][1], t0, t1)
	addAssign(&a[0][1], t2)

	mul(t0, x0, &a[0][2])
	mul(t1, x3, &a[1][1])
	mul(t2, x4, &a[1][0])
	add(&a[0][2], t0, t1)
	addAssign(&a[0][2], t2)

	mul(t0, x0, &a[1][0])
	mul(t1, x3, &a[0][0])
	mul(t2, tmp2, &a[0][2])
	add(&a[1][0], t0, t1)
	addAssign(&a[1][0], t2)

	mul(t0, x0, &a[1][1])
	mul(t1, x3, &a[0][1])
	mul(t2, x4, &a[0][0])
	add(&a[1][1], t0, t1)
	addAssign(&a[1][1], t2)

	mul(t0, x0, &a[1][2])
	mul(t1, x3, &a[0][2])
	mul(t2, x4, &a[0][1])
	add(&a[1][2], t0, t1)
	addAssign(&a[1][2], t2)
}

// TODO: zexe uses different name which is mul by 034
// lets say {a,b} \in Fp6 where
// in sparse multiplication of a*b
// a = (a00 + a01*u + a02*u^2) + (a10 + a11*u + a12*u^2)
// b = (b00 + 0 + 0) + (0 + b11*u + b12*u^2)
func (e *fp6) mulBy014Assign(a *fe6, c0, c1, c4 *fe) {
	z0, z1, z2, z3, z4, z5 := &fe{}, &fe{}, &fe{}, &fe{}, &fe{}, &fe{}
	z0.set(&a[0][0])
	z1.set(&a[0][1])
	z2.set(&a[0][2])
	z3.set(&a[1][0])
	z4.set(&a[1][1])
	z5.set(&a[1][2])

	x0, x1, x4 := &fe{}, &fe{}, &fe{}
	x0.set(c0)
	x1.set(c1)
	x4.set(c4)

	t0, t1, t2 := &fe{}, &fe{}, &fe{}

	tmp1, tmp2 := new(fe), new(fe)

	mul(tmp1, x1, nonResidue1)

	mul(tmp2, x4, nonResidue1)

	mul(t0, x0, z0)

	mul(t1, tmp1, z2)
	mul(t2, tmp2, z4)
	add(&a[0][0], t0, t1)
	addAssign(&a[0][0], t2)

	mul(t0, x0, z1)

	mul(t1, x1, z0)
	mul(t2, tmp2, z5)
	add(&a[0][1], t0, t1)
	addAssign(&a[0][1], t2)

	mul(t0, x0, z2)
	mul(t1, x1, z1)
	mul(t2, x4, z3)
	add(&a[0][2], t0, t1)
	addAssign(&a[0][2], t2)

	mul(t0, x0, z3)
	mul(t1, tmp1, z5)
	mul(t2, tmp2, z2)
	add(&a[1][0], t0, t1)
	addAssign(&a[1][0], t2)

	mul(t0, x0, z4)
	mul(t1, x1, z3)
	mul(t2, x4, z0)
	add(&a[1][1], t0, t1)
	addAssign(&a[1][1], t2)

	mul(t0, x0, z5)
	mul(t1, x1, z4)
	mul(t2, x4, z1)
	add(&a[1][2], t0, t1)
	addAssign(&a[1][2], t2)
}
