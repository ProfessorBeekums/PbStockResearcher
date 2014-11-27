package parser

import (
	"bufio"
	"container/list"
	"encoding/xml"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"os"
	"strings"
	"time"
)

const shortFormDate = "2006-01-02"

type FinancialReportParser struct {
	// TODO add in year/quarter so we can verify that we are parsing the right file
	currentContext, xbrlFileName string
}

type XbrlElementParser func(frp *FinancialReportParser, listOfElementLists *list.List) 

var parseFunctionMap map[string]XbrlElementParser 

func NewFinancialReportParser(xbrlFileName string) *FinancialReportParser {
	parseFunctionMap = map[string]XbrlElementParser {
		"context" : parseContext,
	}

	return &FinancialReportParser{xbrlFileName: xbrlFileName}
}

func (frp *FinancialReportParser) Parse() {
	xbrlFile, fileErr := os.Open(frp.xbrlFileName)
	fileReader := bufio.NewReader(xbrlFile)

	if fileErr != nil {
		log.Println("Failed to read file")
	} else {
		decoder := xml.NewDecoder(fileReader) 

	    var parentElement string = ""
	    elementList := list.New()
	    parserMap := make(map[string]*list.List)

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
			    	if element.Name.Local == "xbrl" {
			    		//no-op
			    	} else if parentElement == "" {
			    		parentElement = element.Name.Local
			    		// log.Println("@@@ Savign poarent ", parentElement)
			    		// include the first parent element so we have access to the attributes
			    		elementList.PushBack(element)
			    	} else {
			    		// log.Println("@@@ Pushing parent; ", parentElement, element.Name.Local)
			    		elementList.PushBack(element)
			    	}

			    	break
			    case xml.CharData:
			    		// log.Println("@@@ Pushing parent; ", parentElement, string(element))
			    	elementList.PushBack(string(element))
			    	break
			    case xml.EndElement:
			    	if element.Name.Local == "xbrl" {
			    		//no-op
			    	} else if element.Name.Local == parentElement {
			    		if parserMap[parentElement] == nil {
			    			parserMap[parentElement] = list.New()
			    		}

			    		// log.Println("@@@ Adding element list ", elementList)
			    		parserMap[parentElement].PushBack(elementList)

			    		parentElement = ""
			    		elementList = list.New()
			    	} else {
			    		// log.Println("@@@ Pushing parent; ", parentElement, element.Name.Local)
			    		elementList.PushBack(element)
			    	}
			    	
			    	break
			}
		}

		log.Println("Our parser map is ", parserMap)

		// context must be parsed first!
		contextList, contextExists := parserMap["context"]

		if !contextExists {
			// don't bother doing anything else
			log.Error("No context in the following report: ", frp.xbrlFileName)

			return
		}

		parseContext(frp, contextList)
		delete(parserMap, "context")

		for mainElementName, list := range parserMap {
			parseFunction, funcExists := parseFunctionMap[mainElementName]

			if funcExists {
				// log.Println("@@@ going to parse: ", mainElementName)
				parseFunction(frp, list)
			}
		}

		log.Println("@@@ context used for this quarter is ", frp.currentContext)
	}
}

func parseContext(frp *FinancialReportParser, listOfElementLists *list.List) {
	// TODO I would love a way to not have to copy/paste this loop and switch statement in every parsing function...
	// var startParsed bool = false
	// var currentStart string
	var currentContext, latestContext string
	var latestEndDate time.Time

	for listElement := listOfElementLists.Front(); listElement != nil; listElement = listElement.Next() {
		elementList := listElement.Value.(*list.List)

		parsingStartDate := false
		parsingEndDate := false

		var startDate time.Time
		var endDate time.Time

		for elementListElement := elementList.Front(); elementListElement != nil; elementListElement = elementListElement.Next() {
			xmlElement := elementListElement.Value

			switch element := xmlElement.(type) { 
				case xml.StartElement: 
				// log.Println("@@@ going to start parse context: ", element.Name.Local)
					if element.Name.Local == "context" {
						for _, attribute := range element.Attr {
				    		if attribute.Name.Local == "id" {
				    			currentContext = attribute.Value

				    			// log.Println("@@@ setting current contest to: ", currentContext)

				    			break
				    		}
				    	}
					} else if element.Name.Local == "startDate" {
						parsingStartDate = true
					} else if element.Name.Local == "endDate" {
						parsingEndDate = true
					}
			    	break
			    case string:
			    	content := strings.TrimSpace(element)
			    	// log.Println("@@@ going to content parse context: ", content)
			    	if parsingStartDate {
			    		startDate, _ = time.Parse(shortFormDate, content)
			    	} else if parsingEndDate {
			    		endDate,_ = time.Parse(shortFormDate, content)
			    	}
			    	break
			    case xml.EndElement:
					// log.Println("@@@ going to end parse context: ", element.Name.Local)
					parsingStartDate = false
					parsingEndDate = false
			    	break
			}
		}

		log.Println("@@@ for contexnt ", currentContext, " we have start ", startDate, " and end ", endDate )

		periodLengthInMonths := int(endDate.Month()) - int(startDate.Month())

		// we only care about the latest quarter for this report
		if periodLengthInMonths == 2 {
			if endDate.Unix() > latestEndDate.Unix() {
				latestEndDate = endDate
				latestContext = currentContext
			}
		}
	}

	frp.currentContext = latestContext
}