package web

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"fmt"
	"net/http"
	"strconv"
)

func addCompany(w http.ResponseWriter, r *http.Request) {
	cikVal := r.FormValue("cik")
	companyVal := r.FormValue("company")

	cik, parseErr := strconv.Atoi(cikVal)

	if parseErr != nil {
		fmt.Fprintln(w, "ERROR parsing cik: ", parseErr)
		return
	}

	company := &filings.Company{CIK: int64(cik), Name: companyVal}
	mysql.InsertUpdateCompany(company)

	ReturnJson(w, company)
}

func getCompanyData(w http.ResponseWriter, r *http.Request) {
	cikVal := r.FormValue("cik")

	cik, parseErr := strconv.Atoi(cikVal)

	if parseErr != nil {
		fmt.Fprintln(w, "ERROR parsing cik: ", parseErr)
		return
	}

	// TODO this seems silly, but only for now!
	company := mysql.GetCompany(int64(cik))
	returnMap := make(map[string]string)

	returnMap["name"] = company.Name

	ReturnJson(w, returnMap)
}

func InitializeCompanyEndpoints() {
	RegisterHttpHandler("company", HttpMethodPost, addCompany)
	RegisterHttpHandler("company", HttpMethodGet, getCompanyData)
}
