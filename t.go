package fx

import (
	"fmt"
	"strings"
)

const (
	frBits = 52
	iBits  = 63 - frBits
	frMask = (1 << frBits) - 1
	iMask  = ((1 << iBits) - 1) << frBits
)

const (
	Sign = -1
	Zero = 0
	One  = 1 << frBits
	Iota = 1
	Max  = iMask | frMask
)

type T int64

func (t T) String() string {
	s := ""
	if t < 0 {
		t = -t
		s = "-"
	}
	ds := decimal(int64(t & frMask))
	f := string(ds[1:])
	i := fmt.Sprintf("%d", (t>>frBits)+T((ds[0]-'0')))
	return strings.Join([]string{s, i, ".", f}, "")
}

func decimal(v int64) []byte {
	u := v
	var buf [32]byte
	buf[0] = '0'
	n := 1
	o := int64(One)
	o1 := o
	for {
		o1 = o / 10
		if o1 <= u {
			break
		}
		buf[n] = '0'
		n++
		o = o1
	}
	u *= 10
	for o > 0 {
		buf[n] = '0' + byte(u/o)
		m := n
		for buf[m] == '0'+10 {
			buf[m] = '0'
			buf[m-1]++
			m--
		}
		u = u % o
		o /= 10
		n++
	}
	m := n - 1
	for m >= 1 && buf[m] == '0' {
		m--
	}
	n = m + 1
	if n == 1 {
		n = 2
	}
	return buf[:n]
}

func (a T) Frac() T {
	return a & (frMask | Sign)
}

func (a T) Int() T {
	return a & (iMask | Sign)
}

func (a T) Mul(b T) T {
	s, aa, ab := mulSign(a, b)
	low, high := mulBits(uint64(aa), uint64(ab))
	f := (high << (iBits + 1)) | (low >> frBits) + ((low >> (frBits - 1)) & 1)
	return s * T(f)
}

func (a T) Div(b T) T {
	s, aa, bb := mulSign(a, b)
	u := &u128{hi: uint64(aa) >> (64 - frBits), lo: uint64(aa) & ((1 << (iBits + 1)) - 1)}
	u.lo <<= frBits
	v := u.divBits(uint64(bb))
	return s * T(v)
}

func (a T) Inv() T {
	return T(One).Div(a)
}

// Sqrt

// Atan

// SinCos

// Pi

// e

// Sqrt2

func mulSign(a, b T) (s, absA, absB T) {
	s = T(1)
	if a^b < 0 {
		s = -1
	}
	if a < 0 {
		a = -a
	}
	absA = a
	if b < 0 {
		b = -b
	}
	absB = b
	return
}

func mulBits(x, y uint64) (low, high uint64) {
	const (
		shift = 32
		mask  = (1 << shift) - 1
	)
	low = x * y

	xl, xh := x&mask, x>>shift
	yl, yh := y&mask, y>>shift

	ll := xl * yl
	lh := xl * yh
	hl := xh * yl
	hh := xh * yh

	t := ll>>shift + lh&mask + hl&mask

	c := t >> shift
	t = lh>>shift + hl>>shift + c + hh&mask
	high = t & mask
	c = t >> shift
	high |= (hh>>shift + c) << shift
	return
}