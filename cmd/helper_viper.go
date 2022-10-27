package cmd

import (
	"github.com/spf13/viper"
)

func viperString(flag string) *string {
	if viper.GetString(flag) == "" {
		value := ""
		return &value
	}
	value := viper.GetString(flag)
	return &value
}

func viperInt(flag string) *int {
	value := viper.GetInt(flag)
	return &value
}

func viperBool(flag string) *bool {
	if !viper.GetBool(flag) {
		value := false
		return &value
	}
	value := viper.GetBool(flag)
	return &value
}
