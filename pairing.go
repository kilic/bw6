package bw6

const (
	TWIST_TYPE_D = 0
	TWIST_TYPE_M = 1
)

type pair struct {
	g1 *Point
	g2 *Point
}

func newPair(g1 *Point, g2 *Point) pair {
	return pair{g1, g2}
}

type Engine struct {
	g   *G
	fp6 *fp6
	fp3 *fp3
	pairingEngineTemp
	pairs       []pair
	twistType   int
	xIsNeg      bool
	ateLoop1Neg bool
	ateLoop2Neg bool
}

// NewEngine creates new pairing engine insteace.
func NewEngine() *Engine {
	fp3 := newFp3()
	fp6 := newFp6(fp3)
	return &Engine{
		fp6:               fp6,
		fp3:               fp3,
		g:                 NewG(),
		twistType:         twistType,
		pairingEngineTemp: newEngineTemp(),
		xIsNeg:            xIsNeg,
		ateLoop1Neg:       ateLoop1Neg,
		ateLoop2Neg:       ateLoop2Neg,
	}
}

type pairingEngineTemp struct {
	t2  [10]*fe
	t12 [9]fe6
}

func newEngineTemp() pairingEngineTemp {
	t2 := [10]*fe{}
	for i := 0; i < 10; i++ {
		t2[i] = &fe{}
	}
	t12 := [9]fe6{}
	return pairingEngineTemp{t2, t12}
}

func (e *Engine) calculate() *fe6 {
	f := e.fp6.one()
	if len(e.pairs) == 0 {
		return f
	}
	e.millerLoop(f)
	e.finalExp(f)
	return f
}

// since original f \in GT = Fp6 then f (a0+a1*x+a2x^2) where (a0,a1,a2) \in Fp
func (e *Engine) doublingStep(coeff *[3]fe, r *Point) {
	// Adaptation of Formula 3 in https://eprint.iacr.org/2010/526.pdf

	// A = X1 * Y1
	a := new(fe)
	mul(a, &r[0], &r[1])

	// B = Y1^2
	bb := new(fe)
	square(bb, &r[1])

	// B = 4 * Y1^2
	b4 := new(fe)
	double(b4, bb)
	double(b4, b4)

	// C = Z1^2
	c := new(fe)
	square(c, &r[2])

	// D = 3 * C
	d := new(fe)
	double(d, c)
	add(d, d, c)

	// E = twist_b * D
	ee := new(fe)
	mul(ee, b2, d)

	// F = 3 * E
	f := new(fe)
	double(f, ee)
	add(f, f, ee)

	// G = B+F
	g := new(fe)
	add(g, bb, f)

	// H = (Y1+Z1)^2-(B+C)
	h := new(fe)
	tmp := new(fe)
	add(h, &r[1], &r[2])
	square(h, h)
	add(tmp, bb, c)
	sub(h, h, tmp)

	// I = E-B
	i := new(fe)
	sub(i, ee, bb)

	// J = X1^2
	j := new(fe)
	square(j, &r[0])

	// E2_squared = (2E)^2
	e2squared := new(fe)
	double(e2squared, ee)
	square(e2squared, e2squared)

	// X3 = 2A * (B-F)
	double(&r[0], a)
	sub(tmp, bb, f)
	mul(&r[0], &r[0], tmp)

	// Y3 = G^2 - 3*E2^2
	square(&r[1], g)
	double(tmp, e2squared)
	add(tmp, tmp, e2squared)
	sub(&r[1], &r[1], tmp)

	// Z3 = 4 * B * H
	mul(&r[2], b4, h)

	if e.twistType == TWIST_TYPE_D { // D
		// c0 = -h
		neg(&coeff[0], h)

		// c1 = 3j
		double(&coeff[1], j)
		add(&coeff[1], &coeff[1], j)

		// c2 = i
		coeff[2].set(i)
	} else { // M
		// c0 = i
		coeff[0].set(i)

		// c1 = 3j
		double(&coeff[1], j)
		add(&coeff[1], &coeff[1], j)

		// c2 = -h
		neg(&coeff[2], h)
	}
}

