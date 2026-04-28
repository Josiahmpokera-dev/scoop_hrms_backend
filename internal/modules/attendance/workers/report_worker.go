package workers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/attendance/services"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/queue"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/pkg/mailer"
)

type ReportTask struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

type ReportWorker struct {
	queueClient *queue.RabbitMQClient
	genService  *services.ReportGenerationService
	mailer      *mailer.Mailer
}

func NewReportWorker() (*ReportWorker, error) {
	qClient, err := queue.NewRabbitMQClient() // Reusing existing RabbitMQ logic
	if err != nil {
		return nil, err
	}

	return &ReportWorker{
		queueClient: qClient,
		genService:  services.NewReportGenerationService(),
		mailer:      mailer.NewMailer(),
	}, nil
}

func (w *ReportWorker) Start() error {
	ch := w.queueClient.GetChannel()
	
	// Ensure queue exists
	qName := "attendance_reports"
	_, err := ch.QueueDeclare(
		qName, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var task ReportTask
			if err := json.Unmarshal(d.Body, &task); err != nil {
				log.Printf("Error unmarshaling report task: %v", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("Generating monthly report for %d-%02d", task.Year, task.Month)
			
			excelData, err := w.genService.GenerateMonthlyReportExcel(task.Year, task.Month)
			if err != nil {
				log.Printf("Failed to generate report: %v", err)
				d.Nack(false, true)
				continue
			}

			// Send email
			to := []string{config.AppConfig.SMTP.HREmail}
			subject := fmt.Sprintf("Monthly Attendance Report - %d/%02d", task.Year, task.Month)
			body := fmt.Sprintf("<p>Attached is the monthly attendance report for <b>%d-%02d</b>.</p><p>Regards,<br>HRMS System</p>", task.Year, task.Month)
			
			attachments := []mailer.Attachment{
				{
					Filename: fmt.Sprintf("Attendance_Report_%d_%02d.xlsx", task.Year, task.Month),
					Content:  excelData,
				},
			}

			if err := w.mailer.SendEmail(to, subject, body, attachments); err != nil {
				log.Printf("Failed to send report email: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Successfully sent monthly report for %d-%02d", task.Year, task.Month)
			d.Ack(false)
		}
	}()

	return nil
}
