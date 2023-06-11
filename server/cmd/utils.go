package cmd

import "github.com/spf13/cobra"

func getConfigDSN(command *cobra.Command, env map[string]string) string {
	configDSN, _ := command.Flags().GetString("config")

	if configDSN == "" {
		configDSN = env["PARTSY_CONFIG"]
	}

	if configDSN == "" {
		configDSN = "config.json"
	}

	return configDSN
}
