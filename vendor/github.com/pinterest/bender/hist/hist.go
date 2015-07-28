/*
Copyright 2014 Pinterest.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hist

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Histogram struct {
	start  int
	end    int
	scale  int
	max    int
	n      int
	errCnt int
	total  int
	values []int
}

func NewHistogram(max int, scale int) *Histogram {
	return &Histogram{0, 0, scale, max, 0, 0, 0, make([]int, max+1)}
}

func (h *Histogram) Start(t int) {
	h.start = t
}

func (h *Histogram) End(t int) {
	h.end = t
}

func (h *Histogram) Add(v int) {
	v = int(float64(v) / float64(h.scale))
	if v < 1 {
		h.values[0]++
	} else if v >= h.max {
		h.values[h.max]++
	} else {
		h.values[v]++
	}
	h.n++
	h.total += v
}

func (h *Histogram) AddError(v int) {
	h.Add(v)
	h.errCnt += 1
}

func (h *Histogram) Percentiles(percentiles ...float64) []int {
	result := make([]int, len(percentiles))
	if percentiles == nil || len(percentiles) == 0 {
		return result
	}

	sort.Sort(sort.Float64Slice(percentiles))

	accum := 0
	p_idx := int(math.Max(1.0, percentiles[0]*float64(h.n)))
	for i, j := 0, 0; i < len(percentiles) && j < len(h.values); j++ {
		accum += h.values[j]

		for accum >= p_idx {
			result[i] = j
			i++
			if i >= len(percentiles) {
				break
			}
			p_idx = int(math.Max(1.0, percentiles[i]*float64(h.n)))
		}
	}

	return result
}

func (h *Histogram) Average() float64 {
	return float64(h.total) / float64(h.n)
}

func (h *Histogram) ErrorPercent() float64 {
	return float64(h.errCnt) / float64(h.n) * 100.0
}

func (h *Histogram) String() string {
	ps := h.Percentiles(0.0, 0.25, 0.5, 0.75, 0.9, 0.95, 0.98, 0.999, 1.0)
	s := "Percentiles:\n" +
		" Min:     %d\n" +
		" 25th:    %d\n" +
		" Median:  %d\n" +
		" 75th:    %d\n" +
		" 90th:    %d\n" +
		" 95th:    %d\n" +
		" 98th:    %d\n" +
		" 99.9th:  %d\n" +
		" Max:     %d\n" +
		"Stats:\n" +
		" Average: %f\n" +
		" Total requests: %d\n" +
		" Elapsed Time (sec): %.4f\n" +
		" Average QPS: %.2f\n" +
		" Errors: %d\n" +
		" Percent errors: %.2f\n"
	elapsedSecs := float64(h.end-h.start) / float64(time.Second)
	averageQPS := float64(h.n) / elapsedSecs
	return fmt.Sprintf(s, ps[0], ps[1], ps[2], ps[3], ps[4], ps[5], ps[6], ps[7], ps[8],
		h.Average(), h.n, elapsedSecs, averageQPS, h.errCnt, h.ErrorPercent())
}
