package cmd

import (
	"github.com/aiden-deloryn/hoist/src/server"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [filename]",
	Short: "Send a file over a local area network",
	Long:  `Send a file over a local area network.`,
	Run:   runSendCmd,
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

func runSendCmd(cmd *cobra.Command, args []string) {
	keepAlive, _ := cmd.Flags().GetBool("keep-alive")
	server.StartServer("localhost:8080", args[0], keepAlive)
}
