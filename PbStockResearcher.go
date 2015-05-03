package main

import (
	"html/template"
	"net/http"
	"sync"
	"time"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
)

var templateLock = new(sync.RWMutex)
var templates *template.Template

func loadTemplates() {
	parsedTemplates, parseErr := template.ParseFiles("ui/index.html")

	if parseErr == nil {
		templateLock.Lock()
		templates = parsedTemplates
		templateLock.Unlock()
	} else {
		log.Error(parseErr)
	}
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	templateLock.RLock()
	defer templateLock.RUnlock()
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	log.Println("Starting web server...")

	loadTemplates()

	go func() {
		// TODO replace with fs-notify after updating go to 1.5
		// https://github.com/go-fsnotify/fsnotify/blob/master/example_test.go
	    for {
	    	time.Sleep(2 * time.Second)
		    loadTemplates()
		    log.Println("Templates Reloaded")
		}
	}()

	http.HandleFunc("/", handleMainPage)

	http.ListenAndServe(":4000", nil)
}