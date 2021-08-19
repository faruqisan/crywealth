package storage

import (
	"encoding/csv"
	"os"
	"path"
)

const (
	dirLoc   = "crywealth"
	fileName = "wealth.csv"
)

type (
	Engine struct{}
	Wealth *os.File
)

func New() *Engine {
	return &Engine{}
}

// Load will load the file and create if not exists
func (*Engine) CreateIfNotExists() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := path.Join(homeDir, dirLoc)
	filePath := path.Join(dirPath, fileName)

	// check if exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// if not exist create
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			return err
		}

		_, err = os.Create(filePath)
		return err
	}
	return nil
}

func (*Engine) Read() (records [][]string, err error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	filePath := path.Join(homeDir, dirLoc, fileName)

	csvFile, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)

	records, err = r.ReadAll()
	if err != nil {
		return
	}

	return
}

func (*Engine) Write(record []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := path.Join(homeDir, dirLoc, fileName)

	csvFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	err = writer.Write(record)
	return err
}
