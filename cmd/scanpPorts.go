package cmd

import (
	"cybot/services"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var scanPortsCmd = &cobra.Command{
	Use:   "scan",
	Short: "Used to scan ports --scan <IP> <startPort> <endPort> <protocol>",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]

		startPort, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("❌ Invalid startPort:", err)
			return
		}

		endPort, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("❌ Invalid endPort:", err)
			return
		}

		protocol := args[3]

		services.PortScanner(ip, startPort, endPort, protocol)
	},
}

func init() {
	rootCmd.AddCommand(scanPortsCmd)
}
