package scheduler

import "github.com/benn0x/backup-cloud/backup"

// Executes all code that should run on a schedule
func schedulerRun() error {
	err := backup.StartBackup()
	if err != nil {
		return err
	}

	return nil
}
