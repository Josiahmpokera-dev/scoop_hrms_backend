package services

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// GetBiometricLocation returns the timezone used for parsing naive BioTime timestamps.
// Priority: BIOMETRIC_TIMEZONE, APP_TIMEZONE, Africa/Dar_es_Salaam, UTC.
func GetBiometricLocation() *time.Location {
	return getBiometricLocation()
}

// getBiometricLocation is the internal alias used by parsers in this package.
func getBiometricLocation() *time.Location {
	if tz := os.Getenv("BIOMETRIC_TIMEZONE"); tz != "" {
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}
	if tz := os.Getenv("APP_TIMEZONE"); tz != "" {
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}
	if loc, err := time.LoadLocation("Africa/Dar_es_Salaam"); err == nil {
		return loc
	}
	return time.UTC
}

// normalizeQueryDateTimeInput treats "YYYY-MM-DD+HH:MM:SS" like "YYYY-MM-DD HH:MM:SS".
// Browsers often send "+" as %2B between date and time; our layouts expect a space.
func normalizeQueryDateTimeInput(value string) string {
	s := strings.TrimSpace(value)
	if len(s) <= 11 || s[10] != '+' {
		return s
	}
	if _, err := time.Parse("2006-01-02", s[:10]); err != nil {
		return s
	}
	rest := strings.TrimSpace(s[11:])
	s = s[:10] + " " + rest
	return strings.Join(strings.Fields(s), " ")
}

// parseBiometricDateTime parses timestamp strings from BioTime/query params.
// Naive layouts are interpreted in business timezone; RFC3339 keeps embedded zone.
func parseBiometricDateTime(value string) (time.Time, error) {
	value = normalizeQueryDateTimeInput(value)
	loc := getBiometricLocation()

	naiveLayouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000000",
	}
	for _, layout := range naiveLayouts {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return t, nil
		}
	}

	zonedLayouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
	}
	for _, layout := range zonedLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", value)
}

// ParseBiometricDateTime parses timestamp strings from BioTime or API query params (exported for handlers).
func ParseBiometricDateTime(value string) (time.Time, error) {
	return parseBiometricDateTime(value)
}
