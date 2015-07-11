package web

import (
	"github.com/ProfessorBeekums/PbStockResearcher/jobs"
	"net/http"
	// "github.com/ProfessorBeekums/PbStockResearcher/log"
	"strconv"
	"github.com/ProfessorBeekums/PbStockResearcher/scraper"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
)

func getJobsData(w http.ResponseWriter, r *http.Request) {
	jm := jobs.GetJobManager()
	allJobs := jm.GetJobs()

	jobArray := make([]*jobs.Job, len(allJobs))

	index := 0

	for _, job := range allJobs {
		jobArray[index] = job
		index++
	}
	ReturnJson(w, jobArray)
}

func startScraperJob(w http.ResponseWriter, r *http.Request) {
	jm := jobs.GetJobManager()

	params := make(map[string]string)
	params["year"] = r.FormValue("year")
	params["quarter"] = r.FormValue("quarter")

	jm.StartJob(scraperJobFunc, jobs.JobTypeScraper, params)
	ReturnJsonSuccess(w)
}

func scraperJobFunc(params map[string]string) {
	yearStr,_ := params["year"]
	quarterStr,_ := params["quarter"]
	year, yearParseErr := strconv.Atoi(yearStr)
	quarter, quarterParseErr := strconv.Atoi(quarterStr)

	if yearParseErr != nil || quarterParseErr != nil {
		log.Error("ERROR parsing year or quarter: ", yearParseErr, quarterParseErr)
		return
	}

	pbScraper := scraper.NewEdgarFullIndexScraper(year, quarter, ts, mysql, mysql)

	pbScraper.ScrapeEdgarQuarterlyIndex()
}

// func parserJobFunc

func InitializeJobsEndpoints() {
	RegisterHttpHandler("jobs", HttpMethodGet, getJobsData)
	RegisterHttpHandler("jobs/scraper", HttpMethodPost, startScraperJob)
}
