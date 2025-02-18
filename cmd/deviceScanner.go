package cmd

import (
	"cybot/services"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "net-scan",
	Short: "Scan for active devices on the local network cybot net-scan",
	Run: func(cmd *cobra.Command, args []string) {
		services.ScanNetworkARP()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
