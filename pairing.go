package bw6

type pair struct {
	g1 *PointG1
	g2 *PointG2
}

func newPair(g1 *PointG1, g2 *PointG2) pair {
	return pair{g1, g2}
}

type Engine struct {
	G1  *G1
	G2  *G2
	fp6 *fp6
	fp3 *fp3
	pairingEngineTemp
	pairs     []pair
	twistType int // 0 D 1 M
}

// NewEngine creates new pairing engine insteace.
func NewEngine(twistType int) *Engine {
	fp3 := newFp3()
	fp6 := newFp6(fp3)
	g1 := NewG1()
	g2 := NewG2()
	return &Engine{
		fp6:               fp6,
		fp3:               fp3,
		G1:                g1,
		G2:                g2,
		twistType:         twistType,
		pairingEngineTemp: newEngineTemp(),
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
func (e *Engine) doublingStep(coeff *[3]fe, r *PointG2) {
	// Adaptation of Formula 3 in https://eprint.iacr.org/2010/526.pdf
	// fp6 := e.fp6
	// t := e.t2

	// x1, y1, z1 are fields of r

	// A = X1 * Y1

	// B = Y1^2

	// B = 4 * Y1^2

	// C = Z1^2

	// D = 3 * C

	// E = twist_b * D  // TODO: prepare twist of "b coeff"

	// F = 3 * E

	// G = B+F

	// H = (Y1+Z1)^2-(B+C)

	// I = E-B

	// J = X1^2

	// E2_squared = (2E)^2

	// X3 = 2A * (B-F) // x3, y3, z3 are elements of result

	// Y3 = G^2 - 3*E2^2

	// Z3 = 4 * B * H

	if e.twistType == 0 { // D
		// ell_0 = xi * I // update coeff 0

		// ell_VW = - H (later: * yP) // coeff 1

		// ell_VV = 3*J (later: * xP)
	} else { // M

		// ell_0 = I

		// ell_VW = -xi * H (later: * yP)

		// ell_VV = 3*J (later: * xP)
	}

}

// since original f \in GT = Fp6 then f (a0+a1*x+a2x^2) where (a0,a1,a2) \in Fp
func (e *Engine) additionStep(coeff *[3]fe, r, q *PointG2) {
	// D = X1 - X2*Z1

	// E = Y1 - Y2*Z1

	// F = D^2

	// G = E^2

	// H = D*F

	// I = X1 * F

	// J = H + Z1*G - (I+I)

	// prepare result: x3, y3, z3

	// X3 = D*J

	// Y3 = E*(I-J)-(H*Y1)

	// Z3 = Z1*H

	// update ell coefss

	if e.twistType == 0 { // D
		// c00 = xi * (E * X2 - D * Y2)

		// c1 = - E (later: * xP)

		// c2 = D (later: * yP)
	} else { // M

		// c0 = E * X2 - D * Y2

		// c1 = - E (later: * xP)

		// c2 = - E (later: * xP)
	}
}

// compute
func (e *Engine) preCompute(ellCoeffs *[129][3]fe, twistPoint *PointG2) {
	if e.G2.IsZero(twistPoint) {
		return
	}
	r := new(PointG2).Set(twistPoint)
	j := 0
	for i := int(x.BitLen() - 2); i >= 0; i-- {
		e.doublingStep(&ellCoeffs[j], r)
		if x.Bit(i) != 0 {
			j++
			ellCoeffs[j] = fe3{}
			e.additionStep(&ellCoeffs[j], r, twistPoint)
		}
		j++
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
// 		bitlen of k is 129 and hamming-weight of k is 19
func (e *Engine) millerLoop(f *fe6) {
	pairs := e.pairs
	ellCoeffs := make([][129][3]fe, len(pairs)) // TODO: why 100
	for i := 0; i < len(pairs); i++ {
		e.preCompute(&ellCoeffs[i], pairs[i].g2)
	}
	// mulby_024 and mul_by_045 sparse multiplications for Fp6
}

// (q^k-1)/r where k = 6
func (e *Engine) finalExp(f *fe6) {
	// (q^6-1)/r

	// easy part f^(q^3-1)*(q+1)
	// - f^(q^3)*f^(-1)
	// - f^q * f

	// hard part (q^2-q+1)/r
	// R_0(x) * q*R_1(x)
	// where
	// R_0(x) = (-103*x^7 + 70*x^6 + 269*x^5 - 197*x^4 - 314*x^3 - 73*x^2 - 263*x - 220)
	// R_1(x) = (103*x^9 - 276*x^8 + 77*x^7 + 492*x^6 - 445*x^5 - 65*x^4 + 452*x^3 - 181*x^2 + 34*x + 229)
}
