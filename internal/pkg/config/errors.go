package config

import "fmt"

type ErrorSet struct {
	Errors []error `json:"errors"`
}
func (err ErrorSet) Error() string {
	return fmt.Sprintf("set of %d errors: %+v", len(err.Errors), err.Errors)
}
