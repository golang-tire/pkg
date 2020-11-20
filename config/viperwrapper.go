package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func getViperString(key string, defaultValue string) (string, error) {
	v := viper.GetString(key)
	if v == "" {
		return defaultValue, nil
	}
	return v, nil
}

func getViperStringArray(key string, defaultValue []string) ([]string, error) {
	v := viper.GetStringSlice(key)
	if v == nil {
		return defaultValue, nil
	}
	return v, nil
}

func getViperInt64(key string, defaultValue int64) (int64, error) {
	v := viper.GetInt64(key)
	if v == 0 {
		return defaultValue, nil
	}
	return v, nil
}

func getViperFloat32(key string, defaultValue float32) (float32, error) {
	v := viper.GetFloat64(key)
	if v == 0 {
		return defaultValue, nil
	}
	return float32(v), nil
}

func getViperFloat64(key string, defaultValue float64) (float64, error) {
	v := viper.GetFloat64(key)
	if v == 0 {
		return defaultValue, nil
	}
	return v, nil
}

func getViperBool(key string, defaultValue bool) (bool, error) {
	v := viper.GetBool(key)
	if !v {
		return defaultValue, nil
	}
	return v, nil
}

func initViper(confName, ext, appName string, onChange func() error) error {
	viper.SetConfigName(confName)                          // name of config file (without extension)
	viper.SetConfigType(ext)                               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", appName))   // path to look for the config file in
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName)) // call multiple times to add many search paths
	viper.AddConfigPath(".")                               // optionally look for config in the working directory
	err := viper.ReadInConfig()                            // Find and read the config file
	if err != nil {                                        // Handle errors reading the config file
		return fmt.Errorf("fatal error config file: %s", err)
	}

	err = onChange()
	if err != nil {
		return err
	}

	go func() {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			err := onChange()
			if err != nil {
				panic(err)
			}
		})
	}()
	return err
}
