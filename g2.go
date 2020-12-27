package bw6

import (
	"errors"
	"math"
	"math/big"
)

// PointG2 is type for point in G2 and used for both affine and Jacobian representation.
// A point is accounted as in affine form if z is equal to one.
type PointG2 [3]fe

// Set sets the point p2 to p
func (p *PointG2) Set(p2 *PointG2) *PointG2 {
	p[0].set(&p2[0])
	p[1].set(&p2[1])
	p[2].set(&p2[2])
	return p
}

// Zero sets point p as point at infinity
func (p *PointG2) Zero() *PointG2 {
	p[0].zero()
	p[1].one()
	p[2].zero()
	return p
}

// IsAffine checks a G2 point whether it is in affine form.
func (p *PointG2) IsAffine() bool {
	return p[2].isOne()
}

type tempG2 struct {
	t [9]*fe
}

// G2 is struct for G2 group.
type G2 struct {
	tempG2
}

// NewG2 constructs a new G2 instance.
func NewG2() *G2 {
	t := newTempG2()
	return &G2{t}
}

func newTempG2() tempG2 {
	t := [9]*fe{}
	for i := 0; i < 9; i++ {
		t[i] = &fe{}
	}
	return tempG2{t}
}

// Q returns group order in big.Int.
func (g *G2) Q() *big.Int {
	return new(big.Int).Set(q)
}

// FromBytes constructs a new point given uncompressed byte input.
// Input string is expected to be equal to 192 bytes and concatenation of x and y cooridanates.
// (0, 0) is considered as infinity.
func (g *G2) FromBytes(in []byte) (*PointG2, error) {
	if len(in) != 2*FE_BYTE_SIZE {
		return nil, errors.New("input string should be 192 bytes")
	}
	x, err := fromBytes(in[:FE_BYTE_SIZE])
	if err != nil {
		return nil, err
	}
	y, err := fromBytes(in[FE_BYTE_SIZE:])
	if err != nil {
		return nil, err
	}
	// check if given input points to infinity
	if x.isZero() && y.isZero() {
		return g.Zero(), nil
	}
	z := new(fe).one()
	p := &PointG2{*x, *y, *z}
	if !g.IsOnCurve(p) {
		return nil, errors.New("point is not on curve")
	}
	return p, nil
}

// ToBytes serializes a point into bytes in uncompressed form.
// It returns (0, 0) if point is infinity.
func (g *G2) ToBytes(p *PointG2) []byte {
	out := make([]byte, 2*FE_BYTE_SIZE)
	if g.IsZero(p) {
		return out
	}
	g.Affine(p)
	copy(out[:FE_BYTE_SIZE], toBytes(&p[0]))
	copy(out[FE_BYTE_SIZE:], toBytes(&p[1]))
	return out
}

// New creates a new G2 Point which is equal to zero in other words point at infinity.
func (g *G2) New() *PointG2 {
	return g.Zero()
}

// Zero returns a new G2 Point which is equal to point at infinity.
func (g *G2) Zero() *PointG2 {
	return new(PointG2).Zero()
}

// One returns a new G2 Point which is equal to generator point.
func (g *G2) One() *PointG2 {
	p := &PointG2{}
	return p.Set(&g2One)
}

// IsZero returns true if given point is equal to zero.
func (g *G2) IsZero(p *PointG2) bool {
	return p[2].isZero()
}

// Equal checks if given two G2 point is equal in their affine form.
func (g *G2) Equal(p1, p2 *PointG2) bool {
	if g.IsZero(p1) {
		return g.IsZero(p2)
	}
	if g.IsZero(p2) {
		return g.IsZero(p1)
	}
	t := g.t
	square(t[0], &p1[2])
	square(t[1], &p2[2])
	mul(t[2], t[0], &p2[0])
	mul(t[3], t[1], &p1[0])
	mul(t[0], t[0], &p1[2])
	mul(t[1], t[1], &p2[2])
	mul(t[1], t[1], &p1[1])
	mul(t[0], t[0], &p2[1])
	return t[0].equal(t[1]) && t[2].equal(t[3])
}

// InCorrectSubgroup checks whether given point is in correct subgroup.
func (g *G2) InCorrectSubgroup(p *PointG2) bool {
	tmp := &PointG2{}
	g.MulScalar(tmp, p, q)
	return g.IsZero(tmp)
}

