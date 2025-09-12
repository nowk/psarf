package psarf

import (
	"math"
	"testing"
	"time"
)

type myChartBar struct {
	high float64
	low  float64
	date time.Time
}

var _ = &myChartBar{}

func (c *myChartBar) High() float64 {
	return c.high
}

func (c *myChartBar) Low() float64 {
	return c.low
}

func (c *myChartBar) Date() *time.Time {
	return &c.date
}

func (c *myChartBar) Open() float64 {
	return 0.0
}

func (c *myChartBar) Close() float64 {
	return 0.0
}

var now = time.Now()

func TestPsar_EP_ReturnsTheHighestHigh(t *testing.T) {
	p := &Psar{
		Series: []ChartBar{
			&myChartBar{date: now, high: 10.0},
		},
		StartDate: &now,
	}
	p.Next()
	pb := p.Bar()
	if got := pb.EP; got != 10.0 {
		t.Fatalf("expected 10.0, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 9.0},
			&myChartBar{date: now, high: 10.0},
		},
		StartDate: &now,
	}
	p.Step(2)
	pb = p.Bar()
	if got := pb.EP; got != 10.0 {
		t.Fatalf("expected 10.0, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 10.0},
			&myChartBar{date: now, high: 9.0},
		},
		StartDate: &now,
	}
	p.Step(2)
	pb = p.Bar()
	if got := pb.EP; got != 10.0 {
		t.Fatalf("expected 10.0, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 10.0},
			&myChartBar{date: now, high: 8.0},
			&myChartBar{date: now.Add(24 * time.Hour), high: 9.0},
		},
		StartDate: &now,
	}
	p.Step(3)
	pb = p.Bar()
	if got := pb.EP; got != 10.0 {
		t.Fatalf("expected 10.0, got %v", got)
	}
}

func TestPsar_AF_IncreasesWithEachIncreasingEP(t *testing.T) {
	p := &Psar{
		Series: []ChartBar{
			&myChartBar{high: 10.0},
		},
		StartDate: &now,
	}
	p.Next()
	pb := p.Bar()
	if got := pb.AF; got != 0.0 {
		t.Fatalf("expected 0.0, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 9.0},
			&myChartBar{date: now, high: 10.0},
		},
		StartDate: &now,
	}
	p.Step(2)
	pb = p.Bar()
	if got := pb.AF; got != 0.02 {
		t.Fatalf("expected 0.02, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 10.0},
			&myChartBar{date: now, high: 9.0},
		},
		StartDate: &now,
	}
	p.Step(2)
	pb = p.Bar()
	if got := pb.AF; got != 0.02 {
		t.Fatalf("expected 0.02, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 8.0},
			&myChartBar{date: now, high: 9.0},
			&myChartBar{date: now.Add(24 * time.Hour), high: 10.0},
		},
		StartDate: &now,
	}
	p.Step(3)
	pb = p.Bar()
	if got := pb.AF; got != 0.04 {
		t.Fatalf("expected 0.04, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 10.0},
			&myChartBar{date: now, high: 9.0},
			&myChartBar{date: now.Add(24 * time.Hour), high: 8.0},
		},
		StartDate: &now,
	}
	p.Step(3)
	pb = p.Bar()
	if got := pb.AF; got != 0.02 {
		t.Fatalf("expected 0.02, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 10.0},
			&myChartBar{date: now, high: 9.0},
			&myChartBar{date: now.Add(24 * time.Hour), high: 10.0},
		},
		StartDate: &now,
	}
	p.Step(3)
	pb = p.Bar()
	if got := pb.AF; got != 0.02 {
		t.Fatalf("expected 0.02, got %v", got)
	}
}

// type mockPsar struct {
// 	Psar
// 	afFunc func() float64
// }

// func (m *mockPsar) af(i int) float64 {
// 	return m.afFunc()
// }

// // FIXME Write this test via a mock.
// func TestPsar_AF_DoesNotExceedMax(t *testing.T) {
// 	t.Skipf("How to test with a mock?")
// 	p := &mockPsar{
// 		Psar: Psar{
// 			Series: []ChartBar{
// 				&myChartBar{high: 9.0},
// 				&myChartBar{high: 10.0},
// 			},
// 		},
// 		afFunc: func() float64 {
// 			return 0.18
// 		},
// 	}
// 	p.Next()
// 	if got := rn2(p.AF); got != 0.20 {
// 		t.Fatalf("expected 0.20, got %v", got)
// 	}
// }

