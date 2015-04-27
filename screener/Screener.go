package screener

import "github.com/ProfessorBeekums/PbStockResearcher/filings"

// Used to screen data for a single quarter
type SingleQuarterFilingsScreener interface {
	ScreenNetMargin(year, quarter int, min, max float64) map[*filings.Company]float64
	ScreenAssetRatio(year, quarter int, min, max float64) map[*filings.Company]float64
	ScreenCurrentRatio(year, quarter int, min, max float64) map[*filings.Company]float64
}
