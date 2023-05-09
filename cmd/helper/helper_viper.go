// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

func ViperString(flag string) *string {
	if viper.GetString(flag) == "" {
		value := ""
		return &value
	}
	value := viper.GetString(flag)
	return &value
}

func ViperInt(flag string) *int {
	value := viper.GetInt(flag)
	return &value
}

func ViperBool(flag string) *bool {
	if !viper.GetBool(flag) {
		value := false
		return &value
	}
	value := viper.GetBool(flag)
	return &value
}

func ViperStringSlice(flag string) []string {
	value := viper.GetStringSlice(flag)
	if len(value) == 0 {
		return []string{}
	}
	return value
}

func ViperStringSliceMap(flag string) (map[string]string, error) {
	m := make(map[string]string)
	values := viper.GetStringSlice(flag)

	for _, v := range values {
		// Expecting each value to be in "a=1" format
		s := strings.SplitN(v, "=", 2)
		if len(s) != 2 {
			return nil, errors.New("invalid env var or configuration file.")
		}
		m[s[0]] = s[1]
		s1 := s[0]
		s2 := s[1]

		if s1 == "" {
			return m, errors.New("invalid input provided on left side of a mapping inside the configuration file")
		}

		if s2 == "" {
			return m, errors.New("invalid input provided on right side of a mapping inside the configuration file")
		}

	}
	return m, nil
}

func ViperMapKeyValuePair(flag string) (string, string, error) {
	//m := make(map[string]string)
	var s1 string
	var s2 string
	values := viper.GetStringSlice(flag)

	for _, v := range values {
		// Expecting each value to be in "a=1" format
		s := strings.SplitN(v, "=", 2)
		if len(s) != 2 {
			return "", "", errors.New("invalid env var")
		}
		s1 := s[0]
		s2 := s[1]

		if s1 == "" {
			return "", "", errors.New("invalid source provided on varsets-map left side")
		}

		if s2 == "" {
			return "", "", errors.New("invalid destination provided on varsets-map right side")
		}

		return s1, s2, nil
	}
	return s1, s2, nil
}
