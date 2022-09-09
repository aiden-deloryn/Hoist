package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/aiden-deloryn/hoist/src/server"
	"github.com/aiden-deloryn/hoist/src/util"
	"github.com/aiden-deloryn/hoist/src/values"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [filename]",
	Short: "Send a file over a local area network",
	Long:  `Send a file over a local area network.`,
	RunE:  runSendCmd,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sendCmd.Flags().BoolP("keep-alive", "k", false, "Keep the connection open for multiple transfers")
}

func runSendCmd(cmd *cobra.Command, args []string) error {
	keepAlive, _ := cmd.Flags().GetBool("keep-alive")
	ip, err := util.GetLocalIPAddress()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get local IP address: %s\n", err.Error())
	}

	filename := filepath.FromSlash(strings.TrimSuffix(args[0], string(filepath.Separator)))

	fmt.Print("Enter a password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return fmt.Errorf("failed to set password: %s", err)
	}

	if len(password) > values.MAX_PASSWORD_LENGTH {
		return fmt.Errorf("Password length must be %d characters or less", values.MAX_PASSWORD_LENGTH)
	}

	err = server.StartServer(fmt.Sprintf("%s:8080", ip), filename, string(password), keepAlive)

	if err != nil {
		return fmt.Errorf("server error: %s", err)
	}

	return nil
}
