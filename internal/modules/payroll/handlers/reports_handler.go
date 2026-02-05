package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/payroll/services"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ReportsHandler handles payroll report-related HTTP requests
type ReportsHandler struct {
	runService     *services.PayrollRunService
	payslipService *services.PayslipService
}

// NewReportsHandler creates a new handler instance
func NewReportsHandler() *ReportsHandler {
	return &ReportsHandler{
		runService:     services.NewPayrollRunService(),
		payslipService: services.NewPayslipService(),
	}
}

// GetPayrollReport retrieves payroll report data
func (h *ReportsHandler) GetPayrollReport(c *gin.Context) {
	tenantID := getTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))

	if payMonth == 0 {
		payMonth = int(time.Now().Month())
	}
	if payYear == 0 {
		payYear = time.Now().Year()
	}

	// Get payroll run for the period
	runs, total, err := h.runService.ListPayrollRuns(tenantID, "", payYear, payMonth, 1, 1)
	if err != nil || total == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no payroll run found for this period"})
		return
	}

	run := runs[0]
	summary, _ := h.runService.GetPayrollRunSummary(run.ID, tenantID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetDepartmentCostReport retrieves payroll cost by department
func (h *ReportsHandler) GetDepartmentCostReport(c *gin.Context) {
	tenantID := getTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))

	if payMonth == 0 {
		payMonth = int(time.Now().Month())
	}
	if payYear == 0 {
		payYear = time.Now().Year()
	}

	// Get payroll run for the period
	runs, total, err := h.runService.ListPayrollRuns(tenantID, "", payYear, payMonth, 1, 1)
	if err != nil || total == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no payroll run found for this period"})
		return
	}

	run := runs[0]
	employees, _, _ := h.runService.GetEmployees(run.ID, "", "", nil, 1, 10000)

	// Group by department
	departmentCosts := make(map[string]map[string]interface{})
	for _, emp := range employees {
		dept := emp.Department
		if dept == "" {
			dept = "Unassigned"
		}

		if _, exists := departmentCosts[dept]; !exists {
			departmentCosts[dept] = map[string]interface{}{
				"department":      dept,
				"employeeCount":   0,
				"totalGross":      0.0,
				"totalDeductions": 0.0,
				"totalNet":        0.0,
			}
		}

		departmentCosts[dept]["employeeCount"] = departmentCosts[dept]["employeeCount"].(int) + 1
		departmentCosts[dept]["totalGross"] = departmentCosts[dept]["totalGross"].(float64) + emp.GrossSalary
		departmentCosts[dept]["totalDeductions"] = departmentCosts[dept]["totalDeductions"].(float64) + emp.TotalDeductions
		departmentCosts[dept]["totalNet"] = departmentCosts[dept]["totalNet"].(float64) + emp.NetPay
	}

	// Convert to slice
	var deptList []map[string]interface{}
	for _, v := range departmentCosts {
		deptList = append(deptList, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"payMonth":    payMonth,
			"payYear":     payYear,
			"departments": deptList,
		},
	})
}

