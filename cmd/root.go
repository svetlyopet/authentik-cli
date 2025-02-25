package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "authentik-cli",
	Short: "A CLI tool for managing Authentik",
	Long: `authentik-cli is a CLI that enables managing resources
in Authentik deployments to create more complex setups
and enable multi-tenancy and automation`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", constants.CfgFilename))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".authentik-cli" (without extension).
		viper.AddConfigPath(homeDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(constants.CfgFilename)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
