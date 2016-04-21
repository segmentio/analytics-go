package analytics

import "fmt"

type MissingFieldError struct {
	Type string
	Name string
}

func (e MissingFieldError) Error() string {
	return fmt.Sprintf("%s.%s must be defined", self.Type, self.Name)
}
