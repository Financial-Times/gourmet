package config

import "os"

// EnvConfigProvider - configuration data provider that loads values from env
// variables
type EnvConfigProvider struct{}

// Get - returns env variable value if exists
func (p *EnvConfigProvider) Get(key string) (string, error) {
	val, exists := os.LookupEnv(key)
	if !exists {
		return "", ErrConfigValNotFound
	}
	return val, nil
}

func NewEnvConfigLoader() *Loader {
	return NewLoader(&EnvConfigProvider{})
}
