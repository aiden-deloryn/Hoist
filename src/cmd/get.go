package cmd

import (
	"fmt"
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
}

func runGetCmd(cmd *cobra.Command, args []string) error {
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return fmt.Errorf("failed to read password: %s", err)
	}

	if err := client.GetFileFromServer(args[0], string(password)); err != nil {
		return err
	}

	return nil
}
