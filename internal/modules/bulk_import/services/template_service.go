package services

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// header styling shared across all templates
func headerStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"2F5496"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	return s
}

func requiredStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"C00000"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	return s
}

func sampleStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "808080", Italic: true},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	return s
}

func instructionStyle(f *excelize.File) int {
	s, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "333333"},
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})
	return s
}

// ── Employee Template ──────────────────────────────────────────────────

type colDef struct {
	header   string
	width    float64
	required bool
	sample   string
}

func GenerateEmployeeTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Employees"
	f.SetSheetName("Sheet1", sheet)

	cols := []colDef{
		{"Employee ID *", 15, true, "EMP001"},
		{"First Name *", 18, true, "John"},
		{"Middle Name", 15, false, "Michael"},
		{"Last Name *", 18, true, "Doe"},
		{"Work Email", 28, false, "john.doe@company.com"},
		{"Personal Email", 28, false, "john@gmail.com"},
		{"Phone Number", 18, false, "+255712345678"},
		{"Gender", 12, false, "Male"},
		{"Date of Birth", 16, false, "1990-05-15"},
		{"Nationality", 15, false, "Tanzanian"},
		{"Marital Status", 15, false, "Single"},
		{"Department Code", 18, false, "ENG"},
		{"Position Code", 18, false, "SE-001"},
		{"Employment Type", 18, false, "full_time"},
		{"Hire Date", 14, false, "2024-01-15"},
		{"Grade", 10, false, "L3"},
		{"Shift", 12, false, "Day"},
		{"Manager Employee ID", 20, false, "EMP000"},
		{"Salary", 14, false, "3500000"},
		{"Currency", 12, false, "TZS"},
		{"Emergency Contact Name", 22, false, "Jane Doe"},
		{"Emergency Contact Phone", 22, false, "+255712345679"},
		{"Emergency Contact Relation", 24, false, "Spouse"},
		{"Notes", 25, false, "Transferred from Arusha branch"},
	}

	hStyle := headerStyle(f)
	rStyle := requiredStyle(f)
	smpl := sampleStyle(f)

	for i, c := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, c.header)
		if c.required {
			f.SetCellStyle(sheet, cell, cell, rStyle)
		} else {
			f.SetCellStyle(sheet, cell, cell, hStyle)
		}
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, c.width)

		sampleCell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, sampleCell, c.sample)
		f.SetCellStyle(sheet, sampleCell, sampleCell, smpl)
	}

	f.SetRowHeight(sheet, 1, 30)

	addInstructionsSheet(f, "Instructions", []string{
		"EMPLOYEE BULK UPLOAD INSTRUCTIONS",
		"",
		"1. Fill in data starting from row 2 (row 1 is the header — do NOT modify it).",
		"2. Columns marked with * (red headers) are REQUIRED.",
		"3. Row 2 has sample data — replace it with your actual data.",
		"",
		"FIELD GUIDE:",
		"• Employee ID — Unique identifier (e.g., EMP001). Must not already exist.",
		"• First Name, Last Name — Required text fields.",
		"• Work Email — Must be unique across the system.",
		"• Gender — Male, Female, or Other.",
		"• Date of Birth / Hire Date — Format: YYYY-MM-DD (e.g., 1990-05-15).",
		"• Nationality — e.g., Tanzanian, Kenyan.",
		"• Marital Status — Single, Married, Divorced, Widowed.",
		"• Department Code — Must match an existing department code (e.g., ENG, FIN).",
		"• Position Code — Must match an existing position code (e.g., SE-001).",
		"• Employment Type — full_time, part_time, contract, or intern.",
		"• Grade — Grade level (e.g., L1, L2, L3).",
		"• Shift — Day, Night, or Rotating.",
		"• Manager Employee ID — Employee ID of the reporting manager.",
		"• Salary — Numeric value (e.g., 3500000).",
		"• Currency — 3-letter code (e.g., TZS, USD, KES).",
		"• Emergency Contact — Name, phone, and relation of emergency contact.",
		"",
		"TIPS:",
		"• Upload departments and positions FIRST, then employees.",
		"• You can leave optional columns blank.",
		"• Maximum 500 rows per upload.",
	})

	return f, nil
}

// ── Department Template ────────────────────────────────────────────────

func GenerateDepartmentTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Departments"
	f.SetSheetName("Sheet1", sheet)

	cols := []colDef{
		{"Code *", 12, true, "ENG"},
		{"Name *", 28, true, "Engineering"},
		{"Description", 35, false, "Software engineering department"},
		{"Level", 16, false, "department"},
		{"Department Type", 18, false, "core"},
		{"Parent Department Code", 22, false, ""},
		{"Employee Capacity", 18, false, "50"},
	}

	hStyle := headerStyle(f)
	rStyle := requiredStyle(f)
	smpl := sampleStyle(f)

	for i, c := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, c.header)
		if c.required {
			f.SetCellStyle(sheet, cell, cell, rStyle)
		} else {
			f.SetCellStyle(sheet, cell, cell, hStyle)
		}
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, c.width)

		sampleCell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, sampleCell, c.sample)
		f.SetCellStyle(sheet, sampleCell, sampleCell, smpl)
	}

	// Add more sample rows
	samples := [][]interface{}{
		{"FIN", "Finance", "Finance and accounting", "department", "support", "", 30},
		{"HR", "Human Resources", "People operations", "department", "support", "", 15},
		{"MKT", "Marketing", "Sales and marketing", "department", "core", "", 25},
		{"IT", "Information Technology", "IT infrastructure & support", "department", "support", "", 20},
	}
	for r, row := range samples {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+3)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, smpl)
		}
	}

	f.SetRowHeight(sheet, 1, 30)

	addInstructionsSheet(f, "Instructions", []string{
		"DEPARTMENT BULK UPLOAD INSTRUCTIONS",
		"",
		"1. Fill in data starting from row 2 (row 1 is the header — do NOT modify it).",
		"2. Columns marked with * (red headers) are REQUIRED.",
		"3. Rows 2-6 have sample data — replace with your actual data.",
		"",
		"FIELD GUIDE:",
		"• Code — Short unique code (e.g., ENG, FIN, HR). Max 20 chars.",
		"• Name — Full department name. Max 100 chars.",
		"• Description — Optional description.",
		"• Level — One of: company, business_unit, department, team.",
		"• Department Type — One of: core, support, operational, strategic.",
		"• Parent Department Code — Code of the parent department (for hierarchy).",
		"  Leave blank for top-level departments.",
		"• Employee Capacity — Expected number of employees.",
		"",
		"TIPS:",
		"• Upload parent departments FIRST (rows without parent), then children.",
		"• The system processes rows top-to-bottom, so parent codes referenced",
		"  in later rows should appear in earlier rows.",
		"• Maximum 200 rows per upload.",
	})

	return f, nil
}

// ── Position Template ──────────────────────────────────────────────────

func GeneratePositionTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Positions"
	f.SetSheetName("Sheet1", sheet)

	cols := []colDef{
		{"Code *", 14, true, "SE-001"},
		{"Title *", 28, true, "Software Engineer"},
		{"Grade", 12, false, "L3"},
		{"Department Code", 18, false, "ENG"},
		{"Reports To Position Code", 24, false, ""},
		{"Budgeted Headcount", 20, false, "10"},
		{"Key Competencies", 35, false, "Go, PostgreSQL, Docker"},
	}

	hStyle := headerStyle(f)
	rStyle := requiredStyle(f)
	smpl := sampleStyle(f)

	for i, c := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, c.header)
		if c.required {
			f.SetCellStyle(sheet, cell, cell, rStyle)
		} else {
			f.SetCellStyle(sheet, cell, cell, hStyle)
		}
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, c.width)

		sampleCell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, sampleCell, c.sample)
		f.SetCellStyle(sheet, sampleCell, sampleCell, smpl)
	}

	samples := [][]interface{}{
		{"SSE-001", "Senior Software Engineer", "L4", "ENG", "SE-001", 5, "Go, Architecture, Leadership"},
		{"ACC-001", "Accountant", "L2", "FIN", "", 8, "Accounting, Excel, Tax"},
		{"HR-001", "HR Officer", "L2", "HR", "", 4, "Recruitment, Policies, HRIS"},
		{"MKT-001", "Marketing Specialist", "L2", "MKT", "", 6, "Digital Marketing, Analytics"},
		{"IT-001", "IT Support Specialist", "L2", "IT", "", 5, "Networking, Hardware, Support"},
	}
	for r, row := range samples {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+3)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, smpl)
		}
	}

	f.SetRowHeight(sheet, 1, 30)

	addInstructionsSheet(f, "Instructions", []string{
		"POSITION BULK UPLOAD INSTRUCTIONS",
		"",
		"1. Fill in data starting from row 2 (row 1 is the header — do NOT modify it).",
		"2. Columns marked with * (red headers) are REQUIRED.",
		"3. Rows 2-7 have sample data — replace with your actual data.",
		"",
		"FIELD GUIDE:",
		"• Code — Short unique code (e.g., SE-001, ACC-001). Max 20 chars.",
		"• Title — Full position title (e.g., Software Engineer). Max 100 chars.",
		"• Grade — Grade level (e.g., L1, L2, L3, L4).",
		"• Department Code — Must match an existing department code.",
		"  Upload departments BEFORE positions.",
		"• Reports To Position Code — Code of the position this reports to.",
		"  Leave blank for top-level positions.",
		"• Budgeted Headcount — Expected number of people in this position.",
		"• Key Competencies — Comma-separated skills/competencies.",
		"",
		"TIPS:",
		"• Upload departments FIRST, then positions, then employees.",
		"• Maximum 200 rows per upload.",
	})

	return f, nil
}

// ── Helper ─────────────────────────────────────────────────────────────

func addInstructionsSheet(f *excelize.File, name string, lines []string) {
	f.NewSheet(name)
	iStyle := instructionStyle(f)

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "2F5496"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	f.SetColWidth(name, "A", "A", 90)
	for i, line := range lines {
		cell := fmt.Sprintf("A%d", i+1)
		f.SetCellValue(name, cell, line)
		if i == 0 {
			f.SetCellStyle(name, cell, cell, titleStyle)
			f.SetRowHeight(name, i+1, 30)
		} else {
			f.SetCellStyle(name, cell, cell, iStyle)
		}
	}
}