func TestPsar_Sar(t *testing.T) {
	enDate := time.Date(2021, 0, 5, 0, 0, 0, 0, time.UTC)

	// Based off of ROIC 01/15/2021 entry bar on 01/08/2021
	p := &Psar{
		Series: []ChartBar{
			&myChartBar{high: 13.23, low: 12.79},
		},
		StartDate: &enDate,
	}
	p.Next()
	pb := p.Bar()
	if got := pb.Sar; got != 12.79 {
		t.Fatalf("expected 12.79, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 13.23, low: 12.79},
			&myChartBar{high: 13.73, low: 13.04},
			&myChartBar{high: 13.66, low: 13.23},
			&myChartBar{date: enDate, high: 13.94, low: 13.4},
		},
		StartDate: &enDate,
	}
	p.Step(4)
	pb = p.Bar()
	if got := pb.Sar; got != 12.79 {
		t.Fatalf("expected 12.79, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 13.23, low: 12.79},
			&myChartBar{high: 13.73, low: 13.04},
			&myChartBar{high: 13.66, low: 13.23},
			&myChartBar{date: enDate, high: 13.94, low: 13.4},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 13.77, low: 13.38},
		},
		StartDate: &enDate,
	}
	p.Step(5)
	pb = p.Bar()
	if got := rn2(pb.Sar); got != 12.81 {
		t.Fatalf("expected 12.81, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 13.23, low: 12.79},
			&myChartBar{high: 13.73, low: 13.04},
			&myChartBar{high: 13.66, low: 13.23},
			&myChartBar{date: enDate, high: 13.94, low: 13.4},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 13.77, low: 13.38},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 14.00, low: 13.34},
		},
		StartDate: &enDate,
	}
	p.Step(6)
	pb = p.Bar()
	if got := rn2(pb.Sar); got != 12.84 {
		t.Fatalf("expected 12.84, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 13.23, low: 12.79},
			&myChartBar{high: 13.73, low: 13.04},
			&myChartBar{high: 13.66, low: 13.23},
			&myChartBar{date: enDate, high: 13.94, low: 13.4},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 13.77, low: 13.38},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 14.00, low: 13.34},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 14.32, low: 13.84},
		},
		StartDate: &enDate,
	}
	p.Step(7)
	pb = p.Bar()
	if got := rn2(pb.Sar); got != 12.88 {
		t.Fatalf("expected 12.88, got %v", got)
	}
}

func TestPsar_Sar_Short(t *testing.T) {
	entryBar, _ := time.Parse("2006-01-02", "2025-08-26")

	// based off of data from TAP from Aug 22, 2025
	// NOTE TV data vs GPT is different, below is TV data
	barSeries := []ChartBar{
		&myChartBar{high: 52.94, low: 51.81},
		&myChartBar{high: 52.81, low: 51.44},
		&myChartBar{date: entryBar, high: 51.43, low: 50.25},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 51.29, low: 50.27},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 51.39, low: 49.79},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.94, low: 50.14},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.40, low: 49.42},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.22, low: 49.63},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.06, low: 49.50},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.66, low: 49.42},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.19, low: 49.46},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 50.31, low: 49.30},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 49.25, low: 48.53},
		&myChartBar{date: entryBar.Add(24 * time.Hour), high: 49.42, low: 48.85},
	}

	expectedPsar := []float64{
		52.94,
		52.94,
		52.94, // starting for the entry bar
		52.89,
		52.83,
		52.71,
		52.59,
		52.40,
		52.23,
		52.06,
		// 51.85, // TV's psar has a reset on the AF which slows the descent
		// 51.65, // question is, do I want that to slow down the acceleration?
		// 51.42,
		// 51.07,
		51.90, // Calculations for the original formula, also confimed on TradeStation
		51.75,
		51.55,
		51.25,
	}

	p := &Psar{
		Series:    barSeries,
		StartDate: &entryBar,
		Direction: IsShort,
	}

	// TODO when normal iteratation is called on Psar, ie p.Next(), the Psar
	// should start athe "StartDate" bar, not the first bar in the series.

	// p.Step(2) // move to the entry bar - 1, we are calling Next
	for i, want := range expectedPsar {
		if !p.Next() {
			t.Fatalf("Next() failed at index %d", i)
		}
		got := rn2(p.Bar().Sar)
		if got != want {
			t.Errorf("bar %d: expected SAR %v, got %v", i, want, got)
		}
	}
}

func TestPsar_Sar_CannotExceedLowOfLast2Bars(t *testing.T) {
	// Based on data from BIG Jan 04, 2021
	enDate := time.Date(2021, 0, 5, 0, 0, 0, 0, time.UTC)

	p := &Psar{
		Series: []ChartBar{
			&myChartBar{high: 43.57, low: 42.05},
			&myChartBar{date: enDate, high: 44.12, low: 42.09},
		},
		StartDate: &enDate,
	}
	p.Step(2)
	pb := p.Bar()
	if got := pb.Sar; got != 42.05 {
		t.Fatalf("expected 42.05, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 43.57, low: 42.05},
			&myChartBar{date: enDate, high: 44.12, low: 42.09},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 46.0, low: 43.45},
		},
		StartDate: &enDate,
	}
	p.Step(3)
	pb = p.Bar()
	if got := rn2(pb.Sar); got != 42.05 {
		t.Fatalf("expected 42.05, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 43.57, low: 42.05},
			&myChartBar{date: enDate, high: 44.12, low: 42.09},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 46.0, low: 43.45},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 45.63, low: 43.74},
		},
		StartDate: &enDate,
	}
	p.Step(4)
	pb = p.Bar()
	if got := rn2(pb.Sar); got != 42.09 {
		t.Fatalf("expected 42.09, got %v", got)
	}

	p = &Psar{
		Series: []ChartBar{
			&myChartBar{high: 43.57, low: 42.05},
			&myChartBar{date: enDate, high: 44.12, low: 42.09},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 46.0, low: 42.08},
			&myChartBar{date: enDate.Add(24 * time.Hour), high: 45.63, low: 43.74},
		},
		StartDate: &enDate,
	}
	p.Step(4)
	pb = p.Bar()
	if got := rn2(pb.Sar); got != 42.08 {
		t.Fatalf("expected 42.08, got %v", got)
	}
}


func rn(f float64, p float64) float64 {
	return math.Round(f*p) / p
}

func rn2(f float64) float64 {
	return rn(f, 100)
}
