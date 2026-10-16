package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/benn0x/backup-cloud/config"
	"github.com/benn0x/backup-cloud/files"
)

// How files are structured, simple but effective
type BackupFile struct {
	File *os.File
	Path string
}

// Writes the contents of a file to the zip given
func AddToZipArchive(archiveFile *os.File, archivePathName string, archiveWriter *zip.Writer) error {
	buffer := make([]byte, 1024*64)

	relative, err := filepath.Rel(config.Config.BackupOrigin, archivePathName)
	if err != nil {
		return err
	}

	writer, err := archiveWriter.Create(relative)
	if err != nil {
		return err
	}

	for {
		n, err := archiveFile.Read(buffer)
		if n > 0 {
			if _, err := writer.Write(buffer[:n]); err != nil {
				return nil
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Create and return zip writer and file for easier usage
func CreateZipArchive(location string) (error, *os.File) {
	dir, _ := path.Split(location)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err, nil
	}

	archive, err := os.Create(location)
	if err != nil {
		return err, nil
	}

	return nil, archive
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
