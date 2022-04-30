package cmd

import (
	"fmt"
	"testing"
)

func Test_NewArgsDataProvider(t *testing.T) {

	tests := []struct {
		name          string
		inputArgs     []string
		expectedArgs  map[string]interface{}
		expectedError error
	}{
		{
			name:      "test that argument value is parsed",
			inputArgs: []string{"--test", "value"},
			expectedArgs: map[string]interface{}{
				"test": "value",
			},
		},
		{
			name:      "test multiple values",
			inputArgs: []string{"--test", "value", "--test2", "value2"},
			expectedArgs: map[string]interface{}{
				"test":  "value",
				"test2": "value2",
			},
		},
		{
			name:      "test bool value as last argument",
			inputArgs: []string{"--test", "value", "--test2"},
			expectedArgs: map[string]interface{}{
				"test":  "value",
				"test2": "true",
			},
		},
		{
			name:      "test bool value as middle argument",
			inputArgs: []string{"--test", "value", "--test2", "--test3", "value3"},
			expectedArgs: map[string]interface{}{
				"test":  "value",
				"test2": "true",
				"test3": "value3",
			},
		},
		{
			name:      "test error is returned when value is missing",
			inputArgs: []string{"--test", "value", "--test2", "--test3", "value3"},
			expectedArgs: map[string]interface{}{
				"missing": "",
			},
			expectedError: fmt.Errorf("value not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := NewArgsDataProvider(WithCustomArgs(test.inputArgs))

			for k, v := range test.expectedArgs {
				res, err := provider.Get(k)
				if test.expectedError != nil {
					if err == nil || test.expectedError.Error() != err.Error() {
						t.Fatalf("expected error '%v' got '%v'", test.expectedError, err)
					}
				} else {
					if res != v {
						t.Fatalf("expected value '%v' got '%v'", v, res)
					}
				}
			}
		})
	}
}
