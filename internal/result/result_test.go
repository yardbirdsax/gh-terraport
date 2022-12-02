package result

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsMarkdownTable(t *testing.T) {
	results, err := FromSlice(
		[]interface{}{
			"id",
			"user",
			"created_date",
		},
		[][]interface{}{
			{
				1,
				"someone",
				time.Date(2022, 11, 07, 0, 0, 0, 0, time.UTC),
			},
			{
				2,
				"someone-else",
				time.Date(2022, 11, 8, 0, 0, 0, 0, time.UTC),
			},
		},
	)
	require.NoError(t, err)
	expectedOuput := `| id | user | created_date |
| ---:| --- | --- |
| 1 | someone | 2022-11-07 00:00:00 +0000 UTC |
| 2 | someone-else | 2022-11-08 00:00:00 +0000 UTC |`
	actualOutput := results.AsMarkdownTable()

	assert.Equal(t, expectedOuput, actualOutput)
}

func TestAsCSV(t *testing.T) {
	results, err := FromSlice(
		[]interface{}{
			"id",
			"user",
			"created_date",
		},
		[][]interface{}{
			{
				1,
				"someone",
				time.Date(2022, 11, 07, 0, 0, 0, 0, time.UTC),
			},
			{
				2,
				"someone-else",
				time.Date(2022, 11, 8, 0, 0, 0, 0, time.UTC),
			},
		},
	)
	require.NoError(t, err)
	expectedOutput := `id,user,created_date
1,someone,2022-11-07 00:00:00 +0000 UTC
2,someone-else,2022-11-08 00:00:00 +0000 UTC`
	actualOutput := results.AsCSV()

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name           string
		columnNames    []interface{}
		data           [][]interface{}
		expectedResult Results
		expectedError  string
	}{
		{
			"equal columns to data",
			[]interface{}{"date", "person", "score"},
			[][]interface{}{
				{time.Date(2022, 11, 1, 1, 23, 45, 678, time.UTC), "John", 1},
				{time.Date(2022, 11, 2, 1, 23, 45, 678, time.UTC), "Jane", 2},
			},
			Results{
				{"date", "person", "score"},
				{time.Date(2022, 11, 1, 1, 23, 45, 678, time.UTC), "John", 1},
				{time.Date(2022, 11, 2, 1, 23, 45, 678, time.UTC), "Jane", 2},
			},
			"",
		},
		{
			"fewer columns than data",
			[]interface{}{"date", "person"},
			[][]interface{}{
				{time.Date(2022, 11, 1, 1, 23, 45, 678, time.UTC), "John", 1},
				{time.Date(2022, 11, 2, 1, 23, 45, 678, time.UTC), "Jane", 2},
			},
			Results{},
			"one or more rows contain a different number of columns than what are named",
		},
		{
			"more columns than data",
			[]interface{}{"date", "person", "score", "extra"},
			[][]interface{}{
				{time.Date(2022, 11, 1, 1, 23, 45, 678, time.UTC), "John", 1},
				{time.Date(2022, 11, 2, 1, 23, 45, 678, time.UTC), "Jane", 2},
			},
			Results{},
			"one or more rows contain a different number of columns than what are named",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actualResult, err := FromSlice(tc.columnNames, tc.data)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, &tc.expectedResult, actualResult)
		})
	}
}
