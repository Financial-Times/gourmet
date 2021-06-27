package structloader

import (
	"encoding/json"
	"fmt"
	"testing"
)

type dataProviderMock struct {
	getF func(string) (string, error)
}

func (m *dataProviderMock) Get(s string) (string, error) {
	if m.getF != nil {
		return m.getF(s)
	}
	panic("dataProviderMock.Get not implemented")
}

func TestLoader(t *testing.T) {
	tests := []struct {
		name         string
		strct        interface{}
		dataProvider Provider
		expectedRes  interface{}
		expectedErr  error
	}{
		{
			name: "test that error is returned when the provided struct is not a pointer",
			strct: struct {
				test string
			}{},
			expectedErr: fmt.Errorf("struct is not a pointer"),
		},
		{
			name: "test that error is returned when one of the properties cannot be set but has conf tag",
			strct: &struct {
				test string `data:"TEST"`
			}{},
			expectedErr: fmt.Errorf("field 'test' cannot be set"),
		},
		{
			name: "test that error is returned if required value is not set and has default",
			strct: &struct {
				Test string `data:"TEST" required:"true"`
			}{},
			dataProvider: &dataProviderMock{
				getF: func(_ string) (string, error) {
					return "", ErrValNotFound
				},
			},
			expectedErr: fmt.Errorf("TEST is required, but is not set"),
		},

		{
			name: "test that default value is set when value is missing",
			strct: &struct {
				Test string `data:"TEST" default:"test1"`
			}{},
			dataProvider: &dataProviderMock{
				getF: func(_ string) (string, error) {
					return "", ErrValNotFound
				},
			},
			expectedRes: &struct {
				Test string
			}{Test: "test1"},
		},
		{
			name: "test that default value is overridden when actual config value is provided",
			strct: &struct {
				Test string `data:"TEST" default:"test1"`
			}{},
			dataProvider: &dataProviderMock{
				getF: func(_ string) (string, error) {
					return "test2", nil
				},
			},
			expectedRes: &struct {
				Test string
			}{Test: "test2"},
		},
		{
			name: "test that error is returned when the type is unsupported",
			strct: &struct {
				Test int64 `data:"TEST" default:"5"`
			}{},
			dataProvider: &dataProviderMock{
				getF: func(_ string) (string, error) {
					return "", ErrValNotFound
				},
			},
			expectedErr: fmt.Errorf("could not set value: unsupported config type for field 'Test'"),
		},
		{
			name: "test that int values are correctly converted",
			strct: &struct {
				Test int `data:"TEST" default:"5"`
			}{},
			dataProvider: &dataProviderMock{
				getF: func(_ string) (string, error) {
					return "666", nil
				},
			},
			expectedRes: &struct {
				Test int
			}{Test: 666},
		},
		{
			name: "test that bool values are correctly converted",
			strct: &struct {
				Test bool `data:"TEST" default:"false"`
			}{},
			dataProvider: &dataProviderMock{
				getF: func(_ string) (string, error) {
					return "true", nil
				},
			},
			expectedRes: &struct {
				Test bool
			}{Test: true},
		},
		{
			name: "test that everything works with nested fields",
			strct: &struct {
				Test        bool `data:"TEST" default:"false"`
				AnotherTest struct {
					MyConfig string `data:"MY_CONFIG"`
				}
			}{},
			dataProvider: &dataProviderMock{
				getF: func(key string) (string, error) {
					switch key {
					case "TEST":
						return "true", nil
					case "MY_CONFIG":
						return "myValue", nil
					}
					return "", ErrValNotFound
				},
			},
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

			loader := New(test.dataProvider)
			err := loader.Load(test.strct)

			if err != nil {
				if test.expectedErr == nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if err.Error() != test.expectedErr.Error() {
					t.Fatalf("expected error '%v', got '%v'", test.expectedErr, err)
				}
			}

			if test.expectedRes != nil {
				res, _ := json.Marshal(test.strct)
				expRes, _ := json.Marshal(test.expectedRes)
				if string(res) != string(expRes) {
					t.Fatalf("expected result '%v', got '%v'", test.expectedRes, res)
				}
			}
		})
	}
}
