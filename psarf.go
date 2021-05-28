// Package psarf provides use of the Psar (Parabolic Sar) formula to use as a
// trailing stop mechanism
//
// In it's current implementation psarf has a few limitions. Noted below will be
// the TODO for items that are in the works but significant and to be noted
//
// TODO
// - [ ] Enable Short side calculations (currently it's long side only)
// - [ ] Support for a variety of timeframes (currently it's only viable on
//       daily timeframes due to Psar#isEntryBar calculation)
//
package psarf
