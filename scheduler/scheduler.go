package scheduler

import (
	"fmt"
	"time"

	"github.com/benn0x/backup-cloud/config"
)

var quit = make(chan struct{})

// Starts the scheduler and keeps an eye on any errors to return to the core program.
func StartScheduler(duration time.Duration) error {
	fmt.Printf("Beginning backup scheduler. Next backup in %g minutes\n", config.Config.BackupInterval)
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			err := schedulerRun()
			if err != nil {
				close(quit)
				return err
			}
		case <-quit:
			ticker.Stop()
			return nil
		}
	}
}
