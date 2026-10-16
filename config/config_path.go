package config

import (
	"fmt"
	"runtime"
)

// Define the structure that Paths should be passed as
type PathDef struct {
	Name         string
	Path         string
	ApplicableOS string
}

// Define a new list for path definitions
var paths []PathDef

// RegisterPaths registers the paths specific to an operating system.
// Returns any errors that happen when adding paths to the listing.
func RegisterPaths() error {
	// List of paths and their operating system
	// used for registration and shouldn't
	// be changed willy nilly

	pathlist := []PathDef{
		{Name: "config", Path: "C:\\Program Files\\Cloud Backups\\config.json", ApplicableOS: "windows"},
		{Name: "temp", Path: "C:\\Program Files\\Cloud Backups\\temp\\", ApplicableOS: "windows"},
		{Name: "logs", Path: "C:\\Program Files\\Cloud Backups\\logs\\", ApplicableOS: "windows"},

		{Name: "config", Path: "/etc/cloud-backup/config.json", ApplicableOS: "linux"},
		{Name: "temp", Path: "/etc/cloud-backup/temp/", ApplicableOS: "linux"},
		{Name: "logs", Path: "/var/cloud-backup/", ApplicableOS: "linux"},
	}

	// Loops through the pathlist and registers them
	// depending on which operating system is found
	for _, path := range pathlist {
		if path.ApplicableOS != runtime.GOOS {
			continue
		}

		if err := addPath(path.Name, path.Path); err != nil {
			return err
		}
	}

	if len(pathlist) == 0 {
		return fmt.Errorf("Your device is running an unsupported OS, check out documentation for more information.")
	}

	return nil
}

// AddPath registers a path under a name.config.json
// It returns an error if the name given already has a path assigned to it.
func addPath(name, path string) error {
	for _, existing := range paths {
		if existing.Name == name {
			return fmt.Errorf("path %s already exists", name)
		}
	}

	paths = append(paths, PathDef{
		Name: name,
		Path: path,
	})
	return nil
}

// FindPath returns an assigned path based off the name.
// It returns an error if there isn't a path match found.
func FindPath(name string) (string, error) {
	for _, path := range paths {
		if path.Name == name {
			return path.Path, nil
		}
	}

	return "", fmt.Errorf("path %s not found", name)
}
