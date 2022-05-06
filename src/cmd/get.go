package cmd

import (
	"fmt"
	"os"

	"github.com/aiden-deloryn/hoist/src/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [address]",
	Short: "Download a file being shared from another computer on a local area network",
	Long:  `Download a file being shared from another computer on a local area network.`,
	Run:   runGetCmd,
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

func runGetCmd(cmd *cobra.Command, args []string) {
	if err := client.GetFileFromServer(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
