package main

import (
	"github.com/ProfessorBeekums/PbStockResearcher/config"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/notes"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	"github.com/ProfessorBeekums/PbStockResearcher/web"
	"net/http"
)

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	web.HandleTemplate(w, r, "index.html")
}

func handleCompanyDash(w http.ResponseWriter, r *http.Request) {
	// TODO this will need to read the cik and start a filter
	web.HandleTemplate(w, r, "companyDash.html")
}

func handleJobDash(w http.ResponseWriter, r *http.Request) {
	web.HandleTemplate(w, r, "jobDash.html")
}

func main() {
	log.Println("Starting web server...")

	c := config.NewConfig("/home/beekums/Projects/stockResearch/config")

	log.Println("Loaded config: ", c)

	mysql := persist.NewMysqlDb(c.MysqlUser, c.MysqlPass, c.MysqlDb)

	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/companyDash", handleCompanyDash)
	http.HandleFunc("/jobDash", handleJobDash)

	noteManager := notes.GetNewNoteManager(mysql)
	noteFilterManager := notes.GetNewNoteFilterManager(mysql)

	web.InitializeJobsEndpoints()
	web.InitializeNotesEndpoints()
	web.InitializeCompanyEndpoints()

	web.StartWebServer(mysql, noteManager, noteFilterManager)
}