// ExportPayrollReport exports payroll report to Excel
func (h *ReportsHandler) ExportPayrollReport(c *gin.Context) {
	tenantID := getTenantID(c)
	payMonth, _ := strconv.Atoi(c.Query("month"))
	payYear, _ := strconv.Atoi(c.Query("year"))
	format := c.DefaultQuery("format", "xlsx")

	if payMonth == 0 {
		payMonth = int(time.Now().Month())
	}
	if payYear == 0 {
		payYear = time.Now().Year()
	}

	// Get payroll run
	runs, total, err := h.runService.ListPayrollRuns(tenantID, "", payYear, payMonth, 1, 1)
	if err != nil || total == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no payroll run found for this period"})
		return
	}

	run := runs[0]
	summary, _ := h.runService.GetPayrollRunSummary(run.ID, tenantID)
	employees, _, _ := h.runService.GetEmployees(run.ID, "", "", nil, 1, 10000)

	if format != "xlsx" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}

	// Create Excel file
	f := excelize.NewFile()

	// Summary Sheet
	summarySheet := "Summary"
	f.SetSheetName("Sheet1", summarySheet)
	f.SetCellValue(summarySheet, "A1", "Payroll Summary Report")
	f.SetCellValue(summarySheet, "A2", fmt.Sprintf("Period: %s %d", time.Month(payMonth).String(), payYear))

	row := 4
	f.SetCellValue(summarySheet, fmt.Sprintf("A%d", row), "Metric")
	f.SetCellValue(summarySheet, fmt.Sprintf("B%d", row), "Value")
	row++

	if financialSummary, ok := summary["financialSummary"].(map[string]interface{}); ok {
		for key, val := range financialSummary {
			f.SetCellValue(summarySheet, fmt.Sprintf("A%d", row), key)
			f.SetCellValue(summarySheet, fmt.Sprintf("B%d", row), val)
			row++
		}
	}

	// Employee Details Sheet
	detailSheet := "Employee Details"
	f.NewSheet(detailSheet)

	headers := []string{"Employee ID", "Name", "Department", "Basic Salary", "Gross Salary", "Deductions", "Net Pay"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(detailSheet, cell, h)
	}

	for i, emp := range employees {
		row := i + 2
		f.SetCellValue(detailSheet, fmt.Sprintf("A%d", row), emp.EmployeeCode)
		f.SetCellValue(detailSheet, fmt.Sprintf("B%d", row), emp.EmployeeName)
		f.SetCellValue(detailSheet, fmt.Sprintf("C%d", row), emp.Department)
		f.SetCellValue(detailSheet, fmt.Sprintf("D%d", row), emp.BasicSalary)
		f.SetCellValue(detailSheet, fmt.Sprintf("E%d", row), emp.GrossSalary)
		f.SetCellValue(detailSheet, fmt.Sprintf("F%d", row), emp.TotalDeductions)
		f.SetCellValue(detailSheet, fmt.Sprintf("G%d", row), emp.NetPay)
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("payroll_report_%d_%02d.xlsx", payYear, payMonth)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// GetPayrollKPIs retrieves payroll KPIs
func (h *ReportsHandler) GetPayrollKPIs(c *gin.Context) {
	tenantID := getTenantID(c)
	payYear, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	// Get all runs for the year
	runs, _, err := h.runService.ListPayrollRuns(tenantID, "", payYear, 0, 1, 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalGross, totalNet, totalDeductions, totalEmployerContribs float64
	var totalEmployees int
	monthlyData := make([]map[string]interface{}, 0)

	for _, run := range runs {
		totalGross += run.TotalGross
		totalNet += run.TotalNet
		totalDeductions += run.TotalDeductions
		totalEmployerContribs += run.TotalEmployerContributions
		totalEmployees += run.TotalEmployees

		monthlyData = append(monthlyData, map[string]interface{}{
			"month":                   run.PayMonth,
			"totalGross":              run.TotalGross,
			"totalNet":                run.TotalNet,
			"totalDeductions":         run.TotalDeductions,
			"employerContributions":   run.TotalEmployerContributions,
			"totalEmployees":          run.TotalEmployees,
		})
	}

	avgMonthlyCost := 0.0
	if len(runs) > 0 {
		avgMonthlyCost = (totalGross + totalEmployerContribs) / float64(len(runs))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"year":                    payYear,
			"totalRuns":               len(runs),
			"totalGross":              totalGross,
			"totalNet":                totalNet,
			"totalDeductions":         totalDeductions,
			"totalEmployerContributions": totalEmployerContribs,
			"averageMonthlyCost":      avgMonthlyCost,
			"monthlyBreakdown":        monthlyData,
		},
	})
}

// GetBankFileReport generates bank file report for disbursement
func (h *ReportsHandler) GetBankFileReport(c *gin.Context) {
	tenantID := getTenantID(c)
	runID, err := strconv.ParseUint(c.Param("runId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run id"})
		return
	}

	run, err := h.runService.GetPayrollRun(uint(runID), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payroll run not found"})
		return
	}

	employees, _, _ := h.runService.GetEmployees(run.ID, "", "", nil, 1, 10000)

	// Create Excel file for bank
	f := excelize.NewFile()
	sheet := "Bank Transfer"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"S/N", "Employee Name", "Bank Name", "Account Number", "Amount", "Reference"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, h)
	}

	reference := fmt.Sprintf("PAY/%d/%02d", run.PayYear, run.PayMonth)
	for i, emp := range employees {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), emp.EmployeeName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "N/A") // Would need bank info
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "N/A") // Would need account number
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), emp.NetPay)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("%s/%s", reference, emp.EmployeeCode))
	}

	// Add totals row
	totalRow := len(employees) + 2
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalRow), run.TotalNet)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("bank_transfer_%d_%02d.xlsx", run.PayYear, run.PayMonth)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
