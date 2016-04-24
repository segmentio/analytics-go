package analytics

import "fmt"

// Returned by the `NewWithConfig` function when the one of the configuration
// fields was set to an impossible value (like a negative duration).
type ConfigError struct {

	// A human-readable message explaining why the configuration field's value
	// is invalid.
	Reason string

	// The name of the configuration field that was carrying an invalid value.
	Field string

	// The value of the configuration field that caused the error.
	Value interface{}
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("analytics.NewWithConfig: %s (analytics.Config.%s: %#v)", e.Reason, e.Field, e.Value)
}

// Instances of this type are used to represent errors returned when a field was
// no initialize properly in a structure passed as argument to one of the
// functions of this package.
type FieldError struct {

	// The human-readable representation of the type of structure that wasn't
	// initialized properly.
	Type string

	// The name of the field that wasn't properly initialized.
	Name string

	// The value of the field that wasn't properly initialized.
	Value interface{}
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s.%s: invalid field value: %#v", e.Type, e.Name, e.Value)
}
