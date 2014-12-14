package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const financialReportsCollection = "FinancialReport"
const companyCollection = "Company"

type MongoDbCompany struct {
	host, database string
}

func NewMongoDbCompany(host, database string) *MongoDbCompany {
	return &MongoDbCompany{host:host, database:database}
}

type MongoDbFinancialReports struct {
	host, database string
}

func NewMongoDbFinancialReports(host, database string) *MongoDbFinancialReports {
	return &MongoDbFinancialReports{host: host, database: database}
}

func getSessionAndCollection(host, database, collStr string) (*mgo.Session, *mgo.Collection) {
	session, mongoErr := mgo.Dial(host)
	if mongoErr != nil {
		log.Error("Failed to create session to MongoDb host: ", host)
		return nil, nil
	} else {
		session.SetMode(mgo.Strong, true)

		db := session.DB(database)
		collection := db.C(collStr)

		return session, collection
	}
}

func (mdc *MongoDbCompany) InsertUpdateCompany(company *filings.Company) {
	session, coll :=
		getSessionAndCollection(mdc.host, mdc.database, companyCollection)

	if session != nil {
		defer session.Close()
		
		_, upErr := coll.Upsert(bson.M{"cik":company.CIK}, company)
		if upErr != nil {
			log.Error("Failed to upsert company with cik <", company.CIK, 
			"> due to error: ", upErr)
		}
	}
}

func (mdc *MongoDbCompany) GetCompany(cik int64) *filings.Company {
	session, coll :=
        getSessionAndCollection(mdc.host, mdc.database, companyCollection)
	company := &filings.Company{}

    if session != nil {
        defer session.Close()
		
		coll.Find(bson.M{"cik": cik}).One(&company)
	}

	return company
}

func (mdfr *MongoDbFinancialReports) CreateFinancialReport(fr *filings.FinancialReport) {
	session, coll := 
		getSessionAndCollection(mdfr.host, mdfr.database, financialReportsCollection)

	if session != nil {
		defer session.Close()

		err := coll.Insert(fr)

		if err != nil {
			log.Error("Failed to insert financial report with cik <",
				fr.CIK, "> and year <", fr.Year, "> and quarter <", fr.Quarter,
				"> with error: ", err)
		}
	}
}

func (mdfr *MongoDbFinancialReports) UpdateFinancialReport(fr *filings.FinancialReport) {
	session, coll := 
	        getSessionAndCollection(mdfr.host, mdfr.database, financialReportsCollection)

	if session != nil {
		defer session.Close()

		err := coll.Update(bson.M{"cik": fr.CIK, "year": fr.Year, "quarter": fr.Quarter}, fr)

		if err != nil {
			log.Error("Failed to update financial report with cik <",
				fr.CIK, "> and year <", fr.Year, "> and quarter <", fr.Quarter,
				"> with error: ", err)
		}
	}
}

func (mdfr *MongoDbFinancialReports) GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport {
	report := &filings.FinancialReport{}
	session, coll := 
	        getSessionAndCollection(mdfr.host, mdfr.database, financialReportsCollection)

	if session != nil {
		defer session.Close()

		coll.Find(bson.M{"cik": cik, "year": year, "quarter": quarter}).One(&report)
	}

	return report
}