// IsOnCurve checks if G2 point is on curve.
func (g *G2) IsOnCurve(p *PointG2) bool {
	if g.IsZero(p) {
		return true
	}
	t := g.t
	square(t[0], &p[1])    // y^2
	square(t[1], &p[0])    // x^2
	mul(t[1], t[1], &p[0]) // x^3
	if p.IsAffine() {
		addAssign(t[1], b2)     // x^2 + b
		return t[0].equal(t[1]) // y^2 ?= x^3 + b
	}
	square(t[2], &p[2])   // z^2
	square(t[3], t[2])    // z^4
	mul(t[2], t[2], t[3]) // z^6
	doubleAssign(t[2])
	doubleAssign(t[2])      // b*z^6
	addAssign(t[1], t[2])   // x^3 + b * z^6
	return t[0].equal(t[1]) // y^2 ?= x^3 + b * z^6

}

// Affine returns the affine representation of the given point
func (g *G2) Affine(p *PointG2) *PointG2 {
	if g.IsZero(p) {
		return p
	}
	if !p.IsAffine() {
		t := g.t
		inverse(t[0], &p[2])
		square(t[1], t[0])
		mul(&p[0], &p[0], t[1])
		mul(t[0], t[0], t[1])
		mul(&p[1], &p[1], t[0])
		p[2].one()
	}
	return p
}

// Neg negates a G2 point p and assigns the result to the point at first argument.
func (g *G2) Neg(r, p *PointG2) *PointG2 {
	r[0].set(&p[0])
	r[2].set(&p[2])
	neg(&r[1], &p[1])
	return r
}

// Add adds two G2 points p1, p2 and assigns the result to point at first argument.
func (g *G2) Add(r, p1, p2 *PointG2) *PointG2 {
	// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html
	// add-2007-bl
	if g.IsZero(p1) {
		return r.Set(p2)
	}
	if g.IsZero(p2) {
		return r.Set(p1)
	}
	t := g.t
	square(t[7], &p1[2])    // z1z1
	mul(t[1], &p2[0], t[7]) // u2 = x2 * z1z1
	mul(t[2], &p1[2], t[7]) // z1z1 * z1
	mul(t[0], &p2[1], t[2]) // s2 = y2 * z1z1 * z1
	square(t[8], &p2[2])    // z2z2
	mul(t[3], &p1[0], t[8]) // u1 = x1 * z2z2
	mul(t[4], &p2[2], t[8]) // z2z2 * z2
	mul(t[2], &p1[1], t[4]) // s1 = y1 * z2z2 * z2
	if t[1].equal(t[3]) {
		if t[0].equal(t[2]) {
			return g.Double(r, p1)
		} else {
			return r.Zero()
		}
	}
	subAssign(t[1], t[3])      // h = u2 - u1
	ldouble(t[4], t[1])        // 2h
	square(t[4], t[4])         // i = 2h^2
	mul(t[5], t[1], t[4])      // j = h*i
	subAssign(t[0], t[2])      // s2 - s1
	ldoubleAssign(t[0])        // r = 2*(s2 - s1)
	square(t[6], t[0])         // r^2
	subAssign(t[6], t[5])      // r^2 - j
	mul(t[3], t[3], t[4])      // v = u1 * i
	double(t[4], t[3])         // 2*v
	sub(&r[0], t[6], t[4])     // x3 = r^2 - j - 2*v
	sub(t[4], t[3], &r[0])     // v - x3
	mul(t[6], t[2], t[5])      // s1 * j
	doubleAssign(t[6])         // 2 * s1 * j
	mul(t[0], t[0], t[4])      // r * (v - x3)
	sub(&r[1], t[0], t[6])     // y3 = r * (v - x3) - (2 * s1 * j)
	ladd(t[0], &p1[2], &p2[2]) // z1 + z2
	square(t[0], t[0])         // (z1 + z2)^2
	subAssign(t[0], t[7])      // (z1 + z2)^2 - z1z1
	subAssign(t[0], t[8])      // (z1 + z2)^2 - z1z1 - z2z2
	mul(&r[2], t[0], t[1])     // z3 = ((z1 + z2)^2 - z1z1 - z2z2) * h
	return r
}

