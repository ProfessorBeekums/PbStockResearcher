package scraper

import (
	"bufio"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// create a function to scrape an index given a year and quarter
// the function will also take in a delay between loading each file
// save results to a data store
// do not re-download things already in the data store!

const EDGAR_FULL_INDEX_URL_PREFIX = "http://www.sec.gov/Archives/edgar/full-index/"
const INDEX_FILE_NAME = "/xbrl.idx"

type EdgarFullIndexScraper struct {
	year, quarter int
}

func NewEdgarFullIndexScraper(year, quarter int) *EdgarFullIndexScraper {
	return &EdgarFullIndexScraper{year: year, quarter: quarter}
}

func (efis *EdgarFullIndexScraper) ScrapeEdgarQuarterlyIndex() {
	log.Println("Starting to scrape the full index for year <", efis.year,
		"> and quarter:", efis.quarter)

	indexUrl := EDGAR_FULL_INDEX_URL_PREFIX + 
		strconv.FormatInt(int64(efis.year), 10) +
		"/QTR" + strconv.FormatInt(int64(efis.quarter),10) + INDEX_FILE_NAME

	getResp, getErr := http.Get(indexUrl)

	if getErr != nil {
		log.Error("Failed to retrieve index for url <", indexUrl, 
		"> with error: ", getErr)
	} else if getResp.StatusCode != 200 {
		log.Error("Received status code <", getResp.Status, "> for url: ", indexUrl)
	} else {
		log.Println("@@@ Success!", indexUrl, getResp)
		defer getResp.Body.Close()

		efis.ParseIndexFile(getResp.Body)
	}
}

// Parses a ReadCloser that contains a Full Index file. The caller is
// responsible for closing the ReadCloser.
func (efis *EdgarFullIndexScraper) ParseIndexFile(fileReader io.ReadCloser) {
	listBegun := false // we need to parse the header before we get the list
	var line []byte = nil
	var readErr error = nil
	var isPrefix bool = false

	reader := bufio.NewReader(fileReader)
	for readErr == nil {
		// none of these lines should be bigger than the buffer 
		line, isPrefix, readErr = reader.ReadLine()
		if isPrefix {
			// don't bother parsing here, just log that we had an error
			log.Error("This index file has a line that's too long!")
			continue
		}

		if line != nil {
			lineStr := string(line)
			if !listBegun && strings.Contains(lineStr, "-------") {
				listBegun = true
				continue
			}

			// headers done, now we can start parsing
			if listBegun {
				elements := strings.Split(lineStr, "|")
				cik := elements[0]
				companyName := elements[1]
				formType := elements[2]
				dateFiled := elements[3]
				filename := elements[4]

				log.Println("CIK: ", cik, " Company Name: ", companyName, " Form type: ", formType, "  Date Filed: ", dateFiled, "  FileName: ", filename)
			}
		}
	}
}

/*

func getXbrl(edgarFilename string) {
    // TODO log if not a txt file
    parts := strings.Split(edgarFilename, "/")
    baseName := strings.Trim(parts[3], ".txt")
    preBase := strings.Replace(baseName, "-", "", -1)
    parts[3] = preBase + "/" + baseName + XBRL_ZIP_SUFFIX

    fullUrl := SEC_EDGAR_BASE_URL + strings.Join(parts, "/")

    logger.Println("getting xbrl from ", fullUrl)

    getResp, getErr := http.Get(fullUrl)

    if getErr != nil {
        logger.Println("Failed get to: ", fullUrl)
    } else {
        defer getResp.Body.Close()

        data, readErr := ioutil.ReadAll(getResp.Body)

        if readErr != nil {
            logger.Println("Failed to read")
        } else {
            outputFileName := time.Now().String() + baseName + XBRL_ZIP_SUFFIX
            ioutil.WriteFile(outputFileName, data, 0777)

            zipReader, zipErr := zip.OpenReader(outputFileName)
            if zipErr != nil {
                logger.Println("Failed to open zip: ", outputFileName)
            } else {
                defer zipReader.Close()

                for _, zippedFile := range zipReader.File {
                    zippedFileName := zippedFile.Name
                    isMatch,_ := regexp.MatchString("[a-z]+-[0-9]{8}.xml", zippedFileName)
                    if isMatch {
                        logger.Println("Found zipped file: ", zippedFileName ) 
                    }       
                }
            }
        }   
    }
}
*/
