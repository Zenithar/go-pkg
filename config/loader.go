package config

import (
	"fmt"
	"os"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"go.zenithar.org/pkg/flags"
	"go.zenithar.org/pkg/log"
)

// Load a config
// Apply defaults first, then environment, then finally config file.
func Load(conf interface{}, envPrefix string, cfgFile string) error {
	// Apply defaults first
	defaults.SetDefaults(conf)

	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// Overrides with environment
	for k := range flags.AsEnvVariables(conf, "", false) {
		envName := fmt.Sprintf("%s_%s", upPrefix, k)
		log.CheckErr("Unable to bind environment variable", viper.BindEnv(strings.ToLower(strings.Replace(k, "_", ".", -1)), envName), zap.String("var", envName))
	}

	// Apply file settings
	switch {
	case cfgFile != "":
		// If the config file doesn't exists, let's exit
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			return xerrors.Errorf("Unable to open non-existing file '%s': %w", cfgFile, err)
		}

		log.Bg().Info("Load settings from file", zap.String("path", cfgFile))

		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return xerrors.Errorf("Unable to decode config file '%s': %w", cfgFile, err)
		}
	}

	// Update viper values
	if err := viper.Unmarshal(conf); err != nil {
		return xerrors.Errorf("Unable to apply config '%s': %w", cfgFile, err)
	}

	// No error
	return nil
}
