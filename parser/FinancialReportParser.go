package parser

import (
	"bufio"
	"container/list"
	"encoding/xml"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	"os"
	"strconv"
	"strings"
	"time"
)

const contextTag = "context"
const revenueTag = "Revenues"

const costsAndExpensesTag = "CostsAndExpenses"
const operatingExpensesTag = "OperatingExpenses"
const netIncomeTag = "NetIncomeLoss"

const shortFormDate = "2006-01-02"

type FinancialReportParser struct {
	// TODO add in year/quarter so we can verify that we are parsing the right file
	currentContext, xbrlFileName string
	financialReport              *filings.FinancialReport
	persister                    persist.PersistFinancialReports
}

type XbrlElementParser func(frp *FinancialReportParser, listOfElementLists *list.List)

var parseFunctionMap map[string]XbrlElementParser
var xmlTagToFieldMap map[string]*int64

func initializeParseFunctionMap() {
	parseFunctionMap = map[string]XbrlElementParser{
		contextTag: parseContext,
		revenueTag: parseInt64Field,
		costsAndExpensesTag: parseInt64Field,
		operatingExpensesTag: parseInt64Field,
		netIncomeTag: parseInt64Field,
	}
}

func initializeXmlTagToFieldMap(parser *FinancialReportParser) {
	// there are potentially multiple possible tags for the same field
	xmlTagToFieldMap = map[string]*int64{
		revenueTag: &parser.financialReport.Revenue,
		costsAndExpensesTag: &parser.financialReport.OperatingExpense,
		operatingExpensesTag: &parser.financialReport.OperatingExpense,
		netIncomeTag: &parser.financialReport.NetIncome,
	}
}

func NewFinancialReportParser(xbrlFileName string, fr *filings.FinancialReport, persister persist.PersistFinancialReports) *FinancialReportParser {
	initializeParseFunctionMap()

	parser := &FinancialReportParser{xbrlFileName: xbrlFileName}
	parser.financialReport = fr
	parser.persister = persister

	initializeXmlTagToFieldMap(parser)

	return parser
}

func (frp *FinancialReportParser) GetFinancialReport() *filings.FinancialReport {
	return frp.financialReport
}

func (frp *FinancialReportParser) Parse() {
	xbrlFile, fileErr := os.Open(frp.xbrlFileName)
	fileReader := bufio.NewReader(xbrlFile)

	if fileErr != nil {
		log.Println("Failed to read file")
	} else {
		decoder := xml.NewDecoder(fileReader)

		// create a map of parent elements which we can match up with functions that do the actual parsing
		parserMap := createParserMap(decoder)

		// context must be parsed first!
		contextList, contextExists := parserMap[contextTag]

		if !contextExists {
			// don't bother doing anything else
			log.Error("No context in the following report: ", frp.xbrlFileName)

			return
		}

		parseContext(frp, contextList)
		delete(parserMap, contextTag)

		for mainElementName, list := range parserMap {

			parseFunction, funcExists := parseFunctionMap[mainElementName]

			if funcExists {
				parseFunction(frp, list)
			}
		}

		frErr := frp.financialReport.IsValid()

		if frErr == nil {
			frp.persister.CreateFinancialReport(frp.financialReport)			
		} else {
			log.Error("FinancialReport with CIK <", frp.financialReport.CIK, 
				"> year <", frp.financialReport.Year, 
				"> quarter <", frp.financialReport.Quarter, "> is mising fields: ", frErr)
		}

	}
}

func createParserMap(decoder *xml.Decoder) map[string]*list.List {
	var parentElement string = ""
	elementList := list.New()
	parserMap := make(map[string]*list.List)

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		/*
			    Pseudo code for the algorithm below:
					check if parent is xbrl, ignore if it is
					if not xbrl, save the start element name
					save every element: start, chardata, and endelement, to a new list (DON'T use map, order is not guaranteed)
					if you encounter an end element that matches the start element, add the element list to the object parser for that element

					now we have a map of elements with their variables.
		*/

		switch element := t.(type) {
		case xml.StartElement:
			if element.Name.Local == "xbrl" {
				//no-op
			} else if parentElement == "" {
				parentElement = element.Name.Local
				// include the first parent element so we have access to the attributes
				elementList.PushBack(element)
			} else {
				elementList.PushBack(element)
			}

			break
		case xml.CharData:
			elementList.PushBack(string(element))
			break
		case xml.EndElement:
			if element.Name.Local == "xbrl" {
				//no-op
			} else if element.Name.Local == parentElement {
				if parserMap[parentElement] == nil {
					parserMap[parentElement] = list.New()
				}

				parserMap[parentElement].PushBack(elementList)

				parentElement = ""
				elementList = list.New()
			} else {
				elementList.PushBack(element)
			}

			break
		}
	}

	return parserMap
}

func verifyContext(correctContext string, attributes []xml.Attr) bool {
	for _, attribute := range attributes {
		if attribute.Name.Local == "contextRef" {
			if attribute.Value == correctContext {
				return true
			}
		}
	}

	return false
}

func parseInt64Field(frp *FinancialReportParser, listOfElementLists *list.List) {
	for listElement := listOfElementLists.Front(); listElement != nil; listElement = listElement.Next() {
		elementList := listElement.Value.(*list.List)

		isCorrectContext := false
		var fieldToUpdate *int64
		var tagName string

		for elementListElement := elementList.Front(); elementListElement != nil; elementListElement = elementListElement.Next() {
			xmlElement := elementListElement.Value

			switch element := xmlElement.(type) {
			case xml.StartElement:
				tagName = element.Name.Local
				filingField, fieldExists := xmlTagToFieldMap[tagName]
				if fieldExists {
					isCorrectContext = verifyContext(frp.currentContext, element.Attr)
					if isCorrectContext {
						fieldToUpdate = filingField
					}
				}

				break
			case string:
				if isCorrectContext {
					fieldStr := strings.TrimSpace(element)
					int64Field, convErr := strconv.ParseInt(fieldStr, 10, 64)

					if convErr != nil {
						log.Error("Failed parsing int64 field <", tagName, "> into an int due to error: ", convErr)
					} else {
						*fieldToUpdate = int64Field
					}
				}

				break
			}
		}
	}
}

func parseContext(frp *FinancialReportParser, listOfElementLists *list.List) {
	// TODO I would love a way to not have to copy/paste this loop and switch statement in every parsing function...

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
				if element.Name.Local == contextTag {
					for _, attribute := range element.Attr {
						if attribute.Name.Local == "id" {
							currentContext = attribute.Value

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
				if parsingStartDate {
					startDate, _ = time.Parse(shortFormDate, content)
				} else if parsingEndDate {
					endDate, _ = time.Parse(shortFormDate, content)
				}
				break
			case xml.EndElement:
				parsingStartDate = false
				parsingEndDate = false
				break
			}
		}

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
