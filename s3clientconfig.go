package analytics

import "fmt"

// Given a config object as argument the function will set all zero-values to
// their defaults and return the modified object.
func makeS3ClientConfig(c S3ClientConfig) (S3ClientConfig, error) {
	if c.Stage == "" {
		c.Stage = "prod"
	}
	if c.Bucket == "" {
		c.Bucket = "fh-analytics-" + c.Stage
	}
	if c.Stream == "" {
		return c, fmt.Errorf("Stream should be provided")
	}

	return c, nil
}
