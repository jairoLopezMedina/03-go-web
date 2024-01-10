package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type JSON struct {
	file *os.File
}

func NewJSON(file *os.File) *JSON {
	return &JSON{
		file: file,
	}
}

func (j *JSON) Read() (p []Product, err error) {

	//create a new decoder
	dc := json.NewDecoder(j.file)

	if err := dc.Decode(&p); err == io.EOF {
		
	} else if err != nil {
		return p , fmt.Errorf("cannot decode: %w", err)
	}

	return
}