package cmds

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pengsrc/docker-tools/constants"
	"github.com/pengsrc/docker-tools/remote"
	"github.com/pengsrc/go-shared/check"
)

// remoteImportCMD represents the remote-import sub-command.
var remoteImportCMD = &cobra.Command{
	Use:   "remote-import <image>[:tag]",
	Short: "Import docker image from one registry to another",
	Long:  "Import docker image from one registry to another",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := remote.NewSSHSession(
			importImageFlagHost, importImageFlagPort, importImageFlagUser,
		)
		check.ErrorForExit(constants.Name, err)

		from := fmt.Sprintf("%s/%s", importImageFlagFrom, args[0])
		to := fmt.Sprintf("%s/%s/%s", importImageFlagTo, importImageFlagFrom, args[0])

		commands := []string{
			fmt.Sprintf("docker pull %s", from),
			fmt.Sprintf("docker tag %s %s", from, to),
			fmt.Sprintf("docker push %s", to),
			fmt.Sprintf("docker rmi %s %s", from, to),
		}
		for _, command := range commands {
			fmt.Printf("Executing: %s\n", command)
		}
		err = session.Run(strings.Join(commands, " && "))
		check.ErrorForExit(constants.Name, err)
	},
	Args: cobra.ExactArgs(1),
}

var (
	importImageFlagHost string
	importImageFlagPort int
	importImageFlagUser string

	importImageFlagFrom string
	importImageFlagTo   string

	importImageFlagHelp bool
)

func initImportImageCMD() {
	remoteImportCMD.Flags().StringVarP(
		&importImageFlagHost, "host", "h", "127.0.0.1", "SSH host to run import procedures",
	)
	remoteImportCMD.Flags().IntVarP(
		&importImageFlagPort, "port", "p", 22, "SSH port to connect",
	)
	remoteImportCMD.Flags().StringVarP(
		&importImageFlagUser, "user", "u", "root", "SSH username",
	)
	remoteImportCMD.Flags().StringVarP(
		&importImageFlagFrom, "from", "f", "docker.io", "Registry to import image from",
	)
	remoteImportCMD.Flags().StringVarP(
		&importImageFlagTo, "to", "t", "registry.example.com", "Registry to export image to",
	)
	remoteImportCMD.Flags().BoolVarP(
		&importImageFlagHelp, "help", "", false, "Show help",
	)

	remoteImportCMD.MarkFlagRequired("host")
	remoteImportCMD.MarkFlagRequired("port")
	remoteImportCMD.MarkFlagRequired("user")
	remoteImportCMD.MarkFlagRequired("from")
	remoteImportCMD.MarkFlagRequired("to")

	if !remoteImportCMD.Flag("host").Changed {
		if host := viper.GetString("builder.host"); host != "" {
			importImageFlagHost = host
		}
	}
	if !remoteImportCMD.Flag("port").Changed {
		if port := viper.GetInt("builder.port"); port != 0 {
			importImageFlagPort = port
		}
	}
	if !remoteImportCMD.Flag("user").Changed {
		if user := viper.GetString("builder.user"); user != "" {
			importImageFlagUser = user
		}
	}

	if !remoteImportCMD.Flag("to").Changed {
		if to := viper.GetString("registry"); to != "" {
			importImageFlagTo = to
		}
	}
}
