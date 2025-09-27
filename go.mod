module github.com/nowk/psarf

go 1.23.11

// Fix for break in piquette/finance-go as of early 2025
// https://github.com/piquette/finance-go/issues/32#issuecomment-2699724078
//
// require github.com/piquette/finance-go v1.1.0
replace github.com/piquette/finance-go => github.com/psanford/finance-go v0.0.0-20250222221941-906a725c60a0

require github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
