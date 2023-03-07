// Package json2csv provides JSON to CSV functions.
package json2csv

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

type JSONStreamReader interface {
	Read() map[string]interface{}
	HasNext() bool
}

type CSVHeader map[string]interface{}

// JSON2CSV converts JSON to CSV.
// Update CSVHeader according to the data provided if csvHeader is not nil
func JSON2CSV(data interface{}, csvHeader CSVHeader) ([]KeyValue, error) {
	results := []KeyValue{}
	v := valueOf(data)
	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			result, err := flatten(v)
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
		if isObjectArray(v) {
			for i := 0; i < v.Len(); i++ {
				result, err := flatten(v.Index(i))
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
			result, err := flatten(v)
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

func JSON2CSVHeader(reader JSONStreamReader) (CSVHeader, error) {
	header := CSVHeader{}
	for reader.HasNext() {
		_, err := JSON2CSV(reader.Read(), header)
		if err != nil {
			return header, nil
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

func JSON2CSVOnline(reader JSONStreamReader, csvHeader CSVHeader, output io.Writer) error {
	writer := NewCSVWriter(output, DotBracketStyle, false)
	err := writer.WriterHeader(csvHeader)
	if err != nil {
		return err
	}
	for reader.HasNext() {
		csvRow, err := JSON2CSV(reader.Read(), nil)
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
