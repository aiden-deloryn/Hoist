package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/aiden-deloryn/hoist/src/client"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [address]",
	Short: "Download a file being shared from another computer on a local area network",
	Long:  `Download a file being shared from another computer on a local area network.`,
	RunE:  runGetCmd,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().Bool("no-password", false, "Do not prompt for a password (password will be blank)")
	getCmd.Flags().StringP("password", "p", "", "Provide the password to use for authentication")
	getCmd.Flags().StringP("output", "o", "", "Set a custom output directory")
}

func runGetCmd(cmd *cobra.Command, args []string) error {
	skipPassword, _ := cmd.Flags().GetBool("no-password")
	password, _ := cmd.Flags().GetString("password")
	outputDirectory, _ := cmd.Flags().GetString("output")

	// Bash doesn't expand "~" if the path is in single or double quotes
	if strings.HasPrefix(outputDirectory, "~") {
		user, err := user.Current()

		if err != nil {
			return fmt.Errorf("failed to expand home directory (~): %s", err)
		}

		outputDirectory = filepath.Join(user.HomeDir, outputDirectory[1:])
	}

	if !skipPassword && password == "" {
		fmt.Print("Enter password: ")
		passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()

		if err != nil {
			return fmt.Errorf("failed to read password: %s", err)
		}

		password = string(passwordBytes)
	}

	if err := client.GetFileFromServer(args[0], string(password), outputDirectory); err != nil {
		return err
	}

	return nil
}
