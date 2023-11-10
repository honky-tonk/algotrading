package indicator

import (
	"time"
)

type Indicator interface {
}

type Price struct {
	T time.Time
	P float64
}

func sum_of_slice(s []float64) float64 {
	var r float64
	for _, v := range s {
		r += v
	}
	return r
}
