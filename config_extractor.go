package goconf

// Hector Oliveros - 2019
// hector.oliveros.leon@gmail.com

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Extract extracts environment variables and sets configuration values for a given sets of structs.
func Extract(args ExtractorArgs) error {
	args.Options.mergeWithDefault()
	if len(args.Options.EnvFile) > 0 {
		err := loadAllEnvValFromEnvFile(args)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(args.Configs); i += 2 {
		err := extractEnvByConfigIdx(i, args)
		if err != nil {
			return err
		}
	}
	return nil
}

// extractEnvByConfigIdx extracts environment variables and sets configuration values for a
// given struct based on the specified index.
// It takes the following parameters:
// - idx: The index of the configuration object within the Configs slice.
// - args: The ExtractorArgs containing the configuration options and values.
// The function iterates over the fields of the specified struct and sets their values based
// on the corresponding environment variables.
// The function returns an error if there is any issue setting the configuration values.
func extractEnvByConfigIdx(idx int, args ExtractorArgs) error {
	v := reflect.ValueOf(args.Configs[idx]).Elem()
	vPtr := v
	// if its a pointer, resolve its value
	if v.Kind() == reflect.Ptr {
		vPtr = reflect.Indirect(v)
	}
	for i := 0; i < vPtr.NumField(); i++ {
		f := v.Field(i)
		// make sure that this field is defined, and can be changed.
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		envName, err := getEnvName(v, i, idx, args)
		if err != nil {
			continue
		}
		envVal, _ := getEnvValuesFromSources(envName, args.Options)
		if envVal == "" {
			continue
		}
		err = setConfigFieldValue(f, envVal)
		if err != nil {
			return err
		}
	}
	return nil
}

// getEnvName retrieves the environment variable name for a given field within a struct.
// It takes the following parameters:
// - v: The reflect.Value of the struct.
// - i: The index of the field within the struct.
// - idx: The current index of the Configs slice.
// - args: The ExtractorArgs containing the configuration options and values.
// It returns the environment variable name as a string and an error.
func getEnvName(v reflect.Value, i int, idx int, args ExtractorArgs) (string, error) {
	envName, haveTagEnv := v.Type().Field(i).Tag.Lookup(envTagName)
	prefix := ""
	if args.Configs[idx+1] != "" {
		var ok bool
		prefix, ok = args.Configs[idx+1].(string)
		if !ok {
			panic("ERROR: Invalid configuration. Expected string")
		}
		if !strings.HasSuffix(prefix, "_") {
			prefix = prefix + "_"
		}
	}
	if !haveTagEnv {
		if args.Options.OmitNotTagged {
			return "", fmt.Errorf("unable to get the value name: %v", v)
		}
		envName = v.Type().Field(i).Name
	}
	return prefix + envName, nil
}

// setConfigFieldValue handles setting the value of a field in a struct based on its kind.
// It takes the following parameters:
// - field: The reflect.Value of the field to be set.
// - value: The value to set the field to.
// If the value cannot be parsed or the field kind is not supported, an error is returned.
func setConfigFieldValue(field reflect.Value, value interface{}) error {
	switch field.Kind() {
	case reflect.Int, reflect.Int64:
		valueInt, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer: %v of type %v", value, reflect.TypeOf(value))
		}
		field.SetInt(valueInt)
	case reflect.String:
		field.SetString(value.(string))
	case reflect.Bool:
		val, err := strconv.ParseBool(value.(string))
		if err != nil {
			panic(fmt.Errorf("invalid boolean value %v", value))
		}
		field.SetBool(val)
	default:
		panic("unhandled default case")
	}
	return nil
}
