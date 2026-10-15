package files

import (
	"io/fs"
	"path/filepath"
)

// Recursively walks the directory given to find files
func Walk(dir string) (error, []string) {
	paths := []string{}
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			paths = append(paths, path)
		}

		return nil
	}); err != nil {
		return err, nil
	}

	return nil, paths
}
