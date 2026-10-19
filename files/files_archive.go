package files

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/benn0x/backup-cloud/config"
)

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
