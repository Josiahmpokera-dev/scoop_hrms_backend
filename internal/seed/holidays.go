package seed

import (
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	holidayModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/leave/models"
)

// tanzaniaHoliday is a helper struct to define holiday templates.
type tanzaniaHoliday struct {
	Name        string
	Month       time.Month
	Day         int
	Type        string // Public, Restricted
	Description string
}

// ─────────────────────────────────────────────────────────────────────────────
// SECTION 1 — Fixed-date public holidays (same date every year)
// ─────────────────────────────────────────────────────────────────────────────

var fixedTanzaniaHolidays = []tanzaniaHoliday{
	// January
	{Name: "New Year's Day", Month: time.January, Day: 1,
		Type: "Public", Description: "New Year's Day – National public holiday"},
	{Name: "Zanzibar Revolution Day", Month: time.January, Day: 12,
		Type: "Public", Description: "Anniversary of the 1964 Zanzibar Revolution"},

	// April
	{Name: "Karume Day", Month: time.April, Day: 7,
		Type: "Public", Description: "Sheikh Abeid Amani Karume Day – First President of Zanzibar"},
	{Name: "Union Day", Month: time.April, Day: 26,
		Type: "Public", Description: "Tanganyika and Zanzibar united to form the United Republic of Tanzania in 1964"},

	// May
	{Name: "International Workers' Day", Month: time.May, Day: 1,
		Type: "Public", Description: "International Labour Day / May Day"},

	// July
	{Name: "Saba Saba Day", Month: time.July, Day: 7,
		Type: "Public", Description: "Saba Saba (7/7) – International Trade Fair Day, anniversary of TANU founding (1954)"},

	// August
	{Name: "Nane Nane Day (Farmers' Day)", Month: time.August, Day: 8,
		Type: "Public", Description: "Nane Nane (8/8) – National Farmers' / Peasants' Day"},

	// October
	{Name: "Nyerere Day", Month: time.October, Day: 14,
		Type: "Public", Description: "Mwalimu Julius Kambarage Nyerere Memorial Day (d. 1999)"},

	// December
	{Name: "Independence Day", Month: time.December, Day: 9,
		Type: "Public", Description: "Tanganyika Independence from Britain in 1961"},
	{Name: "Christmas Day", Month: time.December, Day: 25,
		Type: "Public", Description: "Christmas Day – National public holiday"},
	{Name: "Boxing Day", Month: time.December, Day: 26,
		Type: "Public", Description: "Boxing Day – National public holiday"},
}

// ─────────────────────────────────────────────────────────────────────────────
// SECTION 2 — Easter-based movable holidays (Christian)
// ─────────────────────────────────────────────────────────────────────────────

