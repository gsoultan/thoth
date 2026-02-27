package core

// S3Config contains settings for S3-compatible storage.
type S3Config struct {
	// Endpoint is the S3 endpoint URL.
	Endpoint string
	// AccessKeyID is the access key for S3.
	AccessKeyID string
	// SecretAccessKey is the secret key for S3.
	SecretAccessKey string
	// UseSSL enables SSL/TLS for S3 connections.
	UseSSL bool
	// Region is the S3 region.
	Region string
}
