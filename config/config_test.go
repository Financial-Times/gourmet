package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func setVars(vars map[string]string) {
	for k, v := range vars {
		os.Setenv(k, v)
	}
}

func clearVars(vars map[string]string) {
	for k := range vars {
		os.Unsetenv(k)
	}
}

func TestConfigCreate(t *testing.T) {
	tests := []struct {
		name        string
		conf        interface{}
		envVars     map[string]string
		expectedRes interface{}
		expectedErr error
	}{
		{
			name: "test that error is returned when the provided struct is not a pointer",
			conf: struct {
				test string
			}{},
			expectedErr: fmt.Errorf("struct is not a pointer"),
		},
		{
			name: "test that error is returned when one of the properties cannot be set but has conf tag",
			conf: &struct {
				test string `conf:"TEST"`
			}{},
			expectedErr: fmt.Errorf("field 'test' cannot be set"),
		},
		{
			name: "test that error is returned if required value is not set and has default",
			conf: &struct {
				Test string `conf:"TEST" required:"true"`
			}{},
			expectedErr: fmt.Errorf("TEST is required, but is not set"),
		},
		{
			name: "test that default value is set when value is missing",
			conf: &struct {
				Test string `conf:"TEST" default:"test1"`
			}{},
			expectedRes: &struct {
				Test string
			}{Test: "test1"},
		},
		{
			name:    "test that default value is overridden when actual config value is provided",
			envVars: map[string]string{"TEST": "test2"},
			conf: &struct {
				Test string `conf:"TEST" default:"test1"`
			}{},
			expectedRes: &struct {
				Test string
			}{Test: "test2"},
		},
		{
			name: "test that error is returned when the type is unsupported",
			conf: &struct {
				Test int64 `conf:"TEST" default:"5"`
			}{},
			expectedErr: fmt.Errorf("could not set value: unsupported config type for field 'Test'"),
		},
		{
			name:    "test that int values are correctly converted",
			envVars: map[string]string{"TEST": "666"},
			conf: &struct {
				Test int `conf:"TEST" default:"5"`
			}{},
			expectedRes: &struct {
				Test int
			}{Test: 666},
		},
		{
			name:    "test that bool values are correctly converted",
			envVars: map[string]string{"TEST": "true"},
			conf: &struct {
				Test bool `conf:"TEST" default:"false"`
			}{},
			expectedRes: &struct {
				Test bool
			}{Test: true},
		},
		{
			name: "test that everything works with nested fields",
			envVars: map[string]string{
				"TEST":      "true",
				"MY_CONFIG": "myValue",
			},
			conf: &struct {
				Test        bool `conf:"TEST" default:"false"`
				AnotherTest struct {
					MyConfig string `conf:"MY_CONFIG"`
				}
			}{},
			expectedRes: &struct {
				Test        bool
				AnotherTest struct {
					MyConfig string
				}
			}{
				Test: true,
				AnotherTest: struct {
					MyConfig string
				}{
					MyConfig: "myValue",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setVars(test.envVars)
			defer clearVars(test.envVars)

			confLoader := NewLoader(&EnvConfigProvider{})
			err := confLoader.Load(test.conf)

			if err != nil {
				if test.expectedErr == nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if err.Error() != test.expectedErr.Error() {
					t.Fatalf("expected error '%v', got '%v'", test.expectedErr, err)
				}
			}

			if test.expectedRes != nil {
				res, _ := json.Marshal(test.conf)
				expRes, _ := json.Marshal(test.expectedRes)
				if string(res) != string(expRes) {
					t.Fatalf("expected result '%v', got '%v'", test.expectedRes, test.conf)
				}
			}
		})
	}
}
