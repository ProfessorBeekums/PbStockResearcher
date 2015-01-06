package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
)

type PersistFinancialReports interface {
	CreateFinancialReport(fr *filings.FinancialReport)
	UpdateFinancialReport(fr *filings.FinancialReport)
	GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport
}

type PersistFinancialReportsRaw interface {
	InsertUpdateRawReport(rawReport *filings.FinancialReportRaw)
	GetRawReport(cik, year, quarter int64) *filings.FinancialReportRaw
}

type PersistReportFiles interface {
	InsertUpdateReportFile(reportFile *filings.ReportFile)
	GetNextUnparsedFiles(numToGet int64) *[]filings.ReportFile
}

type PersistCompany interface {
	InsertUpdateCompany(company *filings.Company)
	GetCompany(cik int64) *filings.Company
}
