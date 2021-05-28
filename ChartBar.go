package psarf

import (
	"time"
)

// ChartBar represents the individual chart data interface required to calculate
// the Psar
type ChartBar interface {
	Date() *time.Time
	High() float64
	Low() float64

	// NOTE open and closed values are never used in the Psar formula. Keeping
	// these part of the interface for a more rounded out feel
	Open() float64
	Close() float64
}
