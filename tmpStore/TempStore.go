package tmpStore

import (
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"io"
	"io/ioutil"
	"os"
	"syscall"
)

type TempStore struct {
	storeDir string
}

func NewTempStore(storeDir string) *TempStore {
	return &TempStore{storeDir: storeDir}
}

func (ts *TempStore) StoreFile(bucket, fileName string, fileReader io.Reader) {
	if getPercentDiskFree(ts.storeDir) < .25 {
		log.Error("Ran out of temp store disk")
		return;
	}

	data, readErr := ioutil.ReadAll(fileReader)

	if readErr != nil {
		log.Error("Failed to read bytes for saving to bucket <",
			bucket, "> and file name: ", fileName)
	}

	path := ts.storeDir + "/" + bucket
	_, pathErr := os.Stat(path)

	if pathErr != nil {
		// make the directory
		mkErr := os.Mkdir(path, 0777)
		if mkErr != nil {
			log.Error("Failed to create path: " + path)
			return;
		}
	}

	filePath := path + "/" + fileName
	writeErr := ioutil.WriteFile(filePath, data, 0777)

	if writeErr != nil {
		log.Error("Failed to write <",
            bucket, "> and file name: ", fileName)
	}
}

func getPercentDiskFree(dir string) float64 {
	var stat syscall.Statfs_t
    syscall.Statfs(dir, &stat)

    free := float64(stat.Bfree * uint64(stat.Bsize))
    total := float64(stat.Blocks * uint64(stat.Bsize))

	return free / total
}
