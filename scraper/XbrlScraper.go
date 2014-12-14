package scraper

import (
	"archive/zip"
	"bufio"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/tmpStore"
	"io"
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
	ts            *tmpStore.TempStore
}

func NewEdgarFullIndexScraper(year, quarter int,
	ts *tmpStore.TempStore) *EdgarFullIndexScraper {
	return &EdgarFullIndexScraper{year: year, quarter: quarter, ts: ts}
}

func (efis *EdgarFullIndexScraper) ScrapeEdgarQuarterlyIndex() {
	log.Println("Starting to scrape the full index for year <", efis.year,
		"> and quarter:", efis.quarter)

	indexUrl := EDGAR_FULL_INDEX_URL_PREFIX +
		strconv.FormatInt(int64(efis.year), 10) +
		"/QTR" + strconv.FormatInt(int64(efis.quarter), 10) + INDEX_FILE_NAME

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

				bucket := getBucket(cik)
				fileKey := getKey(formType, efis.year, efis.quarter)

				filePath := efis.ts.GetFilePath(bucket, fileKey)

				if filePath == "" {
					efis.GetXbrl(filename, bucket, fileKey)
					// TODO - temporary hack for testing
					break
				} else {
					log.Println("SKIP <", filename, "> because it already exists in: ", filePath)
				}
			}
		}
	}
}

// The full index provides links to txt files. We want to convert these to retrieve the corresponding zip of xbrl files
// and extract the main xbrl file.
func (efis *EdgarFullIndexScraper) GetXbrl(edgarFilename, bucket, fileKey string) {
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

		outputFileName := strconv.Itoa(int(time.Now().Unix())) + baseName + XBRL_ZIP_SUFFIX
		zipFilePath := efis.ts.StoreFile(bucket, outputFileName, getResp.Body)

		if zipFilePath != "" {
			efis.getXbrlFromZip(zipFilePath, bucket, fileKey)
		}
	}
}

func (efis *EdgarFullIndexScraper) getXbrlFromZip(zipFileName, bucket, fileKey string) {
	zipReader, zipErr := zip.OpenReader(zipFileName)

	if zipErr != nil {
		log.Error("Failed to open zip: ", zipFileName, " with error: ", zipErr)
	} else {
		defer zipReader.Close()

		foundOne := false

		for _, zippedFile := range zipReader.File {
			zippedFileName := zippedFile.Name
			isMatch := isXbrlFileMatch(zippedFileName)
			if isMatch {
				foundOne = true
				log.Println("Found zipped file: ", zippedFileName)

				xbrlFile, xbrlErr := zippedFile.Open()

				defer xbrlFile.Close()

				if xbrlErr != nil {
					log.Error("Failed to open zip file")
				} else {
					efis.ts.StoreFile(bucket, fileKey, xbrlFile)
				}

				// we don't care about the other stuff
				break
			}
		}

		if foundOne == false {
			log.Error("Could not find a match for an xbrl in ", zipFileName)
		}
	}
}

func getBucket(cik string) string {
	return "CIK_" + cik
}

func getKey(formType string, year, quarter int) string {
	return "Y" + strconv.Itoa(year) + "Q" + strconv.Itoa(quarter) + "FT" + formType
}

func isXbrlFileMatch(fileName string) bool {
	isMatch, _ := regexp.MatchString("([a-z]|[0-9])+-[0-9]+.xml", fileName)
	return isMatch
}
