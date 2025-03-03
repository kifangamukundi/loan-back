package exporters

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

func ExportToCSV(data interface{}, filename string, selectedFields []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	sliceValue := reflect.ValueOf(data)
	if sliceValue.Kind() != reflect.Slice || sliceValue.Len() == 0 {
		return fmt.Errorf("data must be a non-empty slice of structs")
	}

	firstItem := sliceValue.Index(0)

	allFields := getStructFields(firstItem)

	filteredFields := []string{}
	for _, field := range selectedFields {
		if contains(allFields, field) {
			filteredFields = append(filteredFields, field)
		}
	}

	if len(filteredFields) == 0 {
		return fmt.Errorf("no valid fields selected for export")
	}

	if err := writer.Write(filteredFields); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	for i := 0; i < sliceValue.Len(); i++ {
		row := sliceValue.Index(i)
		rowValues := getStructValues(row, filteredFields)

		stringValues := make([]string, len(rowValues))
		for j, v := range rowValues {
			stringValues[j] = fmt.Sprintf("%v", v)
		}

		if err := writer.Write(stringValues); err != nil {
			return fmt.Errorf("failed to write row: %v", err)
		}
	}

	return nil
}

func ExportSingleToCSV(data interface{}, filename string, selectedFields []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct")
	}

	allFields := getStructFields(dataValue)

	filteredFields := []string{}
	for _, field := range selectedFields {
		if contains(allFields, field) {
			filteredFields = append(filteredFields, field)
		}
	}

	if len(filteredFields) == 0 {
		return fmt.Errorf("no valid fields selected for export")
	}

	if err := writer.Write(filteredFields); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	rowValues := getStructValues(dataValue, filteredFields)

	stringValues := make([]string, len(rowValues))
	for i, v := range rowValues {
		stringValues[i] = fmt.Sprintf("%v", v)
	}

	if err := writer.Write(stringValues); err != nil {
		return fmt.Errorf("failed to write row: %v", err)
	}

	return nil
}
