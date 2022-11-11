package command

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var settingsGetCommand = &cobra.Command{
	Use:   "settings-get [key]",
	Short: "Show current setting value",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid arguments")
		}

		value, err := manager.GetConfigFile().GetSettingsEntry(args[0])
		if err == nil {
			fmt.Printf("Setting %s value: %s", args[0], value)
		}

		return err
	},
}

func init() {
	RootCommand.AddCommand(settingsGetCommand)
}