// since original f \in GT = Fp6 then f (a0+a1*x+a2x^2) where (a0,a1,a2) \in Fp
func (e *Engine) additionStep(coeff *[3]fe, r, q *Point) {
	tmp := new(fe)

	// theta = Y1 - Y2*Z1
	theta := new(fe)
	mul(tmp, &q[1], &r[2])
	sub(theta, &r[1], tmp)

	// lambda = X1 - X2*Z1
	mul(tmp, &q[0], &r[2])
	lambda := new(fe)
	sub(lambda, &r[0], tmp)

	// c = E^2
	c := new(fe)
	square(c, theta)

	// d = D^2
	d := new(fe)
	square(d, lambda)

	// e = D*F
	ee := new(fe)
	mul(ee, lambda, d)

	// f = Z1 * c
	f := new(fe)
	mul(f, &r[2], c)

	// g = X1 *d
	g := new(fe)
	mul(g, &r[0], d)

	// h = e + f - 2*g
	h := new(fe)
	double(tmp, g)
	add(h, ee, f)
	subAssign(h, tmp)

	// X3 = lambda * h
	mul(&r[0], lambda, h)

	// Y3 = theta * (g- h) - (e * Y1)
	mul(tmp, ee, &r[1])
	sub(&r[1], g, h)
	mul(&r[1], &r[1], theta)
	subAssign(&r[1], tmp)

	// Z3 = X1 * ee
	mul(&r[2], &r[2], ee)

	// j = theta * X2 - (lambda * Y2)
	j := new(fe)
	mul(tmp, lambda, &q[1])
	mul(j, theta, &q[0])
	subAssign(j, tmp)

	if e.twistType == TWIST_TYPE_D {
		// c0 = lambda
		coeff[0].set(lambda)

		// c1 = -theta
		neg(&coeff[1], theta)

		// c2 = j
		coeff[2].set(j)
	} else { // M

		// c0 = j
		coeff[0].set(j)

		// c1 = -theta
		neg(&coeff[1], theta)

		// c2 = lambda
		coeff[2].set(lambda)
	}
}

func (e *Engine) preCompute(ellCoeffs *[288][3]fe, twistPoint *Point) {
	if e.g.IsZero(twistPoint) {
		return
	}
	r1 := new(Point).Set(twistPoint)
	j := 0
	// f_{u+1,Q}(P)
	for i := int(ateLoop1.BitLen() - 2); i >= 0; i-- {
		e.doublingStep(&ellCoeffs[j], r1)
		j++
		if ateLoop1.Bit(i) != 0 {
			ellCoeffs[j] = fe3{}
			e.additionStep(&ellCoeffs[j], r1, twistPoint)
			j++
		}
	}

	r2 := new(Point).Set(twistPoint)
	negTwist := e.g.New()
	e.g.Neg(negTwist, twistPoint)

	// f_{u^3-u^2-u,Q}(P)
	for i := len(ateLoop2) - 2; i >= 0; i-- {
		e.doublingStep(&ellCoeffs[j], r2)
		j++
		switch ateLoop2[i] {
		case 1:
			e.additionStep(&ellCoeffs[j], r2, twistPoint)
			j++
		case -1:
			e.additionStep(&ellCoeffs[j], r2, negTwist)
			j++
		}
	}

}

func (e *Engine) ell(f *fe6, coeffs *[3]fe, p *Point) {
	c0, c1, c2 := new(fe), new(fe), new(fe)
	c0.set(&coeffs[0])
	c1.set(&coeffs[1])
	c2.set(&coeffs[2])

	switch e.twistType {
	case TWIST_TYPE_M:
		mul(c1, c1, &p[0])
		mul(c2, c2, &p[1])
		e.fp6.mulBy014Assign(f, c0, c1, c2)
	case TWIST_TYPE_D:
		//
		mul(c0, c0, &p[1])
		mul(c1, c1, &p[0])
		e.fp6.mulBy034Assign(f, c0, c1, c2)
	}
}

