package notes

import (
	"strconv"
	"time"
)

type Note struct {
	Company, Note string
	Timestamp int64
}

type NoteManager struct {
	Notes map[string]*Note
}

func (nm *NoteManager) GetNotes() map[string]*Note {
	return nm.Notes
}

func (nm *NoteManager) AddNote(company, note string) *Note {
	size := strconv.Itoa(len(nm.Notes))
	currentTime := time.Now().Unix()
	noteObj := &Note{Company: company, Note: note, Timestamp: currentTime}
	nm.Notes[size] = noteObj
	return noteObj
}
