package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

// The structure of the configuration if backupType is "cloudflare"
type cloudflareConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	AccountId       string `json:"accountId"`
}

// The structure of the configuration if backupType is "local"
type localConfig struct {
	BackupLocation string `json:"backupLocation"`
	ArchiveBackup  bool   `json:"archiveBackup"`
}

// The structure of a configuration
type config struct {
	BackupInterval float32          `json:"backupInterval"`
	BackupType     string           `json:"backupType"`
	BackupOrigin   string           `json:"backupOrigin"`
	Cloudflare     cloudflareConfig `json:"cloudflare"`
	Local          localConfig      `json:"local"`
}

// An empty config that can be accessed globally
var Config config

// Checks whether the configuration file actually exists
func configExists(configPath string) (error, bool) {
	if _, err := os.Open(configPath); err != nil && errors.Is(err, os.ErrNotExist) {
		return nil, false
	} else if err != nil {
		return err, false
	}

	return nil, true
}

// Creates a configuration file if it doesn't exist.
func createConfig(configPath string) (*os.File, error) {
	if file, err := os.Create(configPath); err != nil {
		return nil, err
	} else {
		return file, err
	}
}

func writeConfigTemplate(configFile *os.File) error {
	defaultConfig := config{
		BackupInterval: 60,
		BackupType:     "local",
		BackupOrigin:   "path_of_files_to_backup",
		Cloudflare: cloudflareConfig{
			AccessKeyId:     "account_key_id",
			SecretAccessKey: "secret_access_key",
			AccountId:       "account_id",
		},
		Local: localConfig{
			BackupLocation: "backup_output_location",
			ArchiveBackup:  true,
		},
	}

	b, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		configFile.Close()
		return err
	}

	if _, err := configFile.Write(b); err != nil {
		configFile.Close()
		return err
	}

	configFile.Close()
	return nil
}

// Loads the configuration into a globally accessible variable
func LoadConfig() error {
	configPath, err := FindPath("config")
	if err != nil {
		return err
	}

	dirPath, _ := path.Split(configPath)
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	if err, exists := configExists(configPath); err != nil {
		return err
	} else if exists == false {
		file, err := createConfig(configPath)

		if err != nil {
			return err
		}

		defer file.Close()
		err = writeConfigTemplate(file)

		if err != nil {
			return err
		}
	} else {
		file, err := os.Open(configPath)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&Config)
		file.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
