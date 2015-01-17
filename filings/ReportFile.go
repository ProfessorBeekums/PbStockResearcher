package filings

type ReportFile struct {
	CIK, Year, Quarter int64
	Filepath           string
	FormType           string
	Parsed             bool
	ParseError         bool
}

//func (rf *ReportFile) IsFinancialReport() {

//}
