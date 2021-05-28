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

	// psarSeries holds the calculated Psar values for each bar in the series.
	psarSeries []*PsarPeriod

	// pipOffset is a value to offset the sar for the initial Sar and any Sar
	// that is truncated by the lowest last 2 bars rule
	// NOTE this is not part of the psar formula, but is for my own personal use
	pipOffset float64

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
		extLow   = chartBar.Low()
		ep       = chartBar.High()
		af       = AFIncrement

		sar float64
	)

	if i > 0 {
		var prevbp = p.psarSeries[i-1]
		af = prevbp.AF // FIXME write unit test for this

		if ep > prevbp.EP {
			af += AFIncrement
		} else {
			ep = prevbp.EP
		}

		// the sar
		sar = prevbp.Sar + prevbp.AFSarEP

		// get the lowest low (long) of the last 2 bars. This is a cache for
		// performance purposes.
		// FIXME pretty sure this can be reduced to not have to go back 2 bars
		var l1 = prevbp.Low()
		extLow = l1
		if i-2 >= 0 {
			if l2 := p.psarSeries[i-2].Low(); l2 < l1 {
				extLow = l2
			}
		}
	}
	// set the initial sar and af only on the entry bar. Calculate in any
	// additional pipOffset if applicable
	if p.isEntryBar(i) {
		af, sar = AFIncrement, p.Series[0].Low()-p.pipOffset
	}
	// "reset" af and sar values for bars before the start date. Any bar between
	// the first bar in the series and the start bar are defaulted to the first
	// bar in the series. They are not used in the calculation other than using
	// their lows as part of the low of the last 2 bars values.
	if before(*chartBar.Date(), *p.StartDate) {
		af, sar = 0.0, p.Series[0].Low()
	}
	// check to see if that the sar does not exceed (long) the lowest low of the
	// last 2 bars. *Tricky one, not explained in many of the Psar
	// explanations.*
	if sar > extLow {
		sar = extLow - p.pipOffset
	}
	// acceleration factor should not exceed the max af (0.20)
	// TODO needs a unit test
	if af > AFMax {
		af = AFMax
	}

	var (
		sarEp   = ep - sar
		afSarEp = af * sarEp
	)
	p.psarSeries = append(p.psarSeries, &PsarPeriod{
		ChartBar: chartBar,

		EP:      ep,
		AF:      af,
		Sar:     sar,
		SarEP:   sarEp,
		AFSarEP: afSarEp,

		extLow: extLow,
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

// Step "next"s for the number of j from it's current index i. This is primarily
// just for QOL and unit testing purposes
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
func (p *Psar) NextSession() (*PsarPeriod, error) {
	var (
		prevpb = p.psarSeries[p.i-1]
		pb     *PsarPeriod
	)

	// calculate the sar only here (this is why we don't call calculate as we
	// don't have any of the bar data for "tomorrow")
	// TODO is there a way to integrate this into the calculation as a whole
	// the extreme low of last 2 days was missed intially
	var sar = prevpb.Sar + prevpb.AFSarEP

	// must always observer the lowest 2 lows rule
	if sar > prevpb.extLow {
		sar = prevpb.extLow - p.pipOffset
	}

	pb = &PsarPeriod{Sar: sar}

	return pb, nil
}
