package cmd

import (
	"github.com/spf13/cobra"
)

// Provider ...
var Provider string

var (
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "instancepeek",
		Short: "Easily list hardware details of cloud instance types",
		Long:  `TODO`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&Provider, "provider", "", "Specifies the provider (aws, gcp, azure, dummy)")
	rootCmd.MarkFlagRequired("provider")
}
