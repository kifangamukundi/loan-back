package jobs

import (
	"log"

	"github.com/robfig/cron/v3"
)

// InitializeJobs sets up and starts the cron scheduler
func InitializeJobs() {
	c := cron.New()

	// Use the "EVERY_MINUTE" schedule for the PingServer job
	// _, err := c.AddFunc(schedules.Schedules["EVERY_5_MINUTES"], PingServer)
	// if err != nil {
	// 	log.Fatalf("Failed to schedule PingServer job: %v", err)
	// }

	// Start the cron scheduler
	c.Start()

	log.Println("Cron jobs initialized and scheduler started")
}
