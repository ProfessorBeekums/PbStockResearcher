package web

import (
	"net/http"
	"github.com/ProfessorBeekums/PbStockResearcher/jobs"
	// "github.com/ProfessorBeekums/PbStockResearcher/log"
)

func getJobsData(w http.ResponseWriter, r *http.Request) {
	jm := jobs.GetJobManager()
	allJobs := jm.GetJobs()

	ReturnJson(w, allJobs)
}

// TODO start screener job function
// TODO start parser job function

func InitializeJobsEndpoints() {
	RegisterHttpHandler("jobs", HttpMethodGet, getJobsData)
}