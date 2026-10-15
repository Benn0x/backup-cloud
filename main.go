package main

import (
	"fmt"
	"time"

	"github.com/benn0x/backup-cloud/config"
	"github.com/benn0x/backup-cloud/scheduler"
)

func main() {
	if err := config.RegisterPaths(); err != nil {
		fmt.Println(err)
	} else {
		// Begin the project startup
		if err := config.LoadConfig(); err != nil {
			fmt.Println(err)
		}

		if err := scheduler.StartScheduler(time.Second * time.Duration(config.Config.BackupInterval*60)); err != nil {
			fmt.Println(err)
		}
	}
}
