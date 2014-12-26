package filings

import "errors"

type FinancialReport struct {
	CIK, Year, Quarter int64
	Revenue, OperatingExpense, NetIncome  int64
	CurrentAssets, TotalAssets, CurrentLiabilities, TotalLiabilities int64
}

func (fr *FinancialReport) IsValid() error {
	missingFields := ""

	// not looping through every struct field because some may be optional
	if(fr.Revenue == 0) {
		missingFields += "Revenue,"
	}

	if(fr.OperatingExpense == 0) {
		missingFields += "OperatingExpense,"
	}

	if(fr.NetIncome == 0) {
		missingFields += "NetIncome,"
	}


	// if(fr.TotalAssets == 0) {
	// 	missingFields += "TotalAssets,"
	// }


	// if(fr.TotalLiabilities == 0) {
	// 	missingFields += "TotalLiabilities,"
	// }


	if len(missingFields) > 0 {
		return errors.New(missingFields) 
	} else {
		return nil
	}
}