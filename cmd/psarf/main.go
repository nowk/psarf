package main

import (
	"fmt"
	"log"
	"time"

	// "github.com/nowk/psarf"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func main() {

	symbol := "TSN"
	startDate := time.Date(2025, time.July, 9, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.July, 16, 0, 0, 0, 0, time.UTC)

	start := datetime.New(&startDate)
	end := datetime.New(&endDate)

	params := &chart.Params{
		Symbol:   symbol,
		Start:    start,
		End:      end,
		Interval: datetime.OneDay,
	}

	iter := chart.Get(params)
	for iter.Next() {
		bar := iter.Bar()
		fmt.Println(bar)
		// fmt.Printf("%d: Open %.2f, High %.2f, Low %.2f, Close %.2f, Volume %d\n",
		// 	bar.Timestamp,
		// 	bar.Open,
		// 	bar.High,
		// 	bar.Low,
		// 	bar.Close,
		// 	bar.Volume,
		// )
	}
	if err := iter.Err(); err != nil {
		log.Fatalf("Error getting chart data: %v", err)
	}

	// psar := &psarf.Psar{
	// 	Series:    nil,
	// 	StartDate: nil,
	// }
}