// getEasterDate calculates Easter Sunday for a given year using the
// Anonymous Gregorian algorithm (Meeus/Jones/Butcher).
func getEasterDate(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// getChristianMovableHolidays returns Easter-based holidays for a given year.
func getChristianMovableHolidays(year int) []tanzaniaHoliday {
	easter := getEasterDate(year)
	goodFriday := easter.AddDate(0, 0, -2)
	easterMonday := easter.AddDate(0, 0, 1)

	return []tanzaniaHoliday{
		{Name: "Good Friday", Month: goodFriday.Month(), Day: goodFriday.Day(),
			Type: "Public", Description: "Good Friday – Christian holiday observed nationally"},
		{Name: "Easter Sunday", Month: easter.Month(), Day: easter.Day(),
			Type: "Public", Description: "Easter Sunday – Christian holiday observed nationally"},
		{Name: "Easter Monday", Month: easterMonday.Month(), Day: easterMonday.Day(),
			Type: "Public", Description: "Easter Monday – Christian holiday observed nationally"},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SECTION 3 — Islamic holidays (lunar calendar — approximate Gregorian dates)
//
// Islamic holidays follow the Hijri calendar and shift ~10-11 days earlier
// each Gregorian year. Exact dates depend on moon sighting, so these are
// best-estimate dates from astronomical predictions. Update if the official
// government gazette announces different dates.
// ─────────────────────────────────────────────────────────────────────────────

type islamicHolidayDates struct {
	EidAlFitrMonth  time.Month
	EidAlFitrDay    int
	EidAlFitr2Month time.Month // second day
	EidAlFitr2Day   int
	EidAlAdhaMonth  time.Month
	EidAlAdhaDay    int
	MaulidMonth     time.Month
	MaulidDay       int
}

// Approximate Islamic holiday dates per year.
// Sources: timeanddate.com, islamicfinder.org — best astronomical estimates.
var islamicHolidaysByYear = map[int]islamicHolidayDates{
	2025: {
		EidAlFitrMonth: time.March, EidAlFitrDay: 31,
		EidAlFitr2Month: time.April, EidAlFitr2Day: 1,
		EidAlAdhaMonth: time.June, EidAlAdhaDay: 7,
		MaulidMonth: time.September, MaulidDay: 5,
	},
	2026: {
		EidAlFitrMonth: time.March, EidAlFitrDay: 20,
		EidAlFitr2Month: time.March, EidAlFitr2Day: 21,
		EidAlAdhaMonth: time.May, EidAlAdhaDay: 27,
		MaulidMonth: time.August, MaulidDay: 26,
	},
	2027: {
		EidAlFitrMonth: time.March, EidAlFitrDay: 10,
		EidAlFitr2Month: time.March, EidAlFitr2Day: 11,
		EidAlAdhaMonth: time.May, EidAlAdhaDay: 16,
		MaulidMonth: time.August, MaulidDay: 15,
	},
	2028: {
		EidAlFitrMonth: time.February, EidAlFitrDay: 27,
		EidAlFitr2Month: time.February, EidAlFitr2Day: 28,
		EidAlAdhaMonth: time.May, EidAlAdhaDay: 5,
		MaulidMonth: time.August, MaulidDay: 4,
	},
	2029: {
		EidAlFitrMonth: time.February, EidAlFitrDay: 15,
		EidAlFitr2Month: time.February, EidAlFitr2Day: 16,
		EidAlAdhaMonth: time.April, EidAlAdhaDay: 24,
		MaulidMonth: time.July, MaulidDay: 24,
	},
	2030: {
		EidAlFitrMonth: time.February, EidAlFitrDay: 4,
		EidAlFitr2Month: time.February, EidAlFitr2Day: 5,
		EidAlAdhaMonth: time.April, EidAlAdhaDay: 13,
		MaulidMonth: time.July, MaulidDay: 14,
	},
}

// getIslamicHolidays returns Islamic holidays for a given year.
// Returns empty if the year is not in the lookup table.
func getIslamicHolidays(year int) []tanzaniaHoliday {
	dates, ok := islamicHolidaysByYear[year]
	if !ok {
		return nil
	}

	return []tanzaniaHoliday{
		{Name: "Eid al-Fitr (Eid El-Fitry)", Month: dates.EidAlFitrMonth, Day: dates.EidAlFitrDay,
			Type: "Public", Description: "Eid al-Fitr – End of Ramadan, first day. Date is approximate (depends on moon sighting)"},
		{Name: "Eid al-Fitr (Day 2)", Month: dates.EidAlFitr2Month, Day: dates.EidAlFitr2Day,
			Type: "Public", Description: "Eid al-Fitr – End of Ramadan, second day. Date is approximate (depends on moon sighting)"},
		{Name: "Eid al-Adha (Eid El-Hajj)", Month: dates.EidAlAdhaMonth, Day: dates.EidAlAdhaDay,
			Type: "Public", Description: "Eid al-Adha – Festival of Sacrifice. Date is approximate (depends on moon sighting)"},
		{Name: "Maulid Day (Prophet's Birthday)", Month: dates.MaulidMonth, Day: dates.MaulidDay,
			Type: "Public", Description: "Maulid – Birthday of Prophet Muhammad (Mawlid an-Nabi). Date is approximate (depends on moon sighting)"},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SECTION 4 — Seeder runner
// ─────────────────────────────────────────────────────────────────────────────

// RunTanzaniaHolidays seeds all Tanzania public holidays for the current year
// and the next year. It is idempotent — existing holidays (matched by name + date)
// are skipped.
//
// Includes:
//   - 11 fixed-date holidays (New Year, Zanzibar Revolution, Karume Day, Union Day,
//     Workers' Day, Saba Saba, Nane Nane, Nyerere Day, Independence Day, Christmas, Boxing Day)
//   - 3 Christian movable holidays (Good Friday, Easter Sunday, Easter Monday)
//   - 4 Islamic holidays (Eid al-Fitr x2, Eid al-Adha, Maulid Day) — approximate dates
//
// Total: ~18 holidays per year.
func RunTanzaniaHolidays() {
	db := database.GetDB()
	if db == nil {
		log.Println("⚠️  Skipping holiday seed: database not available")
		return
	}

	currentYear := time.Now().Year()
	years := []int{currentYear, currentYear + 1}

	totalCreated := 0
	totalSkipped := 0

	for _, year := range years {
		// Combine all holiday sources
		allHolidays := make([]tanzaniaHoliday, len(fixedTanzaniaHolidays))
		copy(allHolidays, fixedTanzaniaHolidays)
		allHolidays = append(allHolidays, getChristianMovableHolidays(year)...)
		allHolidays = append(allHolidays, getIslamicHolidays(year)...)

		for _, h := range allHolidays {
			date := time.Date(year, h.Month, h.Day, 0, 0, 0, 0, time.UTC)

			// Check if already exists (by name + date) — idempotent
			var count int64
			db.Model(&holidayModels.Holiday{}).
				Where("name = ? AND DATE(date) = DATE(?)", h.Name, date).
				Count(&count)

			if count > 0 {
				totalSkipped++
				continue
			}

			desc := h.Description
			holiday := &holidayModels.Holiday{
				Name:          h.Name,
				Date:          date,
				Type:          h.Type,
				IsFloater:     false,
				Location:      `["All Locations"]`,
				ApplicableFor: `["All Employees"]`,
				Description:   &desc,
			}

			if err := db.Create(holiday).Error; err != nil {
				log.Printf("⚠️  Failed to seed holiday %s (%s): %v", h.Name, date.Format("2006-01-02"), err)
				continue
			}

			totalCreated++
		}
	}

	log.Printf("🇹🇿 Tanzania holidays seeded: %d created, %d already existed (years: %d–%d)",
		totalCreated, totalSkipped, years[0], years[len(years)-1])
}
