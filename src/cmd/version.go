package cmd

import (
	"fmt"

	"github.com/aiden-deloryn/hoist/src/values"
	"github.com/spf13/cobra"
)

// versionCmd represents the get command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the application version",
	Long:  `Show the application version.`,
	RunE:  runVersionCmd,
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s version %s\n", values.APP_NAME, values.APP_VERSION)

	return nil
}
