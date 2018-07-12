package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pagerinc/kongfig/api"
	"github.com/spf13/cobra"
)

var file string

func init() {
	const (
		defaultConfig = "config.json"
		usage         = "Filename that contains the configuration to apply"
	)

	kongfig.Flags().StringVarP(&file, "file", "f", defaultConfig, usage)
	kongfig.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a Kong instance",
	Long:  `Use apply to restore your settings into an existing Kong instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		c := api.Config{}
		json.NewDecoder(file).Decode(&c)

		for _, s := range c.Services {
			c.UpdateService(s)
			c.CreateRoutes(s)
			c.GetRoutes(s)
		}
	},
}
