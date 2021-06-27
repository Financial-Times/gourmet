package structloader

import (
	"fmt"
	"reflect"
	"strconv"
)

// Loader - service that loads data into struct based on rules specified in tag
type Loader struct {
	dataProvider Provider
	opts         loaderOptions
}

// New - create new Loader
func New(dp Provider, opts ...LoaderOption) *Loader {
	options := newLoaderOptions()
	options.applyOptions(opts...)

	return &Loader{
		dataProvider: dp,
		opts:         *options,
	}
}

// Load - fetches values from data provider and fills them in a struct
func (l *Loader) Load(i interface{}) error {
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("%s is not a pointer", v.Kind())
	}
	v = v.Elem()
	t := reflect.TypeOf(i).Elem()

	for i := 0; i < t.NumField(); i++ {
		err := l.process(v.Field(i), t.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Loader) process(v reflect.Value, f reflect.StructField) error {
	if v.Type().Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			err := l.process(v.Field(i), f.Type.Field(i))
			if err != nil {
				return err
			}
		}
		return nil
	}

	confKey, ok := f.Tag.Lookup(l.opts.loaderTagName)
	if !ok {
		// if there is no tag that we are looking for, just skip it
		return nil
	}

	if !v.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", f.Name)
	}

	val, err := l.getValue(confKey, f)
	if err != nil {
		return fmt.Errorf("could not retrieve value: %w", err)
	}

	requiredVal, required := f.Tag.Lookup(l.opts.requiredTagName)
	if val == "" {
		if required && requiredVal == "true" {
			return fmt.Errorf("%s is required, but is not set", confKey)
		}
		// if it is not required and is not set - nothing to do
		return nil
	}

	err = setValue(v, f.Name, val)
	if err != nil {
		return fmt.Errorf("could not set value: %w", err)
	}

	return nil
}

func (l *Loader) getValue(key string, f reflect.StructField) (string, error) {
	newVal, err := l.dataProvider.Get(key)
	if err != nil && err != ErrValNotFound {
		return "", fmt.Errorf("error loading value for '%s': %w", key, err)
	}

	defaultVal, defaultExists := f.Tag.Lookup(l.opts.defaultTagName)
	if newVal == "" && defaultExists {
		newVal = defaultVal
	}
	return newVal, nil
}

func setValue(v reflect.Value, fieldName string, newVal string) error {
	switch v.Type().Kind() {
	case reflect.String:
		v.SetString(newVal)
	case reflect.Int:
		fieldVal, err := strconv.ParseInt(newVal, 10, 64)
		if err != nil {
			return fmt.Errorf("could not cast value '%v' to int: %w", newVal, err)
		}
		v.SetInt(fieldVal)
	case reflect.Bool:
		fieldVal, err := strconv.ParseBool(newVal)
		if err != nil {
			return fmt.Errorf("could not cast value '%v' to bool", newVal)
		}
		v.SetBool(fieldVal)
	default:
		return fmt.Errorf("unsupported config type for field '%s'", fieldName)
	}
	return nil
}
