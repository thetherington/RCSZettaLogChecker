package utils

import (
	"fmt"
	"strings"
	"time"
)

var origin = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)

// helper function to parse the zetta duration string format into a duration object
func ParseToDuration(input string) (time.Duration, error) {
	var layout string

	if strings.Count(input, ":") == 1 {
		layout = "04:05"
	} else {
		layout = "15:04:05"
	}

	t, err := time.Parse(layout, input)
	if err != nil {
		return 0, err
	}

	return t.Sub(origin), nil
}

// helper function to format duration into "HH:MM:SS"
func FmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// help function to create a string pointer from a function return
func Ptr(s string) *string { return &s }
