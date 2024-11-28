/*
Copyright © 2024 Tom Hetherington thetherington@evertz.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string
	BuiltBy string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", Version)
		fmt.Println("Commit:", Commit)
		fmt.Println("Date:", Date)
		fmt.Println("BuiltBy:", BuiltBy)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
