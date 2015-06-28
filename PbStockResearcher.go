package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/ProfessorBeekums/PbStockResearcher/config"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/notes"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
)

const HttpMethodGet = "GET"
const HttpMethodPost = "POST"
const HttpMethodPut = "PUT"

type methodMap map[string]func(w http.ResponseWriter, r *http.Request)

var routeMaps map[string]*methodMap

var templateLock = new(sync.RWMutex)
var templates *template.Template

var noteManager *notes.NoteManager
var noteFilterManager *notes.NoteFilterManager
var mysql *persist.MysqlPbStockResearcher

func loadTemplates() {
	parsedTemplates, parseErr := template.ParseFiles("ui/index.html", "ui/companyDash.html")

	if parseErr == nil {
		templateLock.Lock()
		templates = parsedTemplates
		templateLock.Unlock()
	} else {
		log.Error(parseErr)
	}
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	handleTemplate(w, r, "index.html");
}

func handleCompanyDash(w http.ResponseWriter, r *http.Request) {
	// TODO this will need to read the cik and start a filter
	handleTemplate(w, r, "companyDash.html")
}

func handleTemplate(w http.ResponseWriter, r *http.Request, templateFile string) {
	templateLock.RLock()
	defer templateLock.RUnlock()
	err := templates.ExecuteTemplate(w, templateFile, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	noteMap := noteManager.GetNotes()
	noteArray := make([]*notes.Note, len(noteMap))

	index := 0

	for _, note := range noteMap {
		noteArray[index] = note
		index++
	}

	returnJson(w, noteArray)
}

func postNotes(w http.ResponseWriter, r *http.Request) {
	cikVal := r.FormValue("cik")
	note := r.FormValue("note")

	cik, parseErr := strconv.Atoi(cikVal)

	if parseErr != nil {
		fmt.Fprintln(w, "ERROR parsing cik: ", parseErr)
		return
	}

	noteObj := noteManager.AddNote(int64(cik), note)

	jsonData, err := json.Marshal(noteObj)

	if err != nil {
		fmt.Fprintln(w, "ERROR encoding json data: ", err)
	} else {
		fmt.Fprintln(w, string(jsonData))
	}
}

func getNoteFilters(w http.ResponseWriter, r *http.Request) {
	noteMap := noteFilterManager.GetNoteFilters()
	noteArray := make([]*notes.NoteFilter, len(noteMap))

	index := 0

	for _, note := range noteMap {
		noteArray[index] = note
		index++
	}

	returnJson(w, noteArray)
}

func postNoteFilters(w http.ResponseWriter, r *http.Request) {
	cikVal := r.FormValue("cik")

	cik, parseErr := strconv.Atoi(cikVal)

	if parseErr != nil {
		fmt.Fprintln(w, "ERROR parsing cik: ", parseErr)
		return
	}

	noteFilterObj := noteFilterManager.AddNoteFilter(int64(cik))

	jsonData, err := json.Marshal(noteFilterObj)

	if err != nil {
		fmt.Fprintln(w, "ERROR encoding json data: ", err)
	} else {
		fmt.Fprintln(w, string(jsonData))
	}
}

func addCompany(w http.ResponseWriter, r *http.Request) {
	cikVal := r.FormValue("cik")
	companyVal := r.FormValue("company")

	cik, parseErr := strconv.Atoi(cikVal)

	if parseErr != nil {
		fmt.Fprintln(w, "ERROR parsing cik: ", parseErr)
		return
	}

	company := &filings.Company{CIK:int64(cik), Name:companyVal}
	mysql.InsertUpdateCompany(company)

	returnJson(w, company)
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

	returnJson(w, returnMap)
}

func main() {
	log.Println("Starting web server...")

	c := config.NewConfig("/home/beekums/Projects/stockResearch/config")

	log.Println("Loaded config: ", c)

	mysql = persist.NewMysqlDb(c.MysqlUser, c.MysqlPass, c.MysqlDb)

	loadTemplates()

	go func() {
		// TODO replace with fs-notify after updating go to 1.5
		// https://github.com/go-fsnotify/fsnotify/blob/master/example_test.go
	    for {
	    	time.Sleep(2 * time.Second)
		    loadTemplates()
		}
	}()

	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/companyDash", handleCompanyDash)

	registerHttpHandler("note", HttpMethodGet, getNotes)
	registerHttpHandler("note", HttpMethodPost, postNotes)

	registerHttpHandler("note-filter", HttpMethodGet, getNoteFilters)
	registerHttpHandler("note-filter", HttpMethodPost, postNoteFilters)

	registerHttpHandler("company", HttpMethodPost, addCompany)
	registerHttpHandler("company", HttpMethodGet, getCompanyData)

	noteManager = notes.GetNewNoteManager(mysql)
	noteFilterManager = notes.GetNewNoteFilterManager(mysql)

	// handle assets
	http.Handle("/ui/js/", http.StripPrefix("/ui/js/", http.FileServer(http.Dir("ui/js/"))))
	http.Handle("/ui/css/", http.StripPrefix("/ui/css/", http.FileServer(http.Dir("ui/css/"))))

	http.ListenAndServe(":4000", nil)
}

func returnJson(w http.ResponseWriter, response interface{}) {
	jsonData, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintln(w, "ERROR encoding json data: ", err)
	} else {
		fmt.Fprintln(w, string(jsonData))
	}
}

func registerHttpHandler(uri, httpMethod string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	if routeMaps == nil {
		routeMaps = make(map[string]*methodMap)
	}

	existingRoute, exists := routeMaps[uri]
	if exists {
		if (*existingRoute)[httpMethod] != nil {
			// bad! can't add another route when one already exists for that method
			fmt.Println("ERROR registerHttpHandler: Can't add route because route already exists for <", uri, "> with method <", httpMethod)
		} else {
			// add to the existing route
			(*routeMaps[uri])[httpMethod] = handlerFunc
		}
	} else {
		//make a new route
		newMethodMap := new(methodMap)
		(*newMethodMap) = make(map[string]func(w http.ResponseWriter, r *http.Request))
		(*newMethodMap)[httpMethod] = handlerFunc

		// we're assuming no leading or trailing slashes in the uri
		routeMaps[uri] = newMethodMap

		uri1 := "/" + uri
		uri2 := uri1 + "/"
		http.HandleFunc(uri1, requestHandler)
		http.HandleFunc(uri2, requestHandler)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	// remove leading and trailing slashes
	uri := strings.Trim(r.URL.Path, "/")

	// try to find a mapping directly. if that fails, remove the last mapping and try again
	var routeMap *methodMap = nil
	var exists bool = false

	// TODO this is kind of ghetto. find a better way to do these routes
	for !exists && len(uri) > 0 {
		routeMap, exists = routeMaps[uri]

		if !exists {
			lastSlashIndex := strings.LastIndex(uri, "/")

			if lastSlashIndex > 0 {
				uri = uri[0:lastSlashIndex]
			} else {
				uri = ""
			}
		}
	}

	if routeMap == nil {
		fmt.Fprintln(w, "ERROR requestHandler: no route for <", uri)
	} else {
		handlerFunc, methodExists := (*routeMap)[r.Method]

		if !methodExists {
			fmt.Fprintln(w, "ERROR requestHandler: method does not exist <", r.Method)
		} else {
			handlerFunc(w, r)
		}
	}
}

func webError(w http.ResponseWriter, err string) {
	log.Error(err)
	http.Error(w, err, 500)
}
