package backup

import (
	"archive/zip"
	"fmt"
	"path"
	"time"

	"github.com/benn0x/backup-cloud/config"
	files_service "github.com/benn0x/backup-cloud/files"
)

// Start the procedure that locally backs up a file/folder
func beginLocalBackup(files []BackupFile) error {
	backupName := path.Join(config.Config.BackupLocation, fmt.Sprintf("%s_backup.zip", time.Now().Format("D2006_01_02T15_04_05")))
	err, archiveFile := files_service.CreateZipArchive(backupName)
	zipWriter := zip.NewWriter(archiveFile)
	if err != nil {
		return err
	}
	defer archiveFile.Close()
	defer zipWriter.Close()

	for _, fileToBackup := range files {
		file, path := fileToBackup.File, fileToBackup.Path
		defer file.Close()

		files_service.AddToZipArchive(file, path, zipWriter)
	}

	return nil
}