// optimal ate pairing: ate_opt(P, q) = (f_{u+1,Q}(P)* f_{u^3 - u^2 -u, Q}(P))
// so pairing function can be consideres as two parts
// first part is computation of f_{u+1,Q}(P)
// 	- f_{u+1, Q} = f_{u,Q}*l_{[u]Q,Q}
// second part is computation of f_{u^3 - u^2 - u, Q}(P)
// - second part can be simplified as follows:
// 		u*(u^2-u-1) -> u from first step can be user here
// 		k = u^2 - u - 1 can be written in non-adjecant form(NAF)
//		so evaluation will happen at (-Q) when bit is (-1)
// 		bitlen of k is 288 and hamming-weight of k is 19
func (e *Engine) millerLoop(f *fe6) {
	pairs := e.pairs
	ellCoeffs := make([][288][3]fe, len(pairs)) // TODO: why 100
	// ellCoeffs := make([])
	for i := 0; i < len(pairs); i++ {
		e.preCompute(&ellCoeffs[i], pairs[i].g2)
	}
	f1 := e.fp6.one()

	j := 0
	// f_{u+1,Q}(P)
	for i := int(ateLoop1.BitLen() - 2); i >= 0; i-- {
		e.fp6.square(f1, f1)
		for k := 0; k < len(pairs); k++ {
			e.ell(f1, &ellCoeffs[k][j], pairs[k].g1)
		}
		j++

		if ateLoop1.Bit(i) != 0 {
			for k := 0; k < len(pairs); k++ {
				e.ell(f1, &ellCoeffs[k][j], pairs[k].g1)
			}
			j++
		}

	}

	if e.ateLoop1Neg {
		e.fp6.conjugate(f1, f1)
	}

	f2 := e.fp6.one()
	bitLen := len(ateLoop2) - 2
	// f_{u^3-u^2-u,Q}(P)
	for i := bitLen; i >= 0; i-- {
		if i != bitLen {
			e.fp6.square(f2, f2)
		}
		for k := 0; k < len(pairs); k++ {
			e.ell(f2, &ellCoeffs[k][j], pairs[k].g1)
		}
		j++

		switch ateLoop2[i] {
		case 1:
			for k := 0; k < len(pairs); k++ {
				e.ell(f2, &ellCoeffs[k][j], pairs[k].g1)
			}
			j++
		case -1:
			for k := 0; k < len(pairs); k++ {
				e.ell(f2, &ellCoeffs[k][j], pairs[k].g1)
			}
			j++
		}
	}

	if e.ateLoop2Neg {
		e.fp6.conjugate(f2, f2)
	}
	e.fp6.frobeniusMap(f2, f2, 1)

	e.fp6.mul(f, f1, f2)
}

func (e *Engine) exp(c, a *fe6) {
	fp6 := e.fp6
	t0, t1, t2 := new(fe6).set(a), new(fe6), new(fe6)
	fp6.cyclotomicSquaring(c, t0)
	fp6.mul(t1, t0, c)
	for i := 0; i < 4; i++ {
		fp6.cyclotomicSquaring(c, c)
	}
	fp6.mul(t2, c, t0)
	fp6.cyclotomicSquaring(c, t2)
	for i := 0; i < 6; i++ {
		fp6.cyclotomicSquaring(c, c)
	}
	fp6.mul(c, c, t2)
	for i := 0; i < 5; i++ {
		fp6.cyclotomicSquaring(c, c)
	}
	fp6.mul(c, c, t1)
	for i := 0; i < 46; i++ {
		fp6.cyclotomicSquaring(c, c)
	}
	fp6.mul(c, c, t0)
}