// Double doubles a G2 point p and assigns the result to the point at first argument.
func (g *G2) Double(r, p *PointG2) *PointG2 {
	// http://www.hyperelliptic.org/EFD/gp/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	if g.IsZero(p) {
		return r.Set(p)
	}
	t := g.t
	square(t[0], &p[0])     // a = x^2
	square(t[1], &p[1])     // b = y^2
	square(t[2], t[1])      // c = b^2
	laddAssign(t[1], &p[0]) // b + x1
	square(t[1], t[1])      // (b + x1)^2
	subAssign(t[1], t[0])   // (b + x1)^2 - a
	subAssign(t[1], t[2])   // (b + x1)^2 - a - c
	doubleAssign(t[1])      // d = 2((b+x1)^2 - a - c)
	ldouble(t[3], t[0])     // 2a
	laddAssign(t[0], t[3])  // e = 3a
	square(t[4], t[0])      // f = e^2
	double(t[3], t[1])      // 2d
	sub(&r[0], t[4], t[3])  // x3 = f - 2d
	sub(t[1], t[1], &r[0])  // d-x3
	doubleAssign(t[2])      //
	doubleAssign(t[2])      //
	doubleAssign(t[2])      // 8c
	mul(t[0], t[0], t[1])   // e * (d - x3)
	sub(t[1], t[0], t[2])   // x3 = e * (d - x3) - 8c
	mul(t[0], &p[1], &p[2]) // y1 * z1
	r[1].set(t[1])          //
	double(&r[2], t[0])     // z3 = 2(y1 * z1)
	return r
}

// Sub subtracts two G2 points p1, p2 and assigns the result to point at first argument.
func (g *G2) Sub(c, a, b *PointG2) *PointG2 {
	d := &PointG2{}
	g.Neg(d, b)
	g.Add(c, a, d)
	return c
}

// MulScalar multiplies a point by given scalar value in big.Int and assigns the result to point at first argument.
func (g *G2) MulScalar(c, p *PointG2, e *big.Int) *PointG2 {
	q, n := &PointG2{}, &PointG2{}
	n.Set(p)
	l := e.BitLen()
	for i := 0; i < l; i++ {
		if e.Bit(i) == 1 {
			g.Add(q, q, n)
		}
		g.Double(n, n)
	}
	return c.Set(q)
}

// ClearCofactor maps given a G2 point to correct subgroup
func (g *G2) ClearCofactor(p *PointG2) {
	g.MulScalar(p, p, cofactorG2)
}

// MultiExp calculates multi exponentiation. Given pairs of G2 point and scalar values
// (P_0, e_0), (P_1, e_1), ... (P_n, e_n) calculates r = e_0 * P_0 + e_1 * P_1 + ... + e_n * P_n
// Length of points and scalars are expected to be equal, otherwise an error is returned.
// Result is assigned to point at first argument.
func (g *G2) MultiExp(r *PointG2, points []*PointG2, powers []*big.Int) (*PointG2, error) {
	if len(points) != len(powers) {
		return nil, errors.New("point and scalar vectors should be in same length")
	}
	var c uint32 = 3
	if len(powers) >= 32 {
		c = uint32(math.Ceil(math.Log2(float64(len(powers)))))
	}
	bucketSize, numBits := (1<<c)-1, uint32(g.Q().BitLen())
	windows := make([]*PointG2, numBits/c+1)
	bucket := make([]*PointG2, bucketSize)
	acc, sum := g.New(), g.New()
	for i := 0; i < bucketSize; i++ {
		bucket[i] = g.New()
	}
	mask := (uint64(1) << c) - 1
	j := 0
	var cur uint32
	for cur <= numBits {
		acc.Zero()
		bucket = make([]*PointG2, (1<<c)-1)
		for i := 0; i < len(bucket); i++ {
			bucket[i] = g.New()
		}
		for i := 0; i < len(powers); i++ {
			s0 := powers[i].Uint64()
			index := uint(s0 & mask)
			if index != 0 {
				g.Add(bucket[index-1], bucket[index-1], points[i])
			}
			powers[i] = new(big.Int).Rsh(powers[i], uint(c))
		}
		sum.Zero()
		for i := len(bucket) - 1; i >= 0; i-- {
			g.Add(sum, sum, bucket[i])
			g.Add(acc, acc, sum)
		}
		windows[j] = g.New()
		windows[j].Set(acc)
		j++
		cur += c
	}
	acc.Zero()
	for i := len(windows) - 1; i >= 0; i-- {
		for j := uint32(0); j < c; j++ {
			g.Double(acc, acc)
		}
		g.Add(acc, acc, windows[i])
	}
	return r.Set(acc), nil
}
