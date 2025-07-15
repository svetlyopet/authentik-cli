package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "github.com/svetlyopet/authentik-cli/cmd/create"
	d "github.com/svetlyopet/authentik-cli/cmd/delete"
	g "github.com/svetlyopet/authentik-cli/cmd/get"
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/ak/authentik"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   constants.CmdName,
	Short: "A CLI tool for managing Authentik",
	Long: fmt.Sprintf(`%s is a CLI that enables managing resources
in Authentik deployments to create more complex setups
and enable multi-tenancy and automation`, constants.CmdName),
}

func Execute() {
	loadConfig()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommands() {
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(c.CreateCmd())
	rootCmd.AddCommand(d.DeleteCmd())
	rootCmd.AddCommand(g.GetCmd())
}

func loadConfig() {
	cobra.OnInitialize(readCfgFile)
	addSubcommands()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", constants.CfgFilename))
}

func readCfgFile() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory
		viper.AddConfigPath(homeDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(constants.CfgFilename)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	authentikUrl := viper.GetString("url")
	authentikToken := viper.GetString("token")

	ak.Repo = authentik.New(authentikUrl, authentikToken)
}
