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

package bender

import (
	"math/rand"
	"time"
)

// ExponentialIntervalGenerator creates an IntervalGenerator that outputs
// exponentially distributed intervals. The resulting arrivals constitute
// a Poisson process. The rate parameter is the average queries per second
// for the generator, and corresponds to the reciprocal of the lambda parameter
// to an exponential distribution. In English, if you want to generate 30 QPS
// on average, pass 30 as the value of rate.
func ExponentialIntervalGenerator(rate float64) IntervalGenerator {
	rate = rate / float64(time.Second)
	return func(_ int64) int64 {
		return int64(rand.ExpFloat64() / rate)
	}
}

// RampedExponentialIntervalGenerator increases the average rate linearely
// from 0 to rate during the given ramp duration and keeps rate afterwards.
func RampedExponentialIntervalGenerator(rate float64, ramp time.Duration) IntervalGenerator {
	rate = rate / float64(time.Second)
	framp := float64(ramp)
	start := time.Now()
	return func(_ int64) int64 {
		elapsed := float64(time.Since(start))
		factor := elapsed / framp
		if factor > 1 {
			factor = 1
		} else if factor < 0.05 {
			// Do not let effective rate drop below 20% of the
			// final targe rate: Otherwise the tests does not
			// realy start as basicaly no request are generated.
			factor = 0.05
		}
		return int64(rand.ExpFloat64() / (factor * rate))
	}
}

// UniformIntervalGenerator creates and IntervalGenerator that outputs 1/rate every time it is
// called. Boring, right?
func UniformIntervalGenerator(rate float64) IntervalGenerator {
	irate := int64(rate / float64(time.Second))
	return func(_ int64) int64 {
		return irate
	}
}
