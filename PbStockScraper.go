package main

import (
	"flag"
	"github.com/ProfessorBeekums/PbStockResearcher/config"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	"github.com/ProfessorBeekums/PbStockResearcher/scraper"
	"github.com/ProfessorBeekums/PbStockResearcher/tmpStore"
)

var year int
var quarter int

func init() {
	flag.IntVar(&year, "year", 0, "The year to scrape")
	flag.IntVar(&quarter, "quarter", 0, "The quarter to scrape")
}

func main() {
	log.Println("Starting program")

	flag.Parse()

	c := config.NewConfig("/home/beekums/Projects/stockResearch/config")

	log.Println("Loaded config: ", c)

	ts := tmpStore.NewTempStore(c.TmpDir)

	mysql := persist.NewMysqlDb(c.MysqlUser, c.MysqlPass, c.MysqlDb)

	pbScraper := scraper.NewEdgarFullIndexScraper(year, quarter, ts, mysql, mysql)

	pbScraper.ScrapeEdgarQuarterlyIndex()

	log.Println("Ending program")
}