// (q^k-1)/r where k = 6
func (e *Engine) finalExp(f *fe6) {
	// (q^6-1)/r
	fp6 := e.fp6
	inv := new(fe6)
	fp6.inverse(inv, f)

	// easy part f^(q^3-1)*(q+1)
	//  f1 = f^(q^3)*f^(-1)
	//  f2 = f^q * f1
	easyResult, tmp := new(fe6), new(fe6)
	fp6.conjugate(tmp, f)
	fp6.mul(tmp, tmp, inv)
	fp6.frobeniusMap(easyResult, tmp, 1)
	fp6.mul(easyResult, easyResult, tmp)

	// hard part (q^2-q+1)/r
	// R_0(x) * q*R_1(x)
	// where
	// R_0(x) = (-103*x^7 + 70*x^6 + 269*x^5 - 197*x^4 - 314*x^3 - 73*x^2 - 263*x - 220)
	// R_1(x) = (103*x^9 - 276*x^8 + 77*x^7 + 492*x^6 - 445*x^5 - 65*x^4 + 452*x^3 - 181*x^2 + 34*x + 229)
	f0, f0p := new(fe6), new(fe6)
	f0.set(easyResult)
	fp6.frobeniusMap(f0p, f0, 1)

	f1, f1p := new(fe6), new(fe6)
	e.exp(f1, f0)
	fp6.frobeniusMap(f1p, f1, 1)

	f2, f2p := new(fe6), new(fe6)
	e.exp(f2, f1)
	fp6.frobeniusMap(f2p, f2, 1)

	f3, f3p := new(fe6), new(fe6)
	e.exp(f3, f2)
	fp6.frobeniusMap(f3p, f3, 1)

	f4, f4p := new(fe6), new(fe6)
	e.exp(f4, f3)
	fp6.frobeniusMap(f4p, f4, 1)

	f5, f5p := new(fe6), new(fe6)
	e.exp(f5, f4)
	fp6.frobeniusMap(f5p, f5, 1)

	f6, f6p := new(fe6), new(fe6)
	e.exp(f6, f5)
	fp6.frobeniusMap(f6p, f6, 1)

	f7, f7p := new(fe6), new(fe6)
	e.exp(f7, f6)
	fp6.frobeniusMap(f7p, f7, 1)

	f8p, f9p := new(fe6), new(fe6)
	e.exp(f8p, f7p)
	e.exp(f9p, f8p)

	// step 5
	result1, f5pp3 := new(fe6), new(fe6)
	fp6.conjugate(f5pp3, f5p)
	fp6.mul(result1, f3p, f6p)
	fp6.mul(result1, result1, f5pp3)

	// step 6
	result2, f42p := new(fe6), new(fe6)
	fp6.mul(f42p, f4, f2p)
	fp6.mul(tmp, f0, f1)
	fp6.mul(tmp, tmp, f3)
	fp6.mul(tmp, tmp, f42p)
	fp6.mul(tmp, tmp, f8p)
	fp6.conjugate(tmp, tmp)
	fp6.square(result2, result1)
	fp6.mul(result2, result2, f5)
	fp6.mul(result2, result2, f0p)
	fp6.mul(result2, result2, tmp)

	result3 := new(fe6)
	fp6.conjugate(tmp, f7)
	fp6.square(result3, result2)
	fp6.mul(result3, result3, f9p)
	fp6.mul(result3, result3, tmp)

	result4, f24p, f42p5p := new(fe6), new(fe6), new(fe6)
	fp6.mul(f24p, f2, f4p)
	fp6.mul(f42p5p, f42p, f5p)
	fp6.mul(tmp, f24p, f3)
	fp6.mul(tmp, tmp, f3p)
	fp6.conjugate(tmp, tmp)
	fp6.square(result4, result3)
	fp6.mul(result4, result4, f42p5p)
	fp6.mul(result4, result4, f6)
	fp6.mul(result4, result4, f7p)
	fp6.mul(result4, result4, tmp)

	result5 := new(fe6)
	fp6.mul(tmp, f0p, f9p)
	fp6.conjugate(tmp, tmp)
	fp6.square(result5, result4)
	fp6.mul(result5, result5, f0)
	fp6.mul(result5, result5, f7)
	fp6.mul(result5, result5, f1p)
	fp6.mul(result5, result5, tmp)

	result6, f6p8p, f57p := new(fe6), new(fe6), new(fe6)
	fp6.mul(f6p8p, f6p, f8p)
	fp6.mul(f57p, f5, f7p)
	fp6.conjugate(tmp, f6p8p)
	fp6.square(result6, result5)
	fp6.mul(result6, result6, f57p)
	fp6.mul(result6, result6, f2p)
	fp6.mul(result6, result6, tmp)

	result7, f17, f36 := new(fe6), new(fe6), new(fe6)
	fp6.mul(f36, f3, f6)
	fp6.mul(f17, f1, f7)
	fp6.mul(tmp, f17, f2)
	fp6.conjugate(tmp, tmp)
	fp6.square(result7, result6)
	fp6.mul(result7, result7, f36)
	fp6.mul(result7, result7, f9p)
	fp6.mul(result7, result7, tmp)

	result8 := new(fe6)
	fp6.mul(tmp, f42p, f57p)
	fp6.mul(tmp, tmp, f6p8p)
	fp6.conjugate(tmp, tmp)
	fp6.square(result8, result7)
	fp6.mul(result8, result8, f0)
	fp6.mul(result8, result8, f0p)
	fp6.mul(result8, result8, f3p)
	fp6.mul(result8, result8, f5p)
	fp6.mul(result8, result8, tmp)

	result9 := new(fe6)
	fp6.conjugate(tmp, f36)
	fp6.square(result9, result8)
	fp6.mul(result9, result9, f1p)
	fp6.mul(result9, result9, tmp)

	result10 := new(fe6)
	fp6.mul(tmp, f24p, f42p5p)
	fp6.mul(tmp, tmp, f9p)
	fp6.conjugate(tmp, tmp)
	fp6.square(result10, result9)
	fp6.mul(result10, result10, f17)
	fp6.mul(result10, result10, f57p)
	fp6.mul(result10, result10, f0p)
	fp6.mul(result10, result10, tmp)

	f.set(result10)
}

