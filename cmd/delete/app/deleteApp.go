package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/provider"
)

func DeleteAppCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "app",
		Short: "Delete an app",
		Long: fmt.Sprintf(`Deletes an application and provider pair in Authentik.

Examples:
  # Delete an app
  %s delete app example-app`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err := deleteApp(name)
			cobra.CheckErr(err)
		},
	}

	return c
}

func deleteApp(name string) error {
	app, err := core.GetApplication(name)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
		return err
	}

	if err = core.DeleteApplication(app.Name, app.Slug); err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
		return err
	}

	if err = provider.DeleteProvider(app.Name, app.ProviderPK); err != nil {
		return err
	}

	return nil
}
