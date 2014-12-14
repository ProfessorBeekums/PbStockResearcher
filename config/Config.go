package config

import (
	"encoding/json"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"os"
)

type Config struct {
	TmpDir, MongoHost, MongoDb string
}

func NewConfig(filePath string) *Config {
	file, fileErr := os.Open(filePath)

	log.Println("Loading configuration from: ", filePath)

	if fileErr != nil {
		log.Error("Couldn't open configuration file: ", filePath)
		// don't continue executing if we don't know our config
		panic(fileErr)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}
	decodeErr := decoder.Decode(config)

	if decodeErr != nil {
		log.Error("Couldn't decode configuration file: ", filePath)
		panic(decodeErr)
	}

	return config
}
