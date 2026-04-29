//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/app"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/database"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/models"
)

func main() {
	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}

	var count int64
	database.GetDB().Model(&models.EmailLog{}).Count(&count)
	fmt.Printf("Total Email Logs: %d\n", count)

	var logs []models.EmailLog
	database.GetDB().Order("created_at desc").Limit(5).Find(&logs)
	for _, l := range logs {
		fmt.Printf("ID: %d, Month: %s, Status: %s, URL: %q, Size: %d, Type: %s\n", l.ID, l.ReportMonth, l.Status, l.DownloadURL, l.FileSize, l.FileType)
	}
}
