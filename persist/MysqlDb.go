package persist

import (
	"database/sql"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var driver = "mysql"

// TODO abandon this pattern and move functions into appropriate domains
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

func (mysql *MysqlPbStockResearcher) GetConnection() *sql.DB {
	return mysql.conn
}

////////////////////////////////BEGIN Persistence Functions////////////////////////////////////////////
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

	// we should only have one row if any
	row.Next()
	scanErr := row.Scan(&company.Name)

	if scanErr != nil {
		log.Error("Failed to scan row for cik <", cik, "> due to: ", scanErr)
	} else {
		company.CIK = cik
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
				, parsed=?
				, parse_error=?
			WHERE report_file_id=?`,
			reportFile.Year,
			reportFile.Quarter,
			reportFile.Filepath,
			reportFile.FormType,
			reportFile.Parsed,
			reportFile.ParseError,
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
	rows, err := mysql.conn.Query(`
		SELECT report_file_id as reportFileId
		, cik
		, year
		, quarter
		, filepath
		, form_type
		FROM report_file
		WHERE parsed = 0 and parse_error=0 
		LIMIT ?`, numToGet)

	reportFiles := make([]filings.ReportFile, numToGet)

	if err != nil {
		log.Error("Couldn't retrieve unparsed report files")
		return &[]filings.ReportFile{}
	}

	lastIndex := 0

	for rows.Next() {
		newReportFile := filings.ReportFile{}

		rows.Scan(&newReportFile.ReportFileId, &newReportFile.CIK, &newReportFile.Year,
			&newReportFile.Quarter, &newReportFile.Filepath, &newReportFile.FormType)

		reportFiles[lastIndex] = newReportFile
		lastIndex++
	}

	if int64(lastIndex) < (numToGet - 1) {
		// TODO yuch. Maybe slices aren't the way to go
		if lastIndex < 0 {
			lastIndex++
		}
		reportFiles = reportFiles[0:lastIndex]
	}

	return &reportFiles
}

func (mysql *MysqlPbStockResearcher) InsertUpdateFinancialReport(fr *filings.FinancialReport) {
	result, err := mysql.conn.Exec(
		`INSERT INTO financial_report (cik, year, quarter, report_file_id
			, revenue
			, operating_expense
			, net_income
			, current_assets
			, total_assets
			, current_liabilities
			, total_liabilities
			, operating_cash
			, capital_expenditures
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE report_file_id=VALUES(report_file_id)
			, revenue=VALUES(revenue)
			, operating_expense=VALUES(operating_expense)
			, net_income=VALUES(net_income)
			, current_assets=VALUES(current_assets)
			, total_assets=VALUES(total_assets)
			, current_liabilities=VALUES(current_liabilities)
			, total_liabilities=VALUES(total_liabilities)
			, operating_cash=VALUES(operating_cash)
			, capital_expenditures=VALUES(capital_expenditures)
		`,
		fr.CIK,
		fr.Year,
		fr.Quarter,
		fr.ReportFileId,
		fr.Revenue,
		fr.OperatingExpense,
		fr.NetIncome,
		fr.CurrentAssets,
		fr.TotalAssets,
		fr.CurrentLiabilities,
		fr.TotalLiabilities,
		fr.OperatingCash,
		fr.CapitalExpenditures,
	)

	if err != nil {
		log.Error("Failed to insert/update for ", fr.GetLogStr(),
			" because: ", err)
	} else {
		lastInsertId, insertErr := result.LastInsertId()
		if insertErr != nil {
			log.Error("Failed to get last insert id for ",
				fr.GetLogStr(), " because: ", insertErr)
		}
		fr.FinancialReportId = lastInsertId
	}
}

//func (mysql *MysqlPbStockResearcher) GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport {
//	// TODO unused for now
//	return nil
//}

func (mysql *MysqlPbStockResearcher) InsertUpdateRawReport(rawReport *filings.FinancialReportRaw) {
	numFields := len(rawReport.RawFields)

	if numFields < 1 {
		return
	}

	// create a new financial report to grab the primary key
	fr := &filings.FinancialReport{CIK: rawReport.CIK, Year: rawReport.Year, Quarter: rawReport.Quarter}
	mysql.InsertUpdateFinancialReport(fr)

	args := make([]interface{}, numFields*3)
	query :=
		`INSERT INTO financial_report_raw_fields (financial_report_id, field_name, field_value) VALUES `

	dbArgs := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		dbArgs[i] = "(?,?,?)"
	}

	var argIndex int = 0
	for fieldName, fieldValue := range rawReport.RawFields {
		args[argIndex] = fr.FinancialReportId
		argIndex++
		args[argIndex] = fieldName
		argIndex++
		args[argIndex] = fieldValue
		argIndex++
	}

	query += strings.Join(dbArgs, ",")

	query += " ON DUPLICATE KEY UPDATE field_value=VALUES(field_value)"

	_, err := mysql.conn.Exec(query, args...)

	if err != nil {
		log.Error("Failed to insert/update raw fields for ", fr.GetLogStr(),
			" because: ", err)
	}
}
func (mysql *MysqlPbStockResearcher) GetRawReport(cik, year, quarter int64) *filings.FinancialReportRaw {
	rawReport := new(filings.FinancialReportRaw)
	rawReport.CIK = cik
	rawReport.Year = year
	rawReport.Quarter = quarter

	rows, err := mysql.conn.Query(`
		SELECT financial_report_id
		, field_name
		, field_value
		FROM financial_report_raw_fields raw
		JOIN financial_report fr ON fr.financial_report_id = raw.financial_report_id
		WHERE fr.cik=? AND fr.year=? AND fr.quarter=?`,
		cik, year, quarter,
	)

	rawReport.RawFields = make(map[string]int64)

	if err != nil {
		log.Error("Couldn't retrieve unparsed report files")
		return nil
	}

	for rows.Next() {
		var fieldName string
		var fieldValue int64
		rows.Scan(&fieldName, &fieldValue)

		rawReport.RawFields[fieldName] = fieldValue
	}

	return rawReport
}

////////////////////////////////END Persistence Functions////////////////////////////////////////////

////////////////////////////////BEGIN Screener Functions////////////////////////////////////////////
// Using an implicit limit of 5000 in each query because if it's more than 5000, it's not a very good screen.
// Also... I don't want to deal with the performance of it being higher
const MAX_SCREEN_RESULTS = 5000

func (mysql *MysqlPbStockResearcher) GetRatio(ratioQuery string, year, quarter int, min, max float64) map[*filings.Company]float64 {
	query := "SELECT c.cik, c.name, " + ratioQuery + ` as ratio
				FROM financial_report fr
				JOIN company c on c.cik = fr.cik
				WHERE year = ? AND quarter = ? AND ` +
		ratioQuery + " > ? AND " +
		ratioQuery + ` < ?
				ORDER BY ratio DESC
				LIMIT ?`

	screenResults := make(map[*filings.Company]float64)
	rows, err := mysql.conn.Query(query, year, quarter, min, max, MAX_SCREEN_RESULTS)

	if err != nil {
		log.Error("Bad screen query for ScreenNetMargin: ", err)
	} else {
		for rows.Next() {
			company := &filings.Company{}
			var ratio float64
			rows.Scan(&company.CIK, &company.Name, &ratio)

			screenResults[company] = ratio
		}
	}

	return screenResults
}

func (mysql *MysqlPbStockResearcher) ScreenNetMargin(year, quarter int, min, max float64) map[*filings.Company]float64 {
	return mysql.GetRatio("net_income / revenue", year, quarter, min, max)
}

func (mysql *MysqlPbStockResearcher) ScreenAssetRatio(year, quarter int, min, max float64) map[*filings.Company]float64 {
	return mysql.GetRatio("total_assets / (total_liabilities + total_assets)", year, quarter, min, max)
}

func (mysql *MysqlPbStockResearcher) ScreenCurrentRatio(year, quarter int, min, max float64) map[*filings.Company]float64 {
	return mysql.GetRatio("current_assets / (total_liabilities + current_assets)", year, quarter, min, max)
}

////////////////////////////////END Screener Functions////////////////////////////////////////////
