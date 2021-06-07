package config

import "fmt"

// ErrConfigValNotFound - error that signifies that the required value was not
// found
var ErrConfigValNotFound = fmt.Errorf("config value not found")

// Provider - interface for configuration provider
type Provider interface {
	Get(string) (string, error)
}
