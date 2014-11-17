// Copyright 2014 Volker Dobler.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hist provides integer histograms with exponentially growing bucket sizes.
package hist

// LogHist is a histogram of non-negative integer values whose bin size
// increases exponentially and cover the intervall [0,Max].
// The bins are grouped into blocks of N equal-sized bins. The first block has
// a bin size of 1; in each consecutive block the bin size is doubled.
// The resolution of the histogram is 1/N.
type LogHist struct {
	N         int   // Number of equal sized bins before binsize doubles. A power of two.
	Max       int   // Max is the last value which can be counted.
	Count     []int // Count contains the counts for each bucket.
	Overflow  int   // Number of added values > Max
	Underflow int   // Number of added values < 0

}

// NewLogHist returns a new LogHist capable of storing values from 0 to (at least) max
// with a resolution of bits. N will be 1<<bits.
func NewLogHist(bits uint, max int) *LogHist {
	n := 1 << bits
	h := &LogHist{
		N: n,
	}
	lastBucket := h.Bucket(max)
	_, lastValue := h.Cover(lastBucket)
	h.Count = make([]int, lastValue)
	h.Max = lastValue
	return h
}

// Bucket returns the bucket index the value v belongs to. The value v must
// be in the range [0,h.Max].
func (h *LogHist) Bucket(v int) int {
	// block p covers the value range of [n * 2^p - n , n * 2^(p+1) - n)
	n := h.N
	if v < n {
		return v
	}

	p := uint(0)
	low := n*(1<<p) - n
	for low <= v {
		p++
		low = n*(1<<p) - n
	}
	p--
	low = n*(1<<p) - n
	bucketsize := 1 << p
	return n*int(p) + (v-low)/bucketsize
}

// Cover returns the value intervall [a,b) which is covered by bucket.
func (h *LogHist) Cover(bucket int) (a int, b int) {
	// The bucket z covers the value range [a,b).
	// The bucket z is bin u = z%n in block p = z/n and is w = 1<<p values wide.
	// The bucket starts at a = n*2^p - n + u*w
	n := h.N
	u, p := bucket%n, uint(bucket/n)
	w := 1 << p
	a = n*(1<<p) - n + u*w
	return a, a + w
}

// Add counts the value v.
func (h *LogHist) Add(v int) {
	if v < 0 {
		h.Underflow++
	}
	if v >= h.Max {
		h.Overflow++
	}
	b := h.Bucket(v)
	h.Count[b]++
}

// MustAdd works like Add but will panic if v under- or overflows the
// allowed range of [0,h.Max].
func (h *LogHist) MustAdd(v int) {
	h.Count[h.Bucket(v)]++
}
