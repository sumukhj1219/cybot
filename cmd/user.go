package cmd

import (
	"cybot/services"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user-config",
	Short: "A coomand which displays user details.",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := services.UserDetails()
		if err != nil {
			fmt.Println("Error in displaying user detials")
		}
		fmt.Println("\n🛡️  System & Network Info 🛡️")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Property", "Value"})
		table.SetBorder(true)
		table.SetRowLine(true)

		data := [][]string{
			{"🔹 Local IP", user.LocalIP},
			{"🌍 Public IP", user.PublicIP},
			{"🖥️ OS", user.OS},
			{"💻 Architecture", user.Arch},
			{"🏠 Hostname", user.Hostname},
		}

		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
