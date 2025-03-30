package s3client

import (
	"MediaDispenserGo/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

// Initialize inicializuje S3 klienta
func Initialize() {
	var err error
	Client, err = minio.New(config.AppConfig.S3.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.S3.AccessKey, config.AppConfig.S3.SecretKey, ""),
		Secure: config.AppConfig.S3.Secure,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}
}
