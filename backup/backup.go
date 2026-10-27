package backup

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/benn0x/backup-cloud/config"
	"github.com/benn0x/backup-cloud/files"
)

// How files are structured, simple but effective
type BackupFile struct {
	File *os.File
	Path string
}

// Should return a list of opened files for easier reading
func OpenFiles(paths []string) ([]BackupFile, error) {
	files := []BackupFile{}

	for _, p := range paths {
		dir, _ := path.Split(p)
		if dir != config.Config.BackupLocation {
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

	split := strings.Split(config.Config.BackupType, ",")
	for _, configType := range split {
		switch configType {
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
	}

	return nil
}
