package cmd

import (
	"cybot/services"
	"fmt"

	"github.com/spf13/cobra"
)

var sniffCmd = &cobra.Command{
	Use:   "scan-a",
	Short: "Start packet sniffing",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting packet sniffer...")
		services.ScanAttacks()
	},
}

func init() {
	rootCmd.AddCommand(sniffCmd)
}
