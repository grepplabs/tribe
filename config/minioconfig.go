package config

import "github.com/spf13/pflag"

type MinioConfig struct {
	flagBase

	Endpoint        string
	UseSSL          bool
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	BucketLocation  string
}

func NewMinioConfig() *MinioConfig {
	return &MinioConfig{}
}

func (c *MinioConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Endpoint, "minio-endpoint", "localhost:9000", "Service endpoint")
		c.flagSet.BoolVar(&c.UseSSL, "minio-use-ssl", false, "Set this value to true to enable secure (HTTPS) access")
		c.flagSet.StringVar(&c.AccessKeyID, "minio-access-key-id", "minioadmin", "Access key is like user ID that uniquely identifies your account")
		c.flagSet.StringVar(&c.SecretAccessKey, "minio-secret-access-key", "minioadmin123", "Secret key is the password to your account")
		c.flagSet.StringVar(&c.BucketName, "minio-bucket-name", "tribe", "Bucket name")
		c.flagSet.StringVar(&c.BucketLocation, "minio-bucket-location", "eu-central-1", "Region the bucket resides in")
	}
	return c.flagSet
}
