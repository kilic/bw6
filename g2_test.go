package bw6

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func (g *G2) rand() *PointG2 {
	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		panic(err)
	}
	return g.MulScalar(&PointG2{}, g.New().Set(&g2One), k)
}

func TestG2Serialization(t *testing.T) {
	g := NewG2()
	for i := 0; i < fuz; i++ {
		a := g.rand()
		uncompressed := g.ToBytes(a)
		b, err := g.FromBytes(uncompressed)
		if err != nil {
			t.Fatal(err)
		}
		if !g.Equal(a, b) {
			t.Fatal("bad serialization 3")
		}
	}
}

func TestG2IsOnCurve(t *testing.T) {
	g := NewG2()
	zero := g.Zero()
	if !g.IsOnCurve(zero) {
		t.Fatal("zero must be on curve")
	}
	one := new(fe).one()
	p := &PointG2{*one, *one, *one}
	if g.IsOnCurve(p) {
		t.Fatal("(1, 1) is not on curve")
	}
}

func TestG2AdditiveProperties(t *testing.T) {
	g := NewG2()
	t0, t1 := g.New(), g.New()
	zero := g.Zero()
	for i := 0; i < fuz; i++ {
		a, b := g.rand(), g.rand()
		g.Add(t0, a, zero)
		if !g.Equal(t0, a) {
			t.Fatal("a + 0 == a")
		}
		g.Add(t0, zero, zero)
		if !g.Equal(t0, zero) {
			t.Fatal("0 + 0 == 0")
		}
		g.Sub(t0, a, zero)
		if !g.Equal(t0, a) {
			t.Fatal("a - 0 == a")
		}
		g.Sub(t0, zero, zero)
		if !g.Equal(t0, zero) {
			t.Fatal("0 - 0 == 0")
		}
		g.Neg(t0, zero)
		if !g.Equal(t0, zero) {
			t.Fatal("- 0 == 0")
		}
		g.Sub(t0, zero, a)
		g.Neg(t0, t0)
		if !g.Equal(t0, a) {
			t.Fatal(" - (0 - a) == a")
		}
		g.Double(t0, zero)
		if !g.Equal(t0, zero) {
			t.Fatal("2 * 0 == 0")
		}
		g.Double(t0, a)
		g.Sub(t0, t0, a)
		if !g.Equal(t0, a) {
			t.Fatal("(2 * a) - a == a")
		}
		g.Add(t0, a, b)
		g.Add(t1, b, a)
		if !g.Equal(t0, t1) {
			t.Fatal("a + b == b + a")
		}
		g.Sub(t0, a, b)
		g.Sub(t1, b, a)
		g.Neg(t1, t1)
		if !g.Equal(t0, t1) {
			t.Fatal("a - b == - ( b - a )")
		}
		c := g.rand()
		g.Add(t0, a, b)
		g.Add(t0, t0, c)
		g.Add(t1, a, c)
		g.Add(t1, t1, b)
		if !g.Equal(t0, t1) {
			t.Fatal("(a + b) + c == (a + c ) + b")
		}
		g.Sub(t0, a, b)
		g.Sub(t0, t0, c)
		g.Sub(t1, a, c)
		g.Sub(t1, t1, b)
		if !g.Equal(t0, t1) {
			t.Fatal("(a - b) - c == (a - c) -b")
		}
	}
}

func TestG2MultiplicativeProperties(t *testing.T) {
	g := NewG2()
	t0, t1 := g.New(), g.New()
	zero := g.Zero()
	for i := 0; i < fuz; i++ {
		a := g.rand()
		s1, s2, s3 := randScalar(q), randScalar(q), randScalar(q)
		sone := big.NewInt(1)
		g.MulScalar(t0, zero, s1)
		if !g.Equal(t0, zero) {
			t.Fatal("0 ^ s == 0")
		}
		g.MulScalar(t0, a, sone)
		if !g.Equal(t0, a) {
			t.Fatal("a ^ 1 == a")
		}
		g.MulScalar(t0, zero, s1)
		if !g.Equal(t0, zero) {
			t.Fatal("0 ^ s == a")
		}
		g.MulScalar(t0, a, s1)
		g.MulScalar(t0, t0, s2)
		s3.Mul(s1, s2)
		g.MulScalar(t1, a, s3)
		if !g.Equal(t0, t1) {
			t.Fatal("(a ^ s1) ^ s2 == a ^ (s1 * s2)")
		}
		g.MulScalar(t0, a, s1)
		g.MulScalar(t1, a, s2)
		g.Add(t0, t0, t1)
		s3.Add(s1, s2)
		g.MulScalar(t1, a, s3)
		if !g.Equal(t0, t1) {
			t.Fatal("(a ^ s1) + (a ^ s2) == a ^ (s1 + s2)")
		}
	}
}

func TestG2MultiExpExpected(t *testing.T) {
	g := NewG2()
	one := g.New().Set(&g2One)
	var scalars [2]*big.Int
	var bases [2]*PointG2
	scalars[0] = big.NewInt(2)
	scalars[1] = big.NewInt(3)
	bases[0], bases[1] = new(PointG2).Set(one), new(PointG2).Set(one)
	expected, result := g.New(), g.New()
	g.MulScalar(expected, one, big.NewInt(5))
	_, _ = g.MultiExp(result, bases[:], scalars[:])
	if !g.Equal(expected, result) {
		t.Fatal("bad multi-exponentiation")
	}
}

func TestG2MultiExpBatch(t *testing.T) {
	g := NewG2()
	one := g.New().Set(&g2One)
	n := 1000
	bases := make([]*PointG2, n)
	scalars := make([]*big.Int, n)
	// scalars: [s0,s1 ... s(n-1)]
	// bases: [P0,P1,..P(n-1)] = [s(n-1)*G, s(n-2)*G ... s0*G]
	for i, j := 0, n-1; i < n; i, j = i+1, j-1 {
		scalars[j], _ = rand.Int(rand.Reader, big.NewInt(100000))
		bases[i] = g.New()
		g.MulScalar(bases[i], one, scalars[j])
	}
	// expected: s(n-1)*P0 + s(n-2)*P1 + s0*P(n-1)
	expected, tmp := g.New(), g.New()
	for i := 0; i < n; i++ {
		g.MulScalar(tmp, bases[i], scalars[i])
		g.Add(expected, expected, tmp)
	}
	result := g.New()
	_, _ = g.MultiExp(result, bases, scalars)
	if !g.Equal(expected, result) {
		t.Fatal("bad multi-exponentiation")
	}
}

func BenchmarkG2Add(t *testing.B) {
	g := NewG2()
	a, b, c := g.rand(), g.rand(), PointG2{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g.Add(&c, a, b)
	}
}

func BenchmarkG2Mul(t *testing.B) {
	g := NewG2()
	a, e, c := g.rand(), q, PointG2{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g.MulScalar(&c, a, e)
	}
}
