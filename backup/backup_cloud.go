package backup

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	conf "github.com/benn0x/backup-cloud/config"
)

// Cleans up temporary archive after the upload.
func cleanupTempFile() error {
	fmt.Println("Beginning cleanup of temporary archives.")
	tempPath, err := conf.FindPath("temp")
	if err != nil {
		return err
	}

	if err := os.Remove(tempPath); err != nil {
		return err
	}

	fmt.Println("Finished cleanup of temporary archives.")
	return nil
}

// Create the temporary archive used for an upload
func createTempLocalBackup(files []BackupFile) (error, string) {
	tempPath, err := conf.FindPath("temp")
	if err != nil {
		return err, ""
	}
	backupName := path.Join(tempPath, fmt.Sprintf("%s_backup.zip", time.Now().Format("D2006_01_02T15_04_05")))
	err, archiveFile := CreateZipArchive(backupName)
	if err != nil {
		return err, ""
	}

	zipWriter := zip.NewWriter(archiveFile)
	defer archiveFile.Close()
	defer zipWriter.Close()

	for _, fileToBackup := range files {
		file, path := fileToBackup.File, fileToBackup.Path
		defer file.Close()

		AddToZipArchive(file, path, zipWriter)
	}

	return nil, backupName
}

// Starts the procedure to backup to the cloud, uses the S3 api
func beginCloudBackup(files []BackupFile) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		24*time.Hour,
	)
	defer cancel()

	err, backupPath := createTempLocalBackup(files)
	if err != nil {
		if err := cleanupTempFile(); err != nil {
			return err
		}
		return err
	}

	f, err := os.Open(backupPath)
	if err != nil {
		if err := cleanupTempFile(); err != nil {
			return err
		}
		return err
	}

	file := bufio.NewReader(f)
	defer f.Close()

	accessKeyId, accountId, secretAccessKey := conf.Config.Cloudflare.AccessKeyId, conf.Config.Cloudflare.AccountId, conf.Config.Cloudflare.SecretAccessKey

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKeyId,
				secretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		if err := cleanupTempFile(); err != nil {
			return err
		}
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(
			fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		)
	})

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024
		u.Concurrency = 4
		u.LeavePartsOnError = false
	})

	key := filepath.Base(backupPath)
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(conf.Config.Cloudflare.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		if err := cleanupTempFile(); err != nil {
			return err
		}
		return err
	}

	fmt.Println("Upload successful.", result)
	if err := cleanupTempFile(); err != nil {
		return err
	}

	return nil
}
