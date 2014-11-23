package parser

import (
	"bufio"
	"encoding/xml"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"os"
)

type FinancialReportParser struct {
	// TODO add in year/quarter so we can verify that we are parsing the right file
	xbrlFileName string
}

func NewFinancialReportParser(xbrlFileName string) *FinancialReportParser {
	return &FinancialReportParser{xbrlFileName: xbrlFileName}
}

func (frp *FinancialReportParser) Parse() {
	xbrlFile, fileErr := os.Open(frp.xbrlFileName)
	fileReader := bufio.NewReader(xbrlFile)

	if fileErr != nil {
		log.Println("Failed to read file")
	} else {
		// var xbrlDict string
		// xmlErr := xml.Unmarshal(xbrlFileBytes, &xbrlDict)

		decoder := xml.NewDecoder(fileReader) 

	    parseData := false

		for { 
		    // Read tokens from the XML document in a stream. 
		    t, _ := decoder.Token() 
		    if t == nil { 
		        break 
		    }

			switch element := t.(type) { 
			    case xml.StartElement: 
			    // TODO there are going to be many versions of this. need to parse out contexts and figure out which ones are ok based on the start date. 
			    // note that multiple contexts can have the same start and end date. we also want the latest quarter, not 6 month period
			    	if element.Name.Local == "Revenues" {
			    		parseData = true
			    	}

			    	if parseData {
				    	log.Println(element.Name.Space)
				    	log.Println(element.Name.Local)
				    	log.Println(element.Attr)
			    	}
			    	break
			    case xml.CharData:
			    	if parseData {
			    		log.Println(string(element))
			    	}
			    	break
			    case xml.EndElement:
			    	parseData = false
			    	break
			}
		}
	}
}