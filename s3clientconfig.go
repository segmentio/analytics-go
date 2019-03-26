package analytics

import (
	"fmt"
	"strings"
	"time"
)

// Given a config object as argument the function will set all zero-values to
// their defaults and return the modified object.
func makeS3ClientConfig(c S3ClientConfig) (S3ClientConfig, error) {
	if c.Stage == "" {
		return c, fmt.Errorf("Stage should be provided - dev, prod, ci etc")
	}
	if c.Stage != strings.ToLower(c.Stage) {
		return c, fmt.Errorf("Stage should be lowercased")
	}

	if c.Bucket == "" {
		c.Bucket = "fh-analytics-" + c.Stage
	}

	if c.Stream == "" {
		return c, fmt.Errorf("Stream should be provided")
	}
	if c.Stream != strings.ToLower(c.Stream) {
		return c, fmt.Errorf("Stream should be lowercased")
	}

	if c.KeyConstructor == nil {
		c.KeyConstructor = func(now func() Time, uid func() string) string {
			ts := time.Time(now())
			return fmt.Sprintf(
				"analytics/%s/bulk/%s/json/%d/%02d/%02d/%02d/%d-%s.json.gz",
				strings.ToUpper(c.Stage),
				c.Stream,
				ts.Year(), ts.Month(), ts.Day(), ts.Hour(),
				ts.Unix(),
				uid(),
			)
		}
	}

	return c, nil
}
