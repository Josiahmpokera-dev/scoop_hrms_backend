package scheduler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/queue"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
)

type ReportTask struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

// Scheduler handles scheduled jobs
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	// Schedule monthly report on the 1st of every month at 00:00
	// It will generate for the previous month
	_, err := s.cron.AddFunc("0 0 1 * *", func() {
		now := time.Now()
		// Go to previous month
		prev := now.AddDate(0, -1, 0)
		
		task := ReportTask{
			Year:  prev.Year(),
			Month: int(prev.Month()),
		}

		if err := publishReportTask(task); err != nil {
			log.Printf("Failed to schedule monthly report: %v", err)
		} else {
			log.Printf("Scheduled automated monthly report for %d-%02d", task.Year, task.Month)
		}
	})

	if err != nil {
		log.Printf("Error scheduling monthly report: %v", err)
	}

	s.cron.Start()
	log.Println("✅ Scheduler started successfully")
}

func publishReportTask(task ReportTask) error {
	qClient, err := queue.NewRabbitMQClient()
	if err != nil {
		return err
	}
	defer qClient.Close()

	body, _ := json.Marshal(task)
	ch := qClient.GetChannel()
	
	qName := "attendance_reports"
	_, err = ch.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return ch.Publish("", qName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// ManualTrigger trigger the report manually for a specific date
func ManualTrigger(year, month int) error {
	task := ReportTask{Year: year, Month: month}
	return publishReportTask(task)
}
