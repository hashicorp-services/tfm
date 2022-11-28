package cmd

import (
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
