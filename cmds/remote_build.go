package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pengsrc/docker-tools/build"
	"github.com/pengsrc/docker-tools/constants"
	"github.com/pengsrc/docker-tools/remote"
	"github.com/pengsrc/go-shared/check"
)

// remoteBuildCMD represents the remote-build sub-command.
var remoteBuildCMD = &cobra.Command{
	Use:   "remote-build <image>[:tag]",
	Short: "Build docker image and push to specified registry",
	Long:  "Build docker image and push to specified registry",
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, err := filepath.Abs(buildImageFlagDirectory)
		check.ErrorForExit(constants.Name, err)
		err = check.Dir(sourceDir)
		check.ErrorForExit(constants.Name, err)

		// Prepare image and tag.
		image, tag := build.ParseImageTag(args[0])
		if image == "" {
			check.ErrorForExit(constants.Name, fmt.Errorf("image not valid: %s", args[0]))
		}
		if buildImageFlagGitArchive {
			tag, err = build.ImageTagForGitRepo(sourceDir, tag)
			check.ErrorForExit(constants.Name, err)
		}
		if tag == "" {
			tag = "latest"
		}

		// Pack source code.
		localPath, err := build.CreatePackage(
			sourceDir, buildImageFlagInclude, buildImageFlagExclude,
			buildImageFlagGitArchive, tag,
		)
		defer os.Remove(localPath)
		if err != nil {
			os.Remove(localPath)
			check.ErrorForExit(constants.Name, err)
		}

		// Create remote build directory.
		session, err := remote.NewSSHSession(
			buildImageFlagHost, buildImageFlagPort, buildImageFlagUser,
		)
		if err != nil {
			os.Remove(localPath)
			check.ErrorForExit(constants.Name, err)
		}

		command := fmt.Sprintf("mkdir -p %s", constants.RemoteBuildDir)
		fmt.Printf("Executing: %s\n", command)
		err = session.Run(command)
		if err != nil {
			os.Remove(localPath)
			check.ErrorForExit(constants.Name, err)
		}

		// Upload package.
		packageExt := ".tar.gz"
		remotePath := filepath.Join(constants.RemoteBuildDir, fmt.Sprintf(
			"%s%s", filepath.Base(localPath), packageExt,
		))
		err = remote.Upload(
			buildImageFlagHost, buildImageFlagPort, buildImageFlagUser, localPath, remotePath,
		)
		if err != nil {
			os.Remove(localPath)
			check.ErrorForExit(constants.Name, err)
		}

		// Build and clean.
		session, err = remote.NewSSHSession(
			buildImageFlagHost, buildImageFlagPort, buildImageFlagUser,
		)
		if err != nil {
			os.Remove(localPath)
			check.ErrorForExit(constants.Name, err)
		}

		workDir := strings.TrimRight(remotePath, packageExt)
		commands := []string{
			fmt.Sprintf("mkdir -p %s", workDir),
			fmt.Sprintf("cd %s", workDir),
			fmt.Sprintf("tar -xf %s", remotePath),
		}
		if buildImageFlagBefore != "" {
			commands = append(commands, buildImageFlagBefore)
		}
		commands = append(commands, fmt.Sprintf(
			"docker build -f %s -t %s/%s:%s .",
			buildImageFlagDockerfile, buildImageFlagRegistry, image, tag,
		))
		commands = append(commands, fmt.Sprintf(
			"docker push %s/%s:%s", buildImageFlagRegistry, image, tag),
		)
		if buildImageFlagAfter != "" {
			commands = append(commands, buildImageFlagAfter)
		}
		commands = append(commands, fmt.Sprintf("rm -f %s", remotePath))
		commands = append(commands, fmt.Sprintf("rm -rf %s", workDir))

		for _, command := range commands {
			fmt.Printf("Executing: %s\n", command)
		}
		err = session.Run(strings.Join(commands, " && "))
		if err != nil {
			os.Remove(localPath)
			check.ErrorForExit(constants.Name, err)
		}
	},
	Args: cobra.ExactArgs(1),
}

var (
	buildImageFlagHost string
	buildImageFlagPort int
	buildImageFlagUser string

	buildImageFlagDockerfile string
	buildImageFlagRegistry   string

	buildImageFlagBefore string
	buildImageFlagAfter  string

	buildImageFlagDirectory  string
	buildImageFlagInclude    []string
	buildImageFlagExclude    []string
	buildImageFlagGitArchive bool

	buildImageFlagHelp bool
)

func initBuildImageCMD() {
	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagHost, "host", "h", "127.0.0.1", "SSH host to run import procedures",
	)
	remoteBuildCMD.Flags().IntVarP(
		&buildImageFlagPort, "port", "p", 22, "SSH port to connect",
	)
	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagUser, "user", "u", "root", "SSH username",
	)

	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagDockerfile, "dockerfile", "f", "Dockerfile", "Dockerfile to use in build",
	)
	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagRegistry, "registry", "r", "registry.example.com", "Registry to push image",
	)

	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagBefore, "before", "b", "", "Command to execute before the build",
	)
	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagAfter, "after", "a", "", "Command to execute after the build",
	)

	remoteBuildCMD.Flags().StringVarP(
		&buildImageFlagDirectory, "directory", "d", "", "Source directory to use",
	)
	remoteBuildCMD.Flags().StringArrayVarP(
		&buildImageFlagInclude, "include", "i", []string{}, "Files to include in package",
	)
	remoteBuildCMD.Flags().StringArrayVarP(
		&buildImageFlagExclude, "exclude", "e", []string{}, "Files to exclude in package",
	)
	remoteBuildCMD.Flags().BoolVarP(
		&buildImageFlagGitArchive, "git-archive", "g", false, "Use git archive to pack files",
	)

	remoteBuildCMD.Flags().BoolVarP(
		&buildImageFlagHelp, "help", "", false, "Show help",
	)

	remoteBuildCMD.MarkFlagRequired("host")
	remoteBuildCMD.MarkFlagRequired("port")
	remoteBuildCMD.MarkFlagRequired("user")

	remoteBuildCMD.MarkFlagRequired("dockerfile")
	remoteBuildCMD.MarkFlagRequired("registry")

	if !remoteBuildCMD.Flag("host").Changed {
		if host := viper.GetString("builder.host"); host != "" {
			buildImageFlagHost = host
		}
	}
	if !remoteBuildCMD.Flag("port").Changed {
		if port := viper.GetInt("builder.port"); port != 0 {
			buildImageFlagPort = port
		}
	}
	if !remoteBuildCMD.Flag("user").Changed {
		if user := viper.GetString("builder.user"); user != "" {
			buildImageFlagUser = user
		}
	}

	if !remoteBuildCMD.Flag("registry").Changed {
		if registry := viper.GetString("registry"); registry != "" {
			buildImageFlagRegistry = registry
		}
	}
}
