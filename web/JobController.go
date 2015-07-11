package web

import (
	"github.com/ProfessorBeekums/PbStockResearcher/jobs"
	"net/http"
	"time"
	// "github.com/ProfessorBeekums/PbStockResearcher/log"
//	"strconv"
//	"fmt"
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
	// TODO call actual scraper

//	year, yearParseErr := strconv.Itoa(yearStr)
//	quarter, quarterParseErr := strconv.Itoa(quarterStr)
//
//	if yearParseErr || quarterParseErr {
//		fmt.Fprintln(w, "ERROR parsing year or quarter: ", yearParseErr, quarterParseErr)
//		return
//	}
	time.Sleep(60 * time.Second)
}

// func parserJobFunc

func InitializeJobsEndpoints() {
	RegisterHttpHandler("jobs", HttpMethodGet, getJobsData)
	RegisterHttpHandler("jobs/scraper", HttpMethodPost, startScraperJob)
}
