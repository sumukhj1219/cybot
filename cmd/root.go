package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func printBanner() {
	banner := `
 ██████╗██╗   ██╗██████╗  ██████╗ ████████╗
██╔════╝██║   ██║██╔══██╗██╔═══██╗╚══██╔══╝
██║     ██║   ██║██████╔╝██║   ██║   ██║   
██║     ██║   ██║██╔══██╗██║   ██║   ██║   
╚██████╗╚██████╔╝██║  ██║╚██████╔╝   ██║   
 ╚═════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝    ╚═╝   
`
	fmt.Println(banner)
	fmt.Println("🚀 Welcome to CyBot - Cybersecurity CLI Tool! 🚀\n")
}

func printCommandTable() {
	data := [][]string{
		{"user", "User configurtaion", "cybot user"},
		{"scan", "Scan ports", "cybot scan <IP> <startPort> <endPort> <protocol>"},
		{"net-scan", "Scan devices locally", "cybot net-scan"},
		{"scan-p", "Scan all the running processes.", "cybot scan-p"},
		// {"firewall", "Check firewall status", "cybot --firewall"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Command", "Description", "Usage"})
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

var rootCmd = &cobra.Command{
	Use:   "cybot",
	Short: "CyBot - A Powerful Cybersecurity CLI Tool",
	Long: `CyBot is a command-line cybersecurity tool designed for 
penetration testing, system security analysis, and network scanning.`,
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
		printCommandTable()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
