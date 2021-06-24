package cmd

import (
	"dimitrios_liappis/instancepeek/formatters/terminal"
	"dimitrios_liappis/instancepeek/model"
	"dimitrios_liappis/instancepeek/providers/aws"
	"dimitrios_liappis/instancepeek/providers/dummy"
	"os"

	"github.com/spf13/cobra"
)

// Provider ...
var Provider string

var (
	rootCmd = &cobra.Command{
		Use:   "instancepeek",
		Short: "List hardware details of cloud instance types",
		Long:  `TODO`,
		Args:  cobra.MinimumNArgs(1),
		Run:   entrypoint,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&Provider, "provider", "", "Specifies the provider (aws, gcp, azure, dummy)")
	rootCmd.MarkPersistentFlagRequired("provider")
}

func entrypoint(cmd *cobra.Command, args []string) {
	instanceTypes := args
	processCommand(instanceTypes)
}

func processCommand(instanceTypes []string) error {
	var instanceInfos []model.InstanceInfo

	switch Provider {
	case "dummy":
		instanceInfos, _ = dummy.Convert(instanceTypes)
	case "aws":
		instanceInfos, _ = aws.Convert(instanceTypes)
	}

	terminal.Format(instanceInfos, os.Stdout)
	return nil
}
