package json2csv

import (
	"archive/zip"
	"encoding/json"
	"io"
)

func NewJSONStreamZipReader(zipFileName string) JSONStreamReader {
	zipReader, _ := zip.OpenReader(zipFileName)
	return &JSONStreamZipReader{data: zipReader.File}
}

type JSONStreamZipReader struct {
	data  []*zip.File
	index int
}

func (jz *JSONStreamZipReader) HasNext() bool {
	return jz.index < len(jz.data)
}

func (jz *JSONStreamZipReader) Read() map[string]interface{} {
	child := jz.data[jz.index]
	jz.index++
	res := make(map[string]interface{})
	cfd, _ := child.Open()
	content, _ := io.ReadAll(cfd)
	_ = json.Unmarshal(content, &res)
	_ = cfd.Close()
	return res
}
