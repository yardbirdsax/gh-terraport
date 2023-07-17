// Package result deals with generic ways to return and format the results of API calls.
package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Results [][]interface{}

func FromSlice(columnNames []interface{}, data [][]interface{}) (*Results, error) {
	result := Results{}

	// Write the header row
	columnCount := len(columnNames)
	result = append(result, columnNames)

	// Write the rest of the data
	for _, d := range data {
		if len(d) != columnCount {
			return &Results{}, fmt.Errorf("one or more rows contain a different number of columns than what are named")
		}
		result = append(result, d)
	}

	return &result, nil
}

// AsMarkdownTable renders tabular results as a Markdown table.
func (r *Results) AsMarkdownTable() string {
	var renderedString string

	t := populateTable(r)

	renderedString = t.RenderMarkdown()
	return renderedString
}

func (r *Results) AsCSV() string {
	var renderedString string

	t := populateTable(r)

	renderedString = t.RenderCSV()
	return renderedString
}

func (r *Results) AsJSON() string {
	renderedString := strings.Builder{}

	renderedString.WriteString("[")
	firstRow := (*r)[0]
	isFirstRow := true
	for _, res := range (*r)[1:] {
		if !isFirstRow {
			renderedString.WriteString(",")
		}
		renderedString.WriteString("{")
		isFirstElem := true
		for i := range res {
			if !isFirstElem {
				renderedString.WriteString(",")
			}
			renderedString.WriteString(fmt.Sprintf(`"%s":`, firstRow[i]))
			marshaledValue, _ := json.Marshal(res[i])
			renderedString.WriteString(string(marshaledValue))
			isFirstElem = false
		}
		renderedString.WriteString("}")
		isFirstRow = false
	}
	renderedString.WriteString("]")

	indented := bytes.Buffer{}
	json.Indent(&indented, []byte(renderedString.String()), "", "\t")

	return indented.String()
}

func populateTable(r *Results) table.Writer {
	t := table.NewWriter()

	// Construct the table header
	firstRow := (*r)[0]
	tableHeaderRow := table.Row{}
	for _, k := range firstRow {
		tableHeaderRow = append(tableHeaderRow, k)
	}
	t.AppendHeader(tableHeaderRow)

	for _, result := range (*r)[1:] {
		tableRow := table.Row{}
		for _, colVal := range result {
			typeOf := reflect.TypeOf(colVal)
			if typeOf.Name() == "string" {
				sanitizedValue := strings.ReplaceAll(colVal.(string), "https://api.github.com/", "")
				tableRow = append(tableRow, sanitizedValue)
			} else {
				tableRow = append(tableRow, colVal)
			}
		}
		t.AppendRow(tableRow)
	}
	return t
}
