// Package json2csv provides JSON to CSV functions.
package json2csv

import (
	"errors"
	"fmt"
	"github.com/yukithm/json2csv/jsonpointer"
	"io"
	"math"
	"reflect"
)

type JSONStreamReader interface {
	Read() map[string]interface{}
	HasNext() bool
	Close()
}

type CSVHeader map[string]interface{}

// JSON2CSV converts JSON to CSV.
// Update CSVHeader according to the data provided if csvHeader is not nil
func JSON2CSV(data interface{}, csvHeader CSVHeader, sliceLen int) ([]KeyValue, error) {
	results := []KeyValue{}
	v := valueOf(data)
	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			result, err := flatten(v, sliceLen)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
			if csvHeader != nil {
				for s := range result {
					csvHeader[s] = ""
				}
			}
		}
	case reflect.Slice:
		count := int(math.Min(float64(sliceLen), float64(v.Len())))
		if isObjectArray(v) {
			for i := 0; i < count; i++ {
				result, err := flatten(v.Index(i), sliceLen)
				if err != nil {
					return nil, err
				}
				results = append(results, result)
				if csvHeader != nil {
					for s := range result {
						csvHeader[s] = ""
					}
				}
			}
		} else if v.Len() > 0 {
			result, err := flatten(v, sliceLen)
			if err != nil {
				return nil, err
			}
			if result != nil {
				results = append(results, result)
				if csvHeader != nil {
					for s := range result {
						csvHeader[s] = ""
					}
				}
			}
		}
	default:
		return nil, errors.New("Unsupported JSON structure.")
	}

	return results, nil
}

func JSON2CSVHeader(reader JSONStreamReader, path string, sliceLen int) (CSVHeader, error) {
	header := CSVHeader{}
	var data interface{}
	var err error
	for reader.HasNext() {
		data = reader.Read()
		if path != "" {
			data, err = jsonpointer.Get(data, path)
			if err != nil {
				return header, err
			}
		}
		_, err := JSON2CSV(data, header, sliceLen)
		if err != nil {
			return header, err
		}
	}
	return header, nil
}

// FormatCSVHeaderToDotBracket convert given JSONPointerStyle header to DotBracketStyle.
func FormatCSVHeaderToDotBracket(header string) (string, error) {
	writer := NewCSVWriter(nil, DotBracketStyle, false)
	csvHeader := map[string]interface{}{header: nil}
	result, err := writer.FormatHeader(csvHeader)
	if err != nil {
		return "", fmt.Errorf("Failed to format %+v to DotBracket style: %w ",
			header, err)
	}
	return result[0], nil
}

func JSON2CSVOnline(reader JSONStreamReader, csvHeader CSVHeader, output io.Writer, style KeyStyle, transpose bool, path string, sliceLen int) error {
	writer := NewCSVWriter(output, style, transpose)
	err := writer.WriterHeader(csvHeader)
	if err != nil {
		return err
	}
	var data interface{}
	for reader.HasNext() {
		data = reader.Read()
		if path != "" {
			data, err = jsonpointer.Get(data, path)
			if err != nil {
				return err
			}
		}
		csvRow, err := JSON2CSV(data, nil, sliceLen)
		if err != nil {
			return err
		}
		err = writer.WriteCSVByHeader(
			csvRow,
			csvHeader,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func isObjectArray(obj interface{}) bool {
	value := valueOf(obj)
	if value.Kind() != reflect.Slice {
		return false
	}

	len := value.Len()
	if len == 0 {
		return false
	}
	for i := 0; i < len; i++ {
		if valueOf(value.Index(i)).Kind() != reflect.Map {
			return false
		}
	}

	return true
}
