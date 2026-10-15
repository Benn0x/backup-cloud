package backup

import (
	"fmt"
	"io"
	"path"
	"time"

	"github.com/benn0x/backup-cloud/config"
)

// Start the procedure that locally backs up a file/folder
func beginLocalBackup(files []BackupFile) error {
	if config.Config.Local.ArchiveBackup == true {
		backupName := path.Join(config.Config.Local.BackupLocation, fmt.Sprintf("%s_backup.zip", time.Now().Format("D2006_01_02T15_04_05")))
		err, archiveFile, zipWriter := CreateZipArchive(backupName)
		if err != nil {
			return err
		}
		defer archiveFile.Close()
		defer zipWriter.Close()

		for _, fileToBackup := range files {
			file, path := fileToBackup.File, fileToBackup.Path
			defer file.Close()
			writer, err := zipWriter.Create(path)
			if err != nil {
				return err
			}

			buffer := make([]byte, 1024*128)
			for {
				n, err := file.Read(buffer)
				if n > 0 {
					writer.Write(buffer[:n])
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					return err
				}

			}
		}

		fmt.Println("Finished backup")
	}
	return nil
}
