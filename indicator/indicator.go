package indicator

import (
	"time"
)

type Indicator interface {
}

type Price struct {
	t time.Time
	p float64
}
