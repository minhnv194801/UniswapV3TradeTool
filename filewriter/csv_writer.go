package filewriter

import (
	"encoding/csv"
	"os"
)

type CSVWriter struct {
	filePath string
	fileName string
	fields   []string
}

func (C *CSVWriter) WriteData(data []string) error {
	// mk dir if not exist
	err := os.MkdirAll(C.filePath, os.ModePerm)
	if err != nil {
		return err
	}
	// open, if not existed, create
	file, err := os.OpenFile(C.filePath+C.fileName+".csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// check if file is empty
	writer := csv.NewWriter(file)
	defer writer.Flush()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		err = writer.Write(C.fields)
		if err != nil {
			return err
		}
	}
	err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func NewCSVWriter(filePath, fileName string, fields []string) WriterWriteLine {
	return &CSVWriter{
		fileName: fileName,
		filePath: filePath,
		fields:   fields,
	}
}
