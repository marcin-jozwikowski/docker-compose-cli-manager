package command

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add [project name] [docker-compose-file]",
	Short: "Adds docker-compose file to available files list",
	Long: `Adds docker-compose file to available files list

If no project name is provided current directory name is used.
If no file is provided it look for one in current working directory.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currDir, cdErr := manager.GetFileInfoProvider().GetCurrentDirectory()
		if cdErr != nil {
			return cdErr
		}
		dcFile, dcErr := manager.LocateFileInDirectory(currDir) // assuming file is in current directory
		projectName := ""                                       // project name will be decided on later

		switch len(args) {
		case 0:
			// directory name is used as project name
			projectName = manager.GetFileInfoProvider().GetDirectoryName(currDir)
			if dcErr != nil {
				return dcErr
			}
			break
		case 1:
			if dcErr != nil {
				return dcErr
			} else {
				projectName = args[0] // project name was provided as first argument
			}
			break

		case 2:
			dcFile = manager.GetFileInfoProvider().Expand(args[1])
			if manager.GetFileInfoProvider().IsDir(dcFile) { // if second argument is a directory
				var fileErr error
				dcFile, fileErr = manager.LocateFileInDirectory(dcFile)
				if fileErr != nil {
					return fileErr
				}
			} else if !manager.GetFileInfoProvider().IsFile(dcFile) { // if second argument is not a dir nor a file
				return errors.New("provided file does not exist")
			}
			// second argument was a dir containing a file or was a file itself
			projectName = args[0] // project name was provided as first argument
			break

		default:
			return errors.New("invalid arguments count")
		}

		if err := manager.GetConfigFile().AddDockerComposeFile(dcFile, projectName); err != nil {
			return err
		}

		_, _ = fmt.Fprintf(mainWriter, "File '%s' added to project '%s'", dcFile, projectName)

		return nil
	},
}

func init() {
	RootCommand.AddCommand(addCommand)
}
