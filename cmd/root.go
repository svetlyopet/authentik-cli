package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "github.com/svetlyopet/authentik-cli/cmd/create"
	d "github.com/svetlyopet/authentik-cli/cmd/delete"
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/ak/authentik"
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
	addSubcommands()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", constants.CfgFilename))
}

func addSubcommands() {
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(c.CreateCmd())
	rootCmd.AddCommand(d.DeleteCmd())
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
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found. Run the config command to set a target.")
		os.Exit(1)
	}

	authentikUrl := viper.GetString("url")
	authentikToken := viper.GetString("token")

	ak.Repo = authentik.New(authentikUrl, authentikToken)
}
