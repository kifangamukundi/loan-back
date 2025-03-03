package exporters

import (
	"fmt"
	"reflect"

	"github.com/xuri/excelize/v2"
)

func ExportToXLSX(data interface{}, filename string, selectedFields []string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.NewSheet(sheet)

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

	for i, field := range filteredFields {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, field)
	}

	for i := 0; i < sliceValue.Len(); i++ {
		row := sliceValue.Index(i)
		rowValues := getStructValues(row, filteredFields)

		for j, value := range rowValues {
			cell := fmt.Sprintf("%c%d", 'A'+j, i+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

func ExportSingleToXLSX(data interface{}, filename string, selectedFields []string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.NewSheet(sheet)

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

	for i, field := range filteredFields {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, field)
	}

	rowValues := getStructValues(dataValue, filteredFields)
	for i, value := range rowValues {
		cell := fmt.Sprintf("%c2", 'A'+i)
		f.SetCellValue(sheet, cell, value)
	}

	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}
