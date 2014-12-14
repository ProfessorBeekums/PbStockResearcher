package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
)

type PersistFinancialReports interface {
	CreateFinancialReport(fr *filings.FinancialReport)
	UpdateFinancialReport(fr *filings.FinancialReport)
	GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport
}
