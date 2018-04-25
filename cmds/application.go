package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pengsrc/docker-tools/constants"
	"github.com/pengsrc/go-shared/check"
	"github.com/pengsrc/go-shared/utils"
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	check.ErrorForExit(constants.Name, rootCMD.Execute())
}

// rootCMD represents the base command when called without any subcommands.
var rootCMD = &cobra.Command{
	Use:   constants.Name,
	Short: "Handy tools for Docker",
	Long:  "Handy tools for Docker",
	Run: func(cmd *cobra.Command, args []string) {
		if rootFlagVersion {
			showVersion()
			return
		}
		cmd.Help()
	},
}

var (
	rootFlagConfig  string
	rootFlagVersion bool
	rootFlagHelp    bool
)

func init() {
	rootCMD.SilenceErrors = true
	rootCMD.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	rootCMD.PersistentFlags().StringVarP(
		&rootFlagConfig, "config", "c", "",
		fmt.Sprintf("Configuration file (default is ${HOME}/%s)", constants.DefaultConfigFileName),
	)
	rootCMD.Flags().BoolVarP(
		&rootFlagVersion, "version", "v", false,
		"Show version",
	)
	rootCMD.Flags().BoolVarP(
		&rootFlagHelp, "help", "", false,
		"Show help",
	)

	rootCMD.AddCommand(remoteImportCMD)
	rootCMD.AddCommand(remoteBuildCMD)

	initConfig()
	initImportImageCMD()
	initBuildImageCMD()
}

func initConfig() {
	viper.SupportedExts = []string{"yaml"}

	if rootFlagConfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(rootFlagConfig)
	} else {
		// Don't read default config file if not exists.
		_, err := os.Stat(filepath.Join(utils.GetHome(), constants.DefaultConfigFileName))
		if err != nil {
			return
		}

		// Search config in home directory with name ".docker-tools" (without extension).
		viper.AddConfigPath(filepath.Join(utils.GetHome()))
		viper.SetConfigName(".docker-tools")
	}
	err := viper.ReadInConfig()
	check.ErrorForExit(constants.Name, err)
}

func showVersion() {
	fmt.Printf("%s version %s\n", constants.Name, constants.Version)
}
