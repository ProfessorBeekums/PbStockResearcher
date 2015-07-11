package web

import (
	"encoding/json"
	"fmt"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/notes"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"
	"github.com/ProfessorBeekums/PbStockResearcher/tmpStore"
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
var ts *tmpStore.TempStore

func loadTemplates() {
	parsedTemplates, parseErr := template.ParseFiles("ui/index.html", "ui/companyDash.html", "ui/jobDash.html")

	if parseErr == nil {
		templateLock.Lock()
		templates = parsedTemplates
		templateLock.Unlock()
	} else {
		log.Error(parseErr)
	}
}

func HandleTemplate(w http.ResponseWriter, r *http.Request, templateFile string) {
	templateLock.RLock()
	defer templateLock.RUnlock()
	err := templates.ExecuteTemplate(w, templateFile, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartWebServer(newMysql *persist.MysqlPbStockResearcher,
	newNoteManager *notes.NoteManager,
	newNoteFilterManager *notes.NoteFilterManager,
	newTs *tmpStore.TempStore) {
	log.Println("Starting web server...")
	loadTemplates()

	mysql = newMysql
	noteManager = newNoteManager
	noteFilterManager = newNoteFilterManager
	ts = newTs

	go func() {
		// TODO replace with fs-notify after updating go to 1.5
		// https://github.com/go-fsnotify/fsnotify/blob/master/example_test.go
		for {
			time.Sleep(2 * time.Second)
			loadTemplates()
		}
	}()

	// handle assets
	http.Handle("/ui/js/", http.StripPrefix("/ui/js/", http.FileServer(http.Dir("ui/js/"))))
	http.Handle("/ui/css/", http.StripPrefix("/ui/css/", http.FileServer(http.Dir("ui/css/"))))

	// TODO put in config
	http.ListenAndServe(":4000", nil)
}

func ReturnJsonSuccess(w http.ResponseWriter) {
	jsonData, err := json.Marshal("success")

	if err != nil {
		fmt.Fprintln(w, "ERROR encoding json data: ", err)
	} else {
		fmt.Fprintln(w, string(jsonData))
	}
}

func ReturnJson(w http.ResponseWriter, response interface{}) {
	jsonData, err := json.Marshal(response)

	if err != nil {
		fmt.Fprintln(w, "ERROR encoding json data: ", err)
	} else {
		fmt.Fprintln(w, string(jsonData))
	}
}

func RegisterHttpHandler(uri, httpMethod string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
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

func WebError(w http.ResponseWriter, err string) {
	log.Error(err)
	http.Error(w, err, 500)
}
