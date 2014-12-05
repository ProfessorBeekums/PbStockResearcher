package scraper

import(
	"testing"
)

func TestIsXbrlFileMatch(t *testing.T) {
	possibleMatches := []string{"nick-20121231.xml",
		"ck0001000069-20120930.xml",
	}
	badMatches := []string{"nick-20121231.xsd",
		"nick-20121231_cal.xml",
		"ck0001000069-20120930_cal.xml",
		"ck0001000069-20120930_def.xml",
		"ck0001000069-20120930_lab.xml",
		"ck0001000069-20120930_pre.xml",
		"",
	}

	for _, goodMatch := range possibleMatches {
		isMatch := isXbrlFileMatch(goodMatch)
		if isMatch == false {
			t.Fail()
			t.Log("Should have been a good match: ", goodMatch)
		}
	}

	for _, badMatch := range badMatches {
        isMatch := isXbrlFileMatch(badMatch)
        if isMatch != false {
			t.Fail()
            t.Log("Should have been a bad match: ", badMatch)
        }
    }
}
