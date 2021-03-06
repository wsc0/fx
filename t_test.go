// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package fx

import (
	"math"
	"math/rand"
	"testing"
)

func TestString(t *testing.T) {
	t.Logf("%s", T(One))
	t.Logf("%s", T(Iota))
	t.Logf("%s", T(One>>1))
	t.Logf("%s", T(One>>2))
	t.Logf("%s", T(One+One+One).Mul(One>>1))
}

func TestMulId(t *testing.T) {
	N := 1024
	for i := 0; i < N; i++ {
		n := T(rand.Int63n((1 << 63) - 1))

		m := n.Mul(One)
		if n != m {
			t.Errorf("%d: One isn't identify for %s\no %b\nm %b", i, n, n, m)
		}
	}
}

func TestDivId(t *testing.T) {
	N := 1024
	for i := 0; i < N; i++ {
		n := T(rand.Int63n((1 << 63) - 1))
		m := n.Div(One)
		if n != m {
			t.Errorf("%d: One isn't identify for %s\no %b\nm %b", i, n, n, m)
		}
	}
}

func TestInvPo2(t *testing.T) {
	for i := uint(0); i < iBits; i++ {
		n := T(One << i)
		m := n.Inv().Inv()
		if m != n {
			t.Errorf("%d: inv(inv(%s)) gave %s\n", i, n, m)
		}
	}
}

func TestInvMulClose(t *testing.T) {
	for i := 1; i < (1 << iBits); i++ {
		n := Int(i)
		f := n.Mul(n.Inv())
		e := One - f
		if e > Iota<<10 {
			t.Errorf("%s * %s = %s\n", n, n.Inv(), e)
		}
	}
}

func TestFloat64Conv(t *testing.T) {
	N := 1024
	eps := 1e-16
	for i := 0; i < N; i++ {
		f := rand.Float64() * (1 << iBits)
		n := Float64(f)
		nf := n.Float64()
		if math.Abs(f-nf) > eps {
			t.Errorf("%f -> %s -> %f: e %0.16f", f, n, nf, math.Abs(f-nf))
		}
	}
}

func TestSqrt(t *testing.T) {
	a := Int(9)
	t.Logf("sqrt 9 %s ^2 %s\n", Sqrt(a), Sqrt(a).Mul(Sqrt(a)))
	t.Logf("sqrt2const from float %s calculated %s\n", T(Sqrt2), Sqrt(One+One))
	t.Logf("sqrt(0.5) vs Sqrt2.Inv(): %s %s\n", Sqrt(One/2), T(Sqrt2).Inv())
}
