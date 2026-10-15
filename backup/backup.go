package backup

import (
	"archive/zip"
	"fmt"
	"os"
	"path"

	"github.com/benn0x/backup-cloud/config"
	"github.com/benn0x/backup-cloud/files"
)

// How files are structured, simple but effective
type BackupFile struct {
	File *os.File
	Path string
}

// Create and return zip writer and file for easier usage
func CreateZipArchive(location string) (error, *os.File, *zip.Writer) {
	dir, _ := path.Split(location)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err, nil, nil
	}

	archive, err := os.Create(location)
	if err != nil {
		return err, nil, nil
	}

	zipWriter := zip.NewWriter(archive)

	return nil, archive, zipWriter
}

// Should return a list of opened files for easier reading
func OpenFiles(paths []string) ([]BackupFile, error) {
	files := []BackupFile{}

	for _, p := range paths {
		dir, _ := path.Split(p)
		dir2, _ := path.Split(config.Config.Local.BackupLocation)
		if dir != dir2 {
			file, err := os.Open(p)
			if err != nil {
				return nil, err
			}

			files = append(files, BackupFile{
				Path: p,
				File: file,
			})
		}
	}

	return files, nil
}

// Starts a backup depending on how the program is configured.
func StartBackup() error {
	err, FilePaths := files.Walk(config.Config.BackupOrigin)
	if err != nil {
		return err
	}

	files, err := OpenFiles(FilePaths)
	if err != nil {
		return err
	}

	switch config.Config.BackupType {
	case "local":
		if err := beginLocalBackup(files); err != nil {
			return err
		}
	case "cloudflare":
		if err := beginCloudBackup(files); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported BackupType in your configuration, refer to the documentation.")
	}

	return nil
}
