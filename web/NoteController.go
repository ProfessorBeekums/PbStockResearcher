package web

import (
	"github.com/ProfessorBeekums/PbStockResearcher/notes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getNotes(w http.ResponseWriter, r *http.Request) {
	noteMap := noteManager.GetNotes()
	noteArray := make([]*notes.Note, len(noteMap))

	index := 0

	for _, note := range noteMap {
		noteArray[index] = note
		index++
	}

	ReturnJson(w, noteArray)
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

	ReturnJson(w, noteArray)
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

func InitializeNotesEndpoints() {
	RegisterHttpHandler("note", HttpMethodGet, getNotes)
	RegisterHttpHandler("note", HttpMethodPost, postNotes)

	RegisterHttpHandler("note-filter", HttpMethodGet, getNoteFilters)
	RegisterHttpHandler("note-filter", HttpMethodPost, postNoteFilters)
}
