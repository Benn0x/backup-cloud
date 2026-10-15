package backup

import (
	"archive/zip"
	"fmt"
	"path"
	"time"

	"github.com/benn0x/backup-cloud/config"
)

// Start the procedure that locally backs up a file/folder
func beginLocalBackup(files []BackupFile) error {
	if config.Config.Local.ArchiveBackup == true {
		backupName := path.Join(config.Config.Local.BackupLocation, fmt.Sprintf("%s_backup.zip", time.Now().Format("D2006_01_02T15_04_05")))
		err, archiveFile := CreateZipArchive(backupName)
		zipWriter := zip.NewWriter(archiveFile)
		if err != nil {
			return err
		}
		defer archiveFile.Close()
		defer zipWriter.Close()

		for _, fileToBackup := range files {
			file, path := fileToBackup.File, fileToBackup.Path
			defer file.Close()

			AddToZipArchive(file, path, zipWriter)
		}
	}

	return nil
}
