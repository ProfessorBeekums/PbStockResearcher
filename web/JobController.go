package web

import (
	"net/http"
	"time"
	"github.com/ProfessorBeekums/PbStockResearcher/jobs"
	// "github.com/ProfessorBeekums/PbStockResearcher/log"
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
	// TODO params
	params := make(map[string]string)
	jm.StartJob(scraperJobFunc, jobs.JobTypeScraper, params)
}

func scraperJobFunc(params map[string]string) {
	// TODO call actual scraper
	time.Sleep(60 * time.Second)
}

// func parserJobFunc

func InitializeJobsEndpoints() {
	RegisterHttpHandler("jobs", HttpMethodGet, getJobsData)
	RegisterHttpHandler("jobs/scraper", HttpMethodPost, getJobsData)
}