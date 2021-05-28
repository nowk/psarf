package psarf

// PsarPeriod represents the calculated Psar values for a given day
type PsarPeriod struct {
	// implements ChartBar this is mainly just as a reference to the original
	// bar data
	ChartBar

	// naming schemes are loosly based on the formula itself.
	EP      float64
	AF      float64
	Sar     float64 // This is the Sar value, the one you want
	SarEP   float64
	AFSarEP float64

	// extLow is the extreme low of the last 2 lows
	// NOTE this is here for performance. This is a "cached" value. Lookups on
	// any dataset of more than a few bars becomes noticeable, and over a few
	// weeks it becomes nearly impossible
	extLow float64
}

// ExtLow returns the extLow value (normally not a display value, but available
// this way just in case)
func (p *PsarPeriod) ExtLow() float64 {
	return p.extLow
}
