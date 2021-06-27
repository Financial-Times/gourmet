package config

import (
	"os"

	"github.com/Financial-Times/gourmet/structloader"
)

// EnvConfigProvider - configuration data provider that loads values from env
// variables
type EnvConfigProvider struct{}

// Get - returns env variable value if exists
func (p *EnvConfigProvider) Get(key string) (string, error) {
	val, exists := os.LookupEnv(key)
	if !exists {
		return "", structloader.ErrValNotFound
	}
	return val, nil
}

func NewEnvConfigLoader() *structloader.Loader {
	return structloader.New(
		&EnvConfigProvider{},
		structloader.WithLoaderTagName("conf"),
	)
}
