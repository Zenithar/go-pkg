package config

import (
	"fmt"
	"os"
	"strings"

	"go.zenithar.org/pkg/log"
	"go.zenithar.org/pkg/types"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Load a config
// Apply defaults first, then environment, then finally config file.
func Load(conf interface{}, envPrefix string, cfgFile string) error {
	// Apply defaults first
	defaults.SetDefaults(conf)

	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// Overrides with environment
	env, err := types.AsEnvVariables(conf, "", false)
	if err != nil {
		return err
	}

	for k := range env {
		envName := fmt.Sprintf("%s_%s", upPrefix, k)
		log.CheckErr("Unable to bind environment variable", viper.BindEnv(strings.ToLower(strings.Replace(k, "_", ".", -1)), envName), zap.String("var", envName))
	}

	// Apply file settings
	if cfgFile != "" {
		// If the config file doesn't exists, let's exit
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return fmt.Errorf("Unable to open non-existing file '%s': %w", cfgFile, err)
		}

		log.Bg().Info("Load settings from file", zap.String("path", cfgFile))

		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("Unable to decode config file '%s': %w", cfgFile, err)
		}
	}

	// Update viper values
	if err := viper.Unmarshal(conf); err != nil {
		return fmt.Errorf("Unable to apply config '%s': %w", cfgFile, err)
	}

	// No error
	return nil
}
