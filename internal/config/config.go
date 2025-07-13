package config

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"gopkg.in/yaml.v3"
)

type Config struct {
	URL   string
	Token string
}

var (
	cfgFilepath string
)

func Set() error {
	var c Config
	var err error

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the URL for your target Authentik instace: ")
	scanner.Scan()
	c.URL = scanner.Text()

	fmt.Println("Enter the API token for authenticating against the target Authentik instance: ")
	scanner.Scan()
	c.Token = scanner.Text()

	file, err := os.Create(cfgFilepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close() //nolint

	yamlConfig, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	_, err = file.Write(yamlConfig)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	viper.AutomaticEnv()

	homeDir := viper.GetString("HOME")
	cfgFilepath = fmt.Sprintf("%s/%s", homeDir, constants.CfgFilename)
}
