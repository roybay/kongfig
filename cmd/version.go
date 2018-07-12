package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	VERSION string = "Kongfig Settings Manager v0.0.1-alpha -- HEAD"
)

func init() {
	kongfig.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Kongfig",
	Long:  `All software has versions. This is Kongfig's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
