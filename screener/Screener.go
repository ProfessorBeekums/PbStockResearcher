package screener

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
)

// Used to screen data for a single quarter
type SingleQuarterFilingsScreener struct {
}

// sample query for mongodb:
// db.FinancialReport.aggregate([
//   {$project: {cik: "$cik", test: { $divide: ["$netincome", "$revenue"]}}},
//   {$match : {test : { $gt: .5, $lt : 1}}}
// ])
