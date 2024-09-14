package util

import (
	"encoding/json"
	"log"
	"time"
)

type Date struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Date type
func (d *Date) UnmarshalJSON(b []byte) error {
	dateStr := string(b)

	// Define multiple date formats to support
	formats := []string{
		`"2006-01-02T15:04:05Z"`,
		`"2006-01-02T15:04:05.000Z"`,
		`"2006-01-02T15:04:05"`,
		`"2006-01-02T15:04:05.000"`,
		`"2006-01-02"`,
		`"2006-01-02Z"`,
		`"2006-01-02T15:04:05-07:00"`,
		`"2006-01-02T15:04:05.000-07:00"`,
	}

	var parsedDate time.Time
	var err error

	// if blank, return empty time
	if dateStr == "" || dateStr == "\"\"" {
		d.Time = time.Now()
		return nil
	}

	// Try each format until one succeeds
	for _, format := range formats {
		parsedDate, err = time.Parse(format, dateStr)
		if err == nil {
			d.Time = parsedDate
			return nil
		}
	}

	// If no format worked, return an error
	log.Println("Error parsing date:", err)
	return err
}

// MarshalJSON to ensure consistent date formatting
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(time.RFC3339))
}
