package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
)

type PersistFinancialReports interface {
	CreateFinancialReport(fr *filings.FinancialReport)
	UpdateFinancialReport(fr *filings.FinancialReport)
	GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport
}

type PersistReportFiles interface {
	InsertUpdateReportFile(reportFile *filings.ReportFile)
	GetNextUnparsedFiles(numToGet int64) *[]filings.ReportFile
}

type PersistCompany interface {
	InsertUpdateCompany(company *filings.Company)
	GetCompany(cik int64) *filings.Company
}
