package dateparse

import (
	"testing"
	"time"
)

func TestParseDate_Today(t *testing.T) {
	now := time.Now()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	result, err := ParseDate("today")
	if err != nil {
		t.Fatalf("ParseDate(\"today\") returned error: %v", err)
	}

	if !result.Equal(expected) {
		t.Errorf("ParseDate(\"today\") = %v, want %v", result, expected)
	}
}

func TestParseDate_Tomorrow(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	expected := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	result, err := ParseDate("tomorrow")
	if err != nil {
		t.Fatalf("ParseDate(\"tomorrow\") returned error: %v", err)
	}

	if !result.Equal(expected) {
		t.Errorf("ParseDate(\"tomorrow\") = %v, want %v", result, expected)
	}
}

func TestParseDate_EndOfWeek(t *testing.T) {
	result, err := ParseDate("end-of-week")
	if err != nil {
		t.Fatalf("ParseDate(\"end-of-week\") returned error: %v", err)
	}

	if result.Weekday() != time.Friday {
		t.Errorf("ParseDate(\"end-of-week\") returned %v (weekday: %v), want Friday", result, result.Weekday())
	}

	now := time.Now()
	if result.Before(now) {
		t.Errorf("ParseDate(\"end-of-week\") returned %v which is in the past", result)
	}
}

func TestParseDate_EndOfMonth(t *testing.T) {
	result, err := ParseDate("end-of-month")
	if err != nil {
		t.Fatalf("ParseDate(\"end-of-month\") returned error: %v", err)
	}

	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())

	if result.Month() != now.Month() && result.Before(firstOfNextMonth) {
		t.Errorf("ParseDate(\"end-of-month\") returned %v, expected last day of current month", result)
	}

	dayAfter := result.AddDate(0, 0, 1)
	if dayAfter.Month() == result.Month() {
		t.Errorf("ParseDate(\"end-of-month\") returned %v which is not the last day of the month", result)
	}
}

func TestParseDate_NextWeek(t *testing.T) {
	now := time.Now()
	nextWeek := now.AddDate(0, 0, 7)
	expected := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 0, 0, 0, 0, nextWeek.Location())

	result, err := ParseDate("next-week")
	if err != nil {
		t.Fatalf("ParseDate(\"next-week\") returned error: %v", err)
	}

	if !result.Equal(expected) {
		t.Errorf("ParseDate(\"next-week\") = %v, want %v", result, expected)
	}
}

func TestParseDate_NextMonth(t *testing.T) {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	expected := time.Date(nextMonth.Year(), nextMonth.Month(), nextMonth.Day(), 0, 0, 0, 0, nextMonth.Location())

	result, err := ParseDate("next-month")
	if err != nil {
		t.Fatalf("ParseDate(\"next-month\") returned error: %v", err)
	}

	if !result.Equal(expected) {
		t.Errorf("ParseDate(\"next-month\") = %v, want %v", result, expected)
	}
}

func TestParseDate_SpecificDate_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "standard date",
			input:    "2024-12-25",
			expected: time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "first day of year",
			input:    "2025-01-01",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "leap year date",
			input:    "2024-02-29",
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned error: %v", tt.input, err)
			}

			if !result.Equal(tt.expected) {
				t.Errorf("ParseDate(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDate_SpecificDate_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid format - wrong separator",
			input: "2024/12/25",
		},
		{
			name:  "invalid format - no separators",
			input: "20241225",
		},
		{
			name:  "invalid month",
			input: "2024-13-01",
		},
		{
			name:  "invalid day",
			input: "2024-02-30",
		},
		{
			name:  "random string",
			input: "not-a-date",
		},
		{
			name:  "empty string",
			input: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDate(tt.input)
			if err == nil {
				t.Errorf("ParseDate(%q) expected error, got nil", tt.input)
			}
		})
	}
}

func TestParseDate_AllKeywords(t *testing.T) {
	keywords := []string{
		"today",
		"tomorrow",
		"end-of-week",
		"end-of-month",
		"next-week",
		"next-month",
	}

	for _, keyword := range keywords {
		t.Run(keyword, func(t *testing.T) {
			result, err := ParseDate(keyword)
			if err != nil {
				t.Errorf("ParseDate(%q) returned error: %v", keyword, err)
			}

			if result.IsZero() {
				t.Errorf("ParseDate(%q) returned zero time", keyword)
			}

			if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
				t.Errorf("ParseDate(%q) returned %v with non-zero time component", keyword, result)
			}
		})
	}
}
