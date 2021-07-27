package structloader

import "fmt"

// ErrValNotFound - error that signifies that the required value was not
// found
var ErrValNotFound = fmt.Errorf("value not found")

// Provider - interface for loader data provider
type Provider interface {
	Get(string) (string, error)
}
