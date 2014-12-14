package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
)

type PersistFinancialReports interface {
	CreateFinancialReport(fr *filings.FinancialReport) error
	UpdateFinancialReport(fr *filings.FinancialReport) error
	GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport
}
