package parser

import (
	"bufio"
	"encoding/xml"
	//"fmt"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	// "github.com/ProfessorBeekums/PbStockResearcher/log"
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

func TestParseContext(t *testing.T) {
	xbrlFile, _ := os.Open("../testData/TestCreateParserMap.xml")
	fileReader := bufio.NewReader(xbrlFile)
	decoder := xml.NewDecoder(fileReader)

	parserMap := createParserMap(decoder)

	frp := &FinancialReportParser{}
	frp.contextMap = make(map[string]*context)

	parseContext(frp, parserMap[contextTag])

	// verify each context as added to the map
	contextName := "c0_From1Jun2014To31Aug2014"
	context, exists := frp.contextMap[contextName]

	if !exists {
		t.Fatal("Missing context: ", contextName)
	} else {
		if context.startDate.Month() != 6 || context.startDate.Year() != 2014 {
			t.Fatal("Context <", contextName, "> has wrong startDate of ", context.startDate.Month(), " - ", context.startDate.Year())
		} else if context.endDate.Month() != 8 || context.endDate.Year() != 2014 {
			t.Fatal("Context <", contextName, "> has wrong endDate of ", context.startDate.Month(), " - ", context.startDate.Year())
		}
	}

	contextName = "c3_From1Mar2013To31Aug2013"
	context, exists = frp.contextMap[contextName]

	if !exists {
		t.Fatal("Missing context: ", contextName)
	} else {
		if context.startDate.Month() != 3 || context.startDate.Year() != 2013 {
			t.Fatal("Context <", contextName, "> has wrong startDate of ", context.startDate.Month(), " - ", context.startDate.Year())
		} else if context.endDate.Month() != 8 || context.endDate.Year() != 2013 {
			t.Fatal("Context <", contextName, "> has wrong endDate of ", context.startDate.Month(), " - ", context.startDate.Year())
		}
	}

	contextName = "c4_AsOf31Aug2014"
	context, exists = frp.contextMap[contextName]

	if !exists {
		t.Fatal("Missing context: ", contextName)
	} else {
		if context.instant.Month() != 8 || context.instant.Year() != 2014 {
			t.Fatal("Context <", contextName, "> has wrong instant of ", context.instant.Month(), " - ", context.instant.Year())
		}
	}
}

func TestParseInt64Field(t *testing.T) {
	xbrlFile, _ := os.Open("../testData/TestCreateParserMap.xml")
	fileReader := bufio.NewReader(xbrlFile)
	decoder := xml.NewDecoder(fileReader)

	parserMap := createParserMap(decoder)

	rawReport := &filings.FinancialReportRaw{}
	rawReport.RawFields = make(map[string]int64)

	frp := NewFinancialReportParser("../testData/TestCreateParserMap.xml", rawReport, nil, &filings.BasicRawFieldNameList{})

	frp.contextMap = make(map[string]*context)

	parseContext(frp, parserMap[contextTag])
	parseInt64Field(frp, parserMap["Revenues"])

	rawRev := frp.financialReportRaw.RawFields["Revenues"]
	if rawRev != 9457448 {
		t.Fatal("Expected revenue was 9457448, received: ", rawRev)
	}
}

func TestAlternateTags_Revenue(t *testing.T) {
	// test that we can handle multiple exclusive tags for the same value
	mockPersister := &MockPersister{}
	frp := NewFinancialReportParser("",
		&filings.FinancialReportRaw{},
		mockPersister, &filings.BasicRawFieldNameList{})

	rawFields := make(map[string]int64)
	rawFields["Revenues"] = 12345
	rawFields["SalesRevenueNet"] = 1234

	frp.financialReportRaw.RawFields = rawFields

	fr := frp.GetFinancialReport()

	if fr.Revenue != 12345 {
		t.Fatal("Expected Revenues was 12345, got: ", fr.Revenue)
	}

	rawFields = make(map[string]int64)
	rawFields["Revenues"] = 12345
	rawFields["SalesRevenueNet"] = 123456

	frp.financialReportRaw.RawFields = rawFields

	fr = frp.GetFinancialReport()

	if fr.Revenue != 123456 {
		t.Fatal("Expected Revenues was 123456, got: ", fr.Revenue)
	}
}

func TestContextParseXmlTagAtTop(t *testing.T) {
	// golang by default only decodes utf. test handling ascii
	mockPersister := &MockPersister{}
	frp := NewFinancialReportParser("../testData/TestContext2.xml",
		&filings.FinancialReportRaw{CIK: 95304, Year: 2014, Quarter: 2},
		mockPersister, &filings.BasicRawFieldNameList{})
	frp.financialReportRaw.RawFields = make(map[string]int64)

	frp.Parse()

	if len(frp.contextMap) == 0 {
		t.Fatal("Context map was 0 for TestContext2.xml")
	}
}

func TestCompleteParseRMCF_2014_2(t *testing.T) {
	mockPersister := &MockPersister{}
	frp := NewFinancialReportParser("../testData/rmcf-20140831.xml", &filings.FinancialReportRaw{CIK: 785815, Year: 2014, Quarter: 2},
		mockPersister, &filings.BasicRawFieldNameList{})

	frp.financialReportRaw.RawFields = make(map[string]int64)

	rawFields := make(map[string]int64)
	rawFields["NetCashProvidedByUsedInOperatingActivities"] = 82978
	rawFields["PaymentsToAcquirePropertyPlantAndEquipment"] = 98023

	mockPersister.SetFinancialReport(&filings.FinancialReportRaw{CIK: 785815, Year: 2014, Quarter: 1, RawFields: rawFields})

	frp.Parse()

	fieldsToValidate := frp.financialReportRaw.RawFields

	if fieldsToValidate["Revenues"] != 9457448 {
		t.Fatal("Expected Raw Revenues was 9457448, received: ", fieldsToValidate["Revenues"])
	}

	if fieldsToValidate["CostsAndExpenses"] != 8028307 {
		t.Fatal("Expected Raw CostsAndExpenses was 8028307, received: ", fieldsToValidate["CostsAndExpenses"])
	}

	if fieldsToValidate["NetIncomeLoss"] != 877356 {
		t.Fatal("Expected Raw NetIncomeLoss was 877356, received: ", fieldsToValidate["NetIncomeLoss"])
	}

	if fieldsToValidate["Assets"] != 38651192 {
		t.Fatal("Expected Raw Assets was 38651192, received: ", fieldsToValidate["Assets"])
	}

	if fieldsToValidate["NetCashProvidedByUsedInOperatingActivities"] != 3200000 {
		t.Fatal("Expected Raw NetCashProvidedByUsedInOperatingActivities was 3200000, received: ", fieldsToValidate["NetCashProvidedByUsedInOperatingActivities"])
	}

	fr := frp.GetFinancialReport()
	if fr.Revenue != 9457448 {
		t.Fatal("Expected Revenue was 9457448, received: ", fr.Revenue)
	}

	if fr.OperatingExpense != 8028307 {
		t.Fatal("Expected OperatingExpense was 8028307, received: ", fr.OperatingExpense)
	}

	if fr.TotalAssets != 38651192 {
		t.Fatal("Expected TotalAssets was 38651192, received: ", fr.TotalAssets)
	}

	if fr.OperatingCash != 3200000 {
		t.Fatal("Expected OperatingCash was 3200000, received: ", fr.OperatingCash)
	}

	if fr.CurrentAssets != 19435429 {
		t.Fatal("Expected CurrentAssets was 19435429, received: ", fr.CurrentAssets)
	}

	if fr.CurrentLiabilities != 8223763 {
		t.Fatal("Expected CurrentAssets was 8223763, received: ", fr.CurrentLiabilities)
	}

	if fr.TotalLiabilities != 14979344 {
		t.Fatal("Expected TotalLiabilities was 14979344, received: ", fr.TotalLiabilities)
	}

	if fr.CapitalExpenditures != 300000 {
		t.Fatal("Expected CapitalExpenditures was 300000, received: ", fr.CapitalExpenditures)
	}
}

type MockPersister struct {
	financialReportRaw *filings.FinancialReportRaw
}

func (mp *MockPersister) InsertUpdateRawReport(fr *filings.FinancialReportRaw) {}
func (mp *MockPersister) GetRawReport(cik, year, quarter int64) *filings.FinancialReportRaw {
	return mp.financialReportRaw
}

func (mp *MockPersister) SetFinancialReport(newFinancialReport *filings.FinancialReportRaw) {
	mp.financialReportRaw = newFinancialReport
}
