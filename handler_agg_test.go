package main

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	testCases := []struct {
		input    string
		expected time.Time
	}{
		{
			input:    "Fri, 30 Aug 2019 00:00:00 +0000",
			expected: time.Date(2019, 8, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		dt, err := parseTime(tc.input)
		if err != nil {
			t.Errorf("Failed to parse date: %v\n", err)
		}
		if dt != tc.expected {
			t.Errorf("Wrong date - got %v, wanted %v", dt, tc.expected)
		}
	}
}
