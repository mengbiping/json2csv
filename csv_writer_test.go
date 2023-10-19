package json2csv_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/yukithm/json2csv"
)

func TestKeyWithTrailingSpace(t *testing.T) {
	b := &bytes.Buffer{}
	wr := json2csv.NewCSVWriter(b, json2csv.JSONPointerStyle, false)
	responses := []map[string]interface{}{
		{
			" A":  1,
			"B ":  "foo",
			"C  ": "FOO",
		},
		{
			" A":  2,
			"B ":  "bar",
			"C  ": "BAR",
		},
	}

	csvContent, err := json2csv.JSON2CSV(responses, nil, math.MaxInt) // csvContent seems to be complete!
	if err != nil {
		t.Fatal(err)
	}
	wr.WriteCSV(csvContent)
	wr.Flush()

	got := b.String()
	want := `/ A,/B ,/C  
1,foo,FOO
2,bar,BAR
`

	if got != want {
		t.Errorf("Expected %v, but %v", want, got)
	}
}
