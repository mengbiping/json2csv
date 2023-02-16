// Package json2csv provides JSON to CSV functions.
package json2csv

import (
	"errors"
	"io"
	"reflect"
)

type JSONStreamReader interface{}

type CSVHeader map[string]interface{}

// JSON2CSV converts JSON to CSV.
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
			for s := range result {
				csvHeader[s] = ""
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
				for s := range result {
					csvHeader[s] = ""
				}
			}
		} else if v.Len() > 0 {
			result, err := flatten(v)
			if err != nil {
				return nil, err
			}
			if result != nil {
				results = append(results, result)
				for s := range result {
					csvHeader[s] = ""
				}
			}
		}
	default:
		return nil, errors.New("Unsupported JSON structure.")
	}

	return results, nil
}

func JSON2CSVOnline(jsonStreamReader JSONStreamReader, csvHeader CSVHeader, output io.Writer) error {
	results, err := JSON2CSV(jsonStreamReader, CSVHeader{})
	if err != nil {
		return err
	}
	for _, res := range results {
		for h := range csvHeader {
			if _, exist := res[h]; !exist {
				res[h] = ""
			}
		}
	}
	csv := NewCSVWriter(output)
	if err := csv.WriteCSV(results, false, true); err != nil {
		return err
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
