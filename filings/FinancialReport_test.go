package filings

import (
	"testing"
)

func TestGetPreviousQuarter(t *testing.T) {
	fr := &FinancialReport{Year: 2013, Quarter: 4}

	year, quarter := fr.GetPreviousQuarter()

	if year != 2013 {
		t.Fatal("GetPreviousQuarter expected year 2013, got: ", year)
	}

	if quarter != 3 {
		t.Fatal("GetPreviousQuarter expected quarter 3, got: ", quarter)
	}

	fr = &FinancialReport{Year: 2013, Quarter: 1}

	year, quarter = fr.GetPreviousQuarter()

	if year != 2012 {
		t.Fatal("GetPreviousQuarter expected year 2012, got: ", year)
	}

	if quarter != 4 {
		t.Fatal("GetPreviousQuarter expected quarter 4, got: ", quarter)
	}
}