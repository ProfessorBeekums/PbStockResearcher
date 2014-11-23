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

	    // TODO hack
	    parseContext := false

		for { 
		    // Read tokens from the XML document in a stream. 
		    t, _ := decoder.Token() 
		    if t == nil { 
		        break 
		    }

		    // TODO break out the map builder into a different function and call it here
		    // TODO pseudo code for what this should look like:
		    /*
				check if parent is xbrl, ignore if it is
				if not xbrl, save the start element name
				save every element: start, chardata, and endelement, to a new list (DON'T use map, order is not guaranteed)
				if you encounter an end element that matches the start element, add the element list to the object parser for that element

				now we have a map of elements with their variables. 
				parse all contexts first so that we know which one we want. store that in the FinancialReportParser

				parse every other element
		    */

			switch element := t.(type) { 
			    case xml.StartElement: 
			    	log.Println("start tagf:", element.Name.Local)
			    // TODO there are going to be many versions of this. need to parse out contexts and figure out which ones are ok based on the start date. 
			    // note that multiple contexts can have the same start and end date. we also want the latest quarter, not 6 month period
			    	if element.Name.Local == "Revenues" || element.Name.Local == "endDate" || element.Name.Local == "startDate" {
			    		parseData = true
			    	}

			    	for _, attribute := range element.Attr {
			    		if attribute.Name.Local == "id" && attribute.Value == "Context_6ME__30-Sep-2013_FinancingReceivableRecordedInvestmentByClassOfFinancingReceivableAxis_ConsumerLoanMember" {
			    			log.Println("parsing a context start")
			    			parseContext = true
			    		}
			    	}

			    	if parseData {
				    	log.Println("Space:",element.Name.Space)
				    	log.Println("LOcal:", element.Name.Local)
				    	log.Println("Attr", element.Attr)
			    	}
			    	break
			    case xml.CharData:
			    	if parseData {
			    		log.Println(string(element))
			    	}
			    	break
			    case xml.EndElement:
			    	log.Println(parseContext)
			    	return
			    	// if parseContext && element.Name.Local == "context" {
			    	// 	log.Println("End context")
			    	// 	parseContext = false
			    	// }
			    	// if parseData {
			    	// 	log.Println("end parsig: ", element.Name.Local)
			    	// }
			    	// parseData = false
			    	// break
			}
		}
	}
}