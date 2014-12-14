package persist

import (
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const financialReportsCollection = "FinancialReport"

type MongoDbFinancialReports struct {
	host, database string
}

func NewMongoDbFinancialReports(host, database string) *MongoDbFinancialReports{
	return &MongoDbFinancialReports{host: host, database: database}
}

func (mdfr *MongoDbFinancialReports) getSessionAndCollection() (*mgo.Session, *mgo.Collection) {
	session, mongoErr := mgo.Dial(mdfr.host)
	if mongoErr != nil {
		log.Error("Failed to create session to MongoDb host: ", mdfr.host)
		return nil, nil
	} else {
		session.SetMode(mgo.Strong, true)

		db := session.DB(mdfr.database)
		collection := db.C(financialReportsCollection)	

		return session, collection
	}
}

func (mdfr *MongoDbFinancialReports) CreateFinancialReport(fr *filings.FinancialReport) {
	session, coll := mdfr.getSessionAndCollection()

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
	session, coll := mdfr.getSessionAndCollection()

    if session != nil {
        defer session.Close()

        err := coll.Update(bson.M{"cik" : fr.CIK, "year" : fr.Year, "quarter" : fr.Quarter} ,fr)

        if err != nil {
            log.Error("Failed to update financial report with cik <",
            fr.CIK, "> and year <", fr.Year, "> and quarter <", fr.Quarter,
            "> with error: ", err)
        }
    }
}

func (mdfr *MongoDbFinancialReports) GetFinancialReport(cik, year, quarter int64) *filings.FinancialReport {
	report := &filings.FinancialReport{}
	session, coll := mdfr.getSessionAndCollection()

    if session != nil {
        defer session.Close()
		
		coll.Find(bson.M{"cik" : cik, "year" : year, "quarter" : quarter}).One(&report)
	}

	return report
}
