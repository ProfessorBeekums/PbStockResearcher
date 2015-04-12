package persist

import (
	"database/sql"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	_ "github.com/go-sql-driver/mysql"
)

var driver = "mysql"

// will implement all interfaces in the persist package
type MysqlPbStockResearcher struct {
	user, pass, table string
	conn              *sql.DB
}

func NewMysqlDb(user, pass, table string) *MysqlPbStockResearcher {
	mysqlDb := &MysqlPbStockResearcher{user: user, pass: pass, table: table}
	var err error

	mysqlDb.conn, err = sql.Open(driver,
		mysqlDb.user+":"+mysqlDb.pass+"@tcp/"+mysqlDb.table)

	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}

	return mysqlDb
}

func (mysql *MysqlPbStockResearcher) InsertUpdateCompany(company *filings.Company) {
	_, err := mysql.conn.Exec(
		`INSERT INTO company (cik, name) VALUES (?,?)
		ON DUPLICATE KEY UPDATE name=values(name)`,
		company.CIK,
		company.Name,
	)

	if err != nil {
		log.Error("Failed to upsert cik <",
			company.CIK, "> with name <", company.Name, "> because: ", err)
	}
}

func (mysql *MysqlPbStockResearcher) GetCompany(cik int64) *filings.Company {
	row, err :=
		mysql.conn.Query(`SELECT name FROM company WHERE cik = ?`, cik)

	company := new(filings.Company)

	if err != nil {
		log.Error("Failed to get company <", cik, "> because: ", err)
		return company
	}

	var name string

	// we should only have one row if any
	row.Next()
	scanErr := row.Scan(&name)

	if scanErr != nil {
		log.Error("Failed to scan row for cik <", cik, "> due to: ", scanErr)
	} else {
		company.CIK = cik
		company.Name = name
	}

	return company
}

func (mysql *MysqlPbStockResearcher) InsertUpdateReportFile(reportFile *filings.ReportFile) {

	if reportFile.ReportFileId == 0 {
		result, err := mysql.conn.Exec(
			`INSERT INTO report_file (cik, year, quarter, filepath, form_type) 
			VALUES (?,?,?,?,?)`,
			reportFile.CIK,
			reportFile.Year,
			reportFile.Quarter,
			reportFile.Filepath,
			reportFile.FormType,
		)

		if err != nil {
			log.Error("Failed to insert for ", reportFile.GetLogStr(),
				" because: ", err)
		} else {
			lastInsertId, insertErr := result.LastInsertId()
			if insertErr != nil {
				log.Error("Failed to get last insert id for ",
					reportFile.GetLogStr(), " because: ", insertErr)
			}
			reportFile.ReportFileId = lastInsertId
		}
	} else {
		result, err := mysql.conn.Exec(
			`UPDATE report_file SET
				year=?
				, quarter=?
				, filepath=?
				, form_type=?
			WHERE report_file_id=?`,
			reportFile.Year,
			reportFile.Quarter,
			reportFile.Filepath,
			reportFile.FormType,
			reportFile.ReportFileId,
		)

		if err != nil {
			log.Error("Failed to update for ", reportFile.GetLogStr(),
				" because: ", err)
		} else {
			rowsAffected, insertErr := result.RowsAffected()
			if insertErr != nil {
				log.Error("Failed to get rows affected for ", reportFile.GetLogStr(),
					" because: ", insertErr)
			}

			if rowsAffected < 1 {
				log.Error("Update successful, but rows affected less than 1 for ",
					reportFile.GetLogStr())
			}
		}
	}

}

func (mysql *MysqlPbStockResearcher) GetNextUnparsedFiles(numToGet int64) *[]filings.ReportFile {
	return nil
}
