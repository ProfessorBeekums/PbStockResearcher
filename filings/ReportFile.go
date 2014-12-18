package filings 

type ReportFile struct {
	CIK, Year, Quarter int64
	Filepath string
	Parsed bool
}


