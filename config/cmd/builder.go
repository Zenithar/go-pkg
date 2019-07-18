package cmd

import (
	"fmt"
	"sort"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"go.zenithar.org/pkg/flags"
	"go.zenithar.org/pkg/log"
)

var configNewAsEnvFlag bool

// NewConfigCommand initialize a cobra config command tree
func NewConfigCommand(conf interface{}, envPrefix string) *cobra.Command {
	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// config
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Service Configuration",
	}

	// config new
	configNewCmd := &cobra.Command{
		Use:   "new",
		Short: "Initialize a default configuration",
		Run: func(cmd *cobra.Command, args []string) {
			defaults.SetDefaults(conf)

			if !configNewAsEnvFlag {
				btes, err := toml.Marshal(conf)
				if err != nil {
					log.Bg().Fatal("Error during configuration export", zap.Error(err))
				}
				fmt.Println(string(btes))
			} else {
				m := flags.AsEnvVariables(conf, upPrefix, true)
				keys := []string{}

				for k := range m {
					keys = append(keys, k)
				}

				sort.Strings(keys)
				for _, k := range keys {
					fmt.Printf("export %s=\"%s\"\n", k, m[k])
				}
			}
		},
	}

	// flags
	configNewCmd.Flags().BoolVar(&configNewAsEnvFlag, "env", false, "Print configuration as environment variable")
	configCmd.AddCommand(configNewCmd)

	// Return base command
	return configCmd
}
