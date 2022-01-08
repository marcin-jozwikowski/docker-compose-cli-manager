package command

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes a projects from saved list",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("name at least one project to remove")
		}

		for _, projectToRemove := range args {
			err := manager.GetConfigFile().DeleteProjectByName(projectToRemove)
			if err == nil {
				_, _ = fmt.Fprintf(mainWriter, "Project removed: "+projectToRemove+"\n")
			} else {
				return err
			}
		}

		return nil
	},
	ValidArgsFunction: projectNamesAutocompletion,
}

func init() {
	RootCommand.AddCommand(removeCommand)
}
