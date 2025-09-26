package psarf

import (
	"time"
)

const (
	// These are the Acceleration values used in the standard Psar formula
	AFIncrement = 0.02
	AFMax       = 0.20
)

// before is a utility func that checks if time a is before time b
// Note, hours are reset (this is daily only)
// FIXME allow this to be used with other timeframes
func before(a, b time.Time) bool {
	a = a.Truncate(24 * time.Hour)
	b = b.Truncate(24 * time.Hour)

	return a.Before(b)
}

type Direction int

const (
	IsLong Direction = iota
	IsShort
)

// Psar is an iterative style structure that calculates Psar values for a given
// set of chart data
type Psar struct {
	// Series is the dataset used to calculate the Psar values. The first bar in
	// the series is the base Psar value (usually this is the nearest pivot low
	// prior to the Psar reversal).
	Series []ChartBar

	// StartDate is the date desired to start the Parabolic Sar calculations.
	// THIS IS NOT TO BE CONFUSED WITH THE FIRST BAR IN THE SERIES (DATA). (For
	// my use case this would also be the Entry bar of my position) In the
	// traditional indicator, one Sar being touched begins the Sar in the other
	// direction. But as this is not being directly used a trend indicator a
	// Start date is used to as the trigger. This is also the bar it derives the
	// initial Extreme Price from.
	StartDate *time.Time

	// Direction determines if Psar is calculated for the "long" or "short" side.
	// NOTE if no Direction is set, it defaults to IsLong
	Direction Direction

	// psarSeries holds the calculated Psar values for each bar in the series.
	psarSeries []*PsarPeriod

	// pipOffset is a value to offset the sar for the initial Sar and any Sar
	// that is truncated by the lowest last 2 bars rule
	// NOTE this is not part of the psar formula, but is for my own personal use
	pipOffset float64

	// StopPrice is an optional user defined stop price. This has two use cases:
	// 1) it serves as the "last "pivot" prior to the reversal" that The
	//    Parabolic Sar uses to calculate the initial Sar value. In a normal use
	//    case the Series will start at this pivot point.
	// 2) it serves as a custom stop price, that can be used to "anchor" the
	//    formula to use that value as the initial Sar value. This is the main
	//		use case for the Psar library.
	// NOTE this value will only be used within the isEntryBar call
	StopPrice float64

	// i is an internal iteration counter (no slicing here, we need to be able
	// to go back in history)
	i int
}

// SetPipOffset sets the pipOffset value
func (p *Psar) SetPipOffset(v float64) {
	p.pipOffset = v
}

// Bar returns current PsarPeriod (within the iteration)
func (p *Psar) Bar() *PsarPeriod {
	return p.psarSeries[p.i-1]
}

// isEntryBar returns true if the current bar in the iteration's Date is equal
// to that of the given StartDate
// FIXME allow this to be used with other timeframes
func (p *Psar) isEntryBar(i int) bool {
	var (
		a = p.Series[i].Date()
		b = p.StartDate
	)
	return a.Truncate(24 * time.Hour).Equal(b.Truncate(24 * time.Hour))
}

// calculatePsar calculates the Psar for the given period in the series. This
// must be called each time the Psar is iterated through
func (p *Psar) calculatePsar() {
	var (
		i        = p.i
		chartBar = p.Series[i]
		af       = AFIncrement
		sar      float64
		ep       float64
		extVal   float64

		dir     = p.Direction
		isShort = dir == IsShort
		isLong  = dir == IsLong || !isShort // default to long if not set
	)

	// skip any calculations until we reach the start date
	if before(*chartBar.Date(), *p.StartDate) {
		// FIXME calculations should start at the start date, so this would be
		// technically unnecessary
		p.psarSeries = append(p.psarSeries, &PsarPeriod{ChartBar: chartBar})

		return
	}

	// always default to the current bar's high/low for the ep/extVal
	if isShort {
		ep, extVal = chartBar.Low(), chartBar.High()
	} else {
		ep, extVal = chartBar.High(), chartBar.Low()
	}

	// FIXME checks for start date, etc...

	if p.isEntryBar(i) {
		af = AFIncrement

		if p.StopPrice != 0 {
			// because StopPrice is user defined, there is no pipOffset applied at
			// the entry bar point
			sar = p.StopPrice
		} else if isShort {
			sar = p.Series[0].High() + p.pipOffset
		} else {
			sar = p.Series[0].Low() - p.pipOffset
		}
	} else {
		// technically, at this point, you should always have a previous bar
		var prevbp = p.psarSeries[i-1]
		af = prevbp.AF

		if isShort {
			if ep < prevbp.EP {
				af += AFIncrement
			} else {
				ep = prevbp.EP
			}
			sar = prevbp.Sar - prevbp.AFSarEP

			var h1 = prevbp.High()
			extVal = h1
			if i-2 >= 0 {
				if h2 := p.psarSeries[i-2].High(); h2 > h1 {
					extVal = h2
				}
			}
		} else {
			if ep > prevbp.EP {
				af += AFIncrement
			} else {
				ep = prevbp.EP
			}
			sar = prevbp.Sar + prevbp.AFSarEP

			var l1 = prevbp.Low()
			extVal = l1
			if i-2 >= 0 {
				if l2 := p.psarSeries[i-2].Low(); l2 < l1 {
					extVal = l2
				}
			}
		}
	}

	if isLong && sar > extVal {
		sar = extVal - p.pipOffset
	} else if isShort && sar < extVal {
		sar = extVal + p.pipOffset
	}

	if af > AFMax {
		af = AFMax
	}

	var sarEp float64
	if isShort {
		sarEp = sar - ep
	} else {
		sarEp = ep - sar
	}
	var afSarEp = af * sarEp

	p.psarSeries = append(p.psarSeries, &PsarPeriod{
		ChartBar: chartBar,
		EP:       ep,
		AF:       af,
		Sar:      sar,
		SarEP:    sarEp,
		AFSarEP:  afSarEp,
		extVal:   extVal, // aka extLow or extHigh
	})
}

// Next increments the index counter
func (p *Psar) Next() bool {
	if p.Series == nil {
		return false
	}

	// check if the series is empty or at the end
	if n := len(p.Series); n == 0 || p.i > n-1 {
		return false
	}

	// calculate and count
	p.calculatePsar()
	p.i++

	return true
}

// Step "next"s for the number of j from it's current index i. This is
// primarily just for QOL and unit testing purposes
func (p *Psar) Step(j int) {
	j = p.i + j
	// FIXME check if j exceeds the series count and truncate to the series
	// length
	// FIXME per the above, should this return an error to notifiy of exceeded
	// lengths?
	for p.i < j {
		p.Next()
	}
}

// NextSession returns a PsarBar for "tomorrow". This is for a future value
// calculation. It should be noted the Psar you get for "today" is the Psar is
// calculated on data from "yesterday". So at the EOD you can calculate
// tomorrows Psar before the market opens.
// NOTE this is the "next" bar not necessarily the "future" bar. It will be the
// next session from the current iteration index.
// FIXME why are we returning an error here?
func (p *Psar) NextSession() (*PsarPeriod, error) {
	var (
		prevpb = p.psarSeries[p.i-1]
		extVal = prevpb.extVal

		sar float64
	)

	if p.Direction == IsShort {
		sar = prevpb.Sar - prevpb.AFSarEP

		if sar < extVal {
			sar = extVal + p.pipOffset
		}
	} else {
		sar = prevpb.Sar + prevpb.AFSarEP

		if sar > extVal {
			sar = extVal - p.pipOffset
		}
	}

	return &PsarPeriod{Sar: sar}, nil
}
