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

func (nm *NoteManager) AddNote(company, note string) {
	size := strconv.Itoa(len(nm.Notes))
	currentTime := time.Now().Unix()
	nm.Notes[size] = &Note{Company: company, Note: note, Timestamp: currentTime}
}