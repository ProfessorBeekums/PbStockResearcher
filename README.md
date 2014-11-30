PbStockResearcher
=================

A tool to make researching stocks faster. This is NOT meant to substitute human effort, just reduce it. 

I personally find existing tools to be very useful, but don't fit me very well anymore. I figure that existing tools
must get there data from somewhere. That somewhere is probably the SEC which makes all filings public.

The initial milestone (version 0.1) for this project will be:
- Parse out XBRL files from the SEC full indexes
- Parse out basic financial information for companies from the XBRL files (revenue, income, assets, etc)

Eventual goals for this project (incomplete list):
- Custom stock screeners
- Recommendation engine for stocks to investigate. This is NOT an algorithm for automatic trading. You don't watch everything Netflix recommends right? It is a nice place to start though to find something interesting.
- Custom dashboards for portfolios
- Alerts for actions such as SEC investigations
- Track previous companies key executives have run and evaluate those with the current company performance.

TODOs for 0.1 release:
- figure out a data store for files downloaded from SEC
- figure out a data store for parse data
- create a configuration pattern
- finish parsing out financial data (only revenue is covered)

Design Notes: 
- log package creates a wrapper around golang's log. This is so that errors and info messages can eventually be redirected to separate places
- scraper package will parse an index file from the SEC and start retrieving XBRL files.
- parser package parses XBRL files
- filings package holds the data that's parsed.

Usage:
Currently the project is a collection of libraries. There is nothing tying it all together yet.

To scrape a quarterly index:
```go
    scraper := scraper.NewEdgarFullIndexScraper(2013, 1)
    scraper.ScrapeEdgarQuarterlyIndex()
```
This will retrieve data for the 1st quarter of 2013. NOTE: data before fourth quarter of 2010 is potentially unreliable.

To parse data for an xbrl file retrieved from the scraper:
```go
    frp := parser.NewFinancialReportParser("testData.xml")
    frp.Parse()
    report := frp.GetFinancialReport()
```