// AddPair adds a g1, g2 point pair to pairing engine
func (e *Engine) AddPair(g1 *Point, g2 *Point) *Engine {
	p := newPair(g1, g2)
	if !e.isZero(p) {
		e.affine(p)
		e.pairs = append(e.pairs, p)
	}
	return e
}

// AddPairInv adds a G1, G2 point pair to pairing engine. G1 point is negated.
func (e *Engine) AddPairInv(g1 *Point, g2 *Point) *Engine {
	ng1 := e.g.New().Set(g1)
	e.g.Neg(ng1, g1)
	e.AddPair(ng1, g2)
	return e
}

// Reset deletes added pairs.
func (e *Engine) Reset() *Engine {
	e.pairs = []pair{}
	return e
}

func (e *Engine) isZero(p pair) bool {
	return e.g.IsZero(p.g1) || e.g.IsZero(p.g2)
}

func (e *Engine) affine(p pair) {
	e.g.Affine(p.g1)
	e.g.Affine(p.g2)
}

// Result computes pairing and returns target group element as result.
func (e *Engine) Result() *E {
	r := e.calculate()
	e.Reset()
	return r
}

// GT returns target group instance.
func (e *Engine) GT() *GT {
	return NewGT()
}

// Check computes pairing and checks if result is equal to one
func (e *Engine) Check() bool {
	return e.calculate().isOne()
}
