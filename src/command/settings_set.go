package command

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var settingsSetCommand = &cobra.Command{
	Use:   "settings-set [key] [velue]",
	Short: "Sets a settings value",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid arguments")
		}

		err := manager.GetConfigFile().StoreSettingsEntry(args[0], args[1])
		if err == nil {
			fmt.Println("Settings set")
		}

		return err
	},
}

func init() {
	RootCommand.AddCommand(settingsSetCommand)
}
