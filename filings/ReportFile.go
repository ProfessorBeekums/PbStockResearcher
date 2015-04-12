package filings

import "strconv"

type ReportFile struct {
	ReportFileId       int64
	CIK, Year, Quarter int64
	Filepath           string
	FormType           string
	Parsed             bool
	ParseError         bool
}

func (rf *ReportFile) GetLogStr() string {
	return "report file cik <" + strconv.FormatInt(rf.CIK, 10) +
		"> with year <" + strconv.FormatInt(rf.Year, 10) +
		"> and quarter <" + strconv.FormatInt(rf.Quarter, 10) +
		"> and file path <" + rf.Filepath + ">"
}

//func (rf *ReportFile) IsFinancialReport() {

//}
