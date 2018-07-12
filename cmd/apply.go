package cmd

import (
	"os"

	"github.com/pagerinc/kongfig/api"
	"github.com/spf13/cobra"
)

var (
	fileVar   string
	dryRunVar bool
)

func init() {
	const (
		defaultConfig = "config.json"
		configUsage   = "Filename that contains the configuration to apply"
		defaultDryRun = false
		dryRunUsage   = "simulate an install"
	)

	applyCmd.Flags().StringVarP(&fileVar, "file", "f", defaultConfig, configUsage)
	applyCmd.Flags().BoolVar(&dryRunVar, "dry-run", defaultDryRun, dryRunUsage)
	kongfig.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a Kong instance",
	Long:  `Use apply to restore your settings into an existing Kong instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(fileVar)
		if err != nil {
			return err
		}
		defer file.Close()

		c := api.NewConfig(file)

		for _, s := range c.Services {
			c.UpdateService(s)
			c.CreateRoutes(s)
			c.GetRoutes(s)
		}

		return nil
	},
}
