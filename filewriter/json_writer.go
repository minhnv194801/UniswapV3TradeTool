package filewriter

import (
	"io/ioutil"
	"os"
)

type JSONWriter struct {
	filePath string
	fileName string
}

func (C *JSONWriter) GetData() []byte {
	file, err := os.OpenFile(C.filePath+C.fileName+".json", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	b, _ := ioutil.ReadAll(file)
	return b
}

func (C *JSONWriter) WriteData(data []byte) error {
	// mk dir if not exist
	err := os.MkdirAll(C.filePath, os.ModePerm)
	if err != nil {
		return err
	}
	// open, if not existed, create
	file, err := os.OpenFile(C.filePath+C.fileName+".json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// check if file is empty
	_, err = file.Write(data)

	return err
}

func NewJSONWriter(filePath, fileName string) WriterWriteObject {
	return &JSONWriter{
		fileName: fileName,
		filePath: filePath,
	}
}
