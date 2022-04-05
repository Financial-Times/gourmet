package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/gotha/gourmet/structloader"
)

type DotEnvFileConfigProvider struct {
	data map[string]string
}

// Get - returns env variable value if exists
func (p *DotEnvFileConfigProvider) Get(key string) (string, error) {
	val, exists := p.data[key]
	if !exists {
		return "", structloader.ErrValNotFound
	}
	return val, nil
}

func NewDotEnvFileConfigProvider(filePath string) (*DotEnvFileConfigProvider, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filePath, err)
	}
	lines := strings.Split(string(data), "\n")

	p := &DotEnvFileConfigProvider{data: make(map[string]string, 0)}
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		pos := strings.Index(line, "=")
		if pos == -1 {
			return nil, fmt.Errorf("malformed config file at line %d", i+1)
		}
		k := line[0:pos]
		v := line[pos+1:]
		p.data[k] = v

	}
	return p, nil
}

func NewDotEnvFileConfigLoader(filePath string) (*structloader.Loader, error) {
	p, err := NewDotEnvFileConfigProvider(filePath)
	if err != nil {
		return nil, err
	}
	return structloader.New(
		p,
		structloader.WithLoaderTagName("conf"),
	), nil
}
