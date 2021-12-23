package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var addCommand = &cobra.Command{
	Use:   "add [project name] [docker-compose-file]",
	Short: "Adds docker-compose file to available files list",
	Long: `Adds docker-compose file to available files list

If no project name is provided current directory name is used.
If no file is provided it look for one in current working directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		currDir, cdErr := manager.GetFileInfoProvider().GetCurrentDirectory()
		if cdErr != nil {
			log.Fatal(cdErr)
		}
		dcFile, dcErr := manager.LocateFileInDirectory(currDir)
		projectName := ""

		switch len(args) {
		case 0:
			if dcErr != nil {
				log.Fatal(dcErr)
			}
			break
		case 1:
			if dcErr != nil {
				log.Fatal(dcErr)
			} else {
				projectName = args[0]
			}
			break

		case 2:
			dcFile = manager.GetFileInfoProvider().Expand(args[1])
			if manager.GetFileInfoProvider().IsDir(dcFile) {
				var fileErr error
				dcFile, fileErr = manager.LocateFileInDirectory(dcFile)
				if fileErr != nil {
					log.Fatal(fileErr)
				}
			} else if !manager.GetFileInfoProvider().IsFile(dcFile) {
				log.Fatal("Provided file does not exist")
			}
			projectName = args[0]
			break

		default:
			fmt.Println("Invalid arguments count")
			os.Exit(2)

		}

		if err := manager.GetConfigFile().AddDockerComposeFile(dcFile, projectName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCommand.AddCommand(addCommand)
}
