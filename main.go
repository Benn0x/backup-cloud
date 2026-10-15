package main

import (
	"fmt"
	"time"

	"github.com/benn0x/backup-cloud/config"
)

func main() {
	if err := config.RegisterPaths(); err != nil {
		fmt.Println(err)
	} else {
		// Begin the project startup
		if err := config.LoadConfig(); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Beginning backup scheduler. Next backup in %g minutes\n", config.Config.BackupInterval)
		tick := time.Tick(time.Second * time.Duration(config.Config.BackupInterval*60))
		for range tick {
			fmt.Println("Should start backup")
		}
	}
}
