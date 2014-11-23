package scraper

import (
	"archive/zip"
	"bufio"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// create a function to scrape an index given a year and quarter
// the function will also take in a delay between loading each file
// save results to a data store
// do not re-download things already in the data store!

const EDGAR_FULL_INDEX_URL_PREFIX = "http://www.sec.gov/Archives/edgar/full-index/"
const INDEX_FILE_NAME = "/xbrl.idx"

const SEC_EDGAR_BASE_URL = "http://www.sec.gov/Archives/"
const XBRL_ZIP_SUFFIX = "-xbrl.zip"

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

				efis.GetXbrl(filename)
				// TODO - temporary hack for testing
				break
			}
		}
	}
}

// The full index provides links to txt files. We want to convert these to retrieve the corresponding zip of xbrl files 
// and extract the main xbrl file.
func (efis *EdgarFullIndexScraper) GetXbrl(edgarFilename string) {
	if !strings.Contains(edgarFilename, ".txt") {
		log.Error("Unexpected file type: ", edgarFilename)
		return
	}

    parts := strings.Split(edgarFilename, "/")
    baseName := strings.Trim(parts[3], ".txt")
    preBase := strings.Replace(baseName, "-", "", -1)
    parts[3] = preBase + "/" + baseName + XBRL_ZIP_SUFFIX

    fullUrl := SEC_EDGAR_BASE_URL + strings.Join(parts, "/")

    log.Println("Getting xbrl zip from ", fullUrl)

    getResp, getErr := http.Get(fullUrl)

    if getErr != nil {
        log.Error("Failed get to: ", fullUrl)
    } else {
        defer getResp.Body.Close()

        data, readErr := ioutil.ReadAll(getResp.Body)

        if readErr != nil {
            log.Error("Failed to read")
        } else {
            outputFileName := strconv.Itoa(int(time.Now().Unix()) )+ baseName + XBRL_ZIP_SUFFIX
            // TODO configure a data directory
            writeErr := ioutil.WriteFile(outputFileName, data, 0777)

            
            if writeErr != nil {
            	log.Error("Failed to write file: ", writeErr)
            } else {
            	efis.getXbrlFromZip(outputFileName)
            }
        }   
    }
}

func (efis *EdgarFullIndexScraper) getXbrlFromZip(zipFileName string) {
	zipReader, zipErr := zip.OpenReader(zipFileName)

	if zipErr != nil {
        log.Error("Failed to open zip: ", zipFileName, " with error: ", zipErr)
    } else {
        defer zipReader.Close()

        for _, zippedFile := range zipReader.File {
            zippedFileName := zippedFile.Name
            isMatch,_ := regexp.MatchString("[a-z]+-[0-9]{8}.xml", zippedFileName)
            if isMatch {
                log.Println("Found zipped file: ", zippedFileName ) 

                xbrlFile, xbrlErr := zippedFile.Open()

                if xbrlErr != nil {
                	log.Error("Failed to open zip file")
                } else {
                	data, readErr := ioutil.ReadAll(xbrlFile)
                	if readErr != nil {
			            log.Error("Failed to read")
			        } else {
			        	writeErr := ioutil.WriteFile(zippedFileName, data, 0777)

			        	if writeErr != nil {
			            	log.Error("Failed to write file: ", writeErr)
			            }
			        }
                }

                // we don't care about the other stuff
                break
            }       
        }
    }
}
