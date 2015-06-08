package notes

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
)

type NoteFilter struct {
	CompanyName string
	NoteFilterId, CIK int64
}

type NoteFilterManager struct {
	noteFilters map[int64]*NoteFilter
	persister *persist.MysqlPbStockResearcher
}

func GetNewNoteFilterManager(persister *persist.MysqlPbStockResearcher) *NoteFilterManager {
	noteManager := &NoteFilterManager{}
	noteManager.noteFilters = make(map[int64]*NoteFilter)
	noteManager.persister = persister

	return noteManager
}

func (nm *NoteFilterManager) GetNoteFilters() map[int64]*NoteFilter {
	rows, err := nm.persister.GetConnection().Query(`
		SELECT nf.note_filter_id, nf.cik, c.name
		FROM note_filters nf
		JOIN company c on nf.cik = c.cik
		`)

	if err != nil {
		log.Error("Failed to load note filters due to: ", err)
		return nm.noteFilters
	}

	for rows.Next() {
		loadedNoteFilter := NoteFilter{}

		rows.Scan(&loadedNoteFilter.NoteFilterId, &loadedNoteFilter.CIK, &loadedNoteFilter.CompanyName)

		nm.noteFilters[loadedNoteFilter.NoteFilterId] = &loadedNoteFilter
	}

	return nm.noteFilters
}

func (nm *NoteFilterManager) AddNoteFilter(cik int64) *NoteFilter {
	company := nm.persister.GetCompany(cik);

	noteFilterObj := &NoteFilter{CIK: company.CIK}

	_, err := nm.persister.GetConnection().Exec(
		`INSERT INTO note_filters (
			cik
		) VALUES (?)
		`,
		noteFilterObj.CIK,
	)

	if err != nil {
		log.Error("Failed to insert note for ", company.Name,
			" because: ", err)
		noteFilterObj = nil
	} else {
		//		nm.notes[noteObj.NoteId] = noteObj
	}

	return noteFilterObj
}
