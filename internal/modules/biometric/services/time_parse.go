package services

import (
	"fmt"
	"os"
	"time"
)

// getBiometricLocation returns the timezone used for parsing naive BioTime timestamps.
// Priority:
// 1) BIOMETRIC_TIMEZONE
// 2) APP_TIMEZONE
// 3) Africa/Dar_es_Salaam (default business timezone)
// 4) UTC fallback
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

// parseBiometricDateTime parses timestamp strings from BioTime/query params.
// Naive layouts are interpreted in business timezone; RFC3339 keeps embedded zone.
func parseBiometricDateTime(value string) (time.Time, error) {
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
