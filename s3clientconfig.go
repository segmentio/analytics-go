package analytics

import (
	"fmt"
	"strings"
	"time"
)

// MB is the number of bytes in one megabyte.
const MB = 1024 * 1024

// Given a config object as argument the function will set all zero-values to
// their defaults and return the modified object.
func makeS3ClientConfig(c S3ClientConfig) (S3ClientConfig, error) {
	c.Config = makeConfig(c.Config)
	if c.S3.Stage != strings.ToLower(c.S3.Stage) {
		return c, fmt.Errorf("Stage should be lowercased")
	}
	if c.S3.Stage == "" {
		c.S3.Stage = "dev"
	}

	if c.S3.MaxBatchBytes == 0 {
		c.S3.MaxBatchBytes = 20 * MB
	}

	if c.S3.Bucket == "" {
		c.S3.Bucket = "fh-analytics-" + c.S3.Stage
	}

	if c.S3.Stream != strings.ToLower(c.S3.Stream) {
		return c, fmt.Errorf("Stream should be lowercased")
	}

	if c.S3.Stream == "" {
		c.S3.Stream = "haring"
	}

	if c.S3.KeyConstructor == nil {
		c.S3.KeyConstructor = func(now func() Time, uid func() string) string {
			ts := time.Time(now())
			return fmt.Sprintf(
				"analytics/%s/bulk/%s/json/%d/%02d/%02d/%02d/%d-%s.json.gz",
				strings.ToUpper(c.S3.Stage),
				c.S3.Stream,
				ts.Year(), ts.Month(), ts.Day(), ts.Hour(),
				ts.Unix(),
				uid(),
			)
		}
	}

	return c, nil
}
