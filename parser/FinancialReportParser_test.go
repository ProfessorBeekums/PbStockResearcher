package parser

import (
	"bufio"
	"encoding/xml"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"os"
	"testing"
)

func TestCreateParserMap(t *testing.T) {
	xbrlFile, _ := os.Open("../testData/TestCreateParserMap.xml")
	fileReader := bufio.NewReader(xbrlFile)
	decoder := xml.NewDecoder(fileReader)

	parserMap := createParserMap(decoder)

	if len(parserMap) != 6 {
		t.Fatal("Expected 6 elements in parser map! Got: ", len(parserMap))
	}

	parentElements := [...]string{"context", "SalesRevenueNet", "FranchiseRevenue", "Revenues", "CostOfRevenue", "schemaRef"}

	for _, parentElement := range parentElements {
		_, parentExists := parserMap[parentElement]
		if parentExists == false {
			t.Fatal("Missing element from parser map: ", parentElement)
		}
	}
}

func TestVerifyContext(t *testing.T) {
	context := "myContext"

	xmlName1 := xml.Name{Local: "foo"}
	xmlName2 := xml.Name{Local: "contextRef"}

	attr1 := xml.Attr{Name: xmlName1, Value: "foo"}
	attr2 := xml.Attr{Name: xmlName2, Value: context}

	var attributes = []xml.Attr{attr1, attr2}

	exists := verifyContext(context, attributes)

	if exists == false {
		t.Fatal("Context should exist in attributes: ", context)
	}

	missingContext := "missing"

	exists = verifyContext(missingContext, attributes)

	if exists {
		t.Fatal("Context should not exists: ", missingContext)
	}
}

func TestParseInt64Field(t *testing.T) {
	xbrlFile, _ := os.Open("../testData/TestCreateParserMap.xml")
	fileReader := bufio.NewReader(xbrlFile)
	decoder := xml.NewDecoder(fileReader)

	parserMap := createParserMap(decoder)

	frp := NewFinancialReportParser("../testData/TestCreateParserMap.xml", &filings.FinancialReport{}, nil)
	frp.currentContext = "c0_From1Jun2014To31Aug2014"

	parseInt64Field(frp, parserMap[revenueTag])

	if frp.financialReport.Revenue != 9457448 {
		t.Fatal("Expected revenue was 9457448, received: ", frp.financialReport.Revenue)
	}
}

func TestParseContext(t *testing.T) {
	xbrlFile, _ := os.Open("../testData/TestCreateParserMap.xml")
	fileReader := bufio.NewReader(xbrlFile)
	decoder := xml.NewDecoder(fileReader)

	parserMap := createParserMap(decoder)

	frp := &FinancialReportParser{}

	parseContext(frp, parserMap[contextTag])

	if frp.currentContext != "c0_From1Jun2014To31Aug2014" {
		t.Fatal("Expected context was c0_From1Jun2014To31Aug2014, received: ",
			frp.currentContext)
	}
}

func TestCompleteParse(t *testing.T) {
	frp := NewFinancialReportParser("../testData/TestCreateParserMap.xml", &filings.FinancialReport{}, &MockPersister{})

	frp.Parse()

	fr := frp.financialReport

	if fr.Revenue != 9457448 {
		t.Fatal("Expected revenue was 9457448, received: ", fr.Revenue)
	}
}

type MockPersister struct{}

func (mp *MockPersister) CreateFinancialReport(fr *filings.FinancialReport) {}
func (mp *MockPersister) UpdateFinancialReport(fr *filings.FinancialReport) {}
func (mp *MockPersister) GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport {
	return nil
}
