package notes

import (
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Note struct {
	CompanyName, Note      string
	NoteId, CIK, Timestamp int64
}

type NoteManager struct {
	notes     map[int64]*Note
	persister *persist.MysqlPbStockResearcher
}

func GetNewNoteManager(persister *persist.MysqlPbStockResearcher) *NoteManager {
	noteManager := &NoteManager{}
	noteManager.notes = make(map[int64]*Note)
	noteManager.persister = persister

	return noteManager
}

func (nm *NoteManager) GetNotes() map[int64]*Note {
	rows, err := nm.persister.GetConnection().Query(`
		SELECT n.note_id, c.name, n.cik, n.note_text, n.timestamp
		FROM notes n
		JOIN company c on c.cik = n.cik
		ORDER BY n.timestamp desc
		LIMIT 100`) // TODO make limit configurable

	if err != nil {
		log.Error("Failed to load notes due to: ", err)
		return nm.notes
	}

	for rows.Next() {
		loadedNote := Note{}

		rows.Scan(&loadedNote.NoteId, &loadedNote.CompanyName, &loadedNote.CIK, &loadedNote.Note, &loadedNote.Timestamp)

		nm.notes[loadedNote.NoteId] = &loadedNote
	}

	return nm.notes
}

func (nm *NoteManager) AddNote(cik int64, note string) *Note {
	company := nm.persister.GetCompany(cik)

	currentTime := time.Now().Unix()
	noteObj := &Note{CIK: company.CIK, CompanyName: company.Name, Note: note, Timestamp: currentTime}

	_, err := nm.persister.GetConnection().Exec(
		`INSERT INTO notes (
			cik
			, note_text
			, timestamp
		) VALUES (?,?,?)
		`,
		noteObj.CIK,
		noteObj.Note,
		noteObj.Timestamp,
	)

	if err != nil {
		log.Error("Failed to insert note for ", noteObj.CompanyName,
			" because: ", err)
		noteObj = nil
	} else {
		//		nm.notes[noteObj.NoteId] = noteObj
	}

	return noteObj
}
