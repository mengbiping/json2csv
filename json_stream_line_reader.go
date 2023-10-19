package json2csv

import (
	"encoding/json"
	"bufio"
	"os"
	"log"
)

func NewJSONStreamLineReader(f *os.File) JSONStreamReader {
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 1024 * 1024), 20* 1024 * 1024)
	jr := &JSONStreamLineReader{
		f: f,
		scanner: s,
	}
	jr.end = !s.Scan()
	return jr
}

type JSONStreamLineReader struct {
	f *os.File
	scanner *bufio.Scanner
	end bool
}

func (jr *JSONStreamLineReader) HasNext() bool {
	return !jr.end
}

func (jr *JSONStreamLineReader) Close() {
    jr.f.Close()
}

func (jr *JSONStreamLineReader) Read() map[string]interface{} {
	res := make(map[string]interface{})
	_ = json.Unmarshal(jr.scanner.Bytes(), &res)
	jr.end = !jr.scanner.Scan()
	if jr.end && jr.scanner.Err() != nil{
		log.Println(jr.scanner.Err())
	}
	return res
}
