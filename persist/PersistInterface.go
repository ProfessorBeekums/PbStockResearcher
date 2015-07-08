package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
)

type PersistFinancialReports interface {
	InsertUpdateFinancialReport(fr *filings.FinancialReport)
	// TODO unused for now
	//	GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport
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
