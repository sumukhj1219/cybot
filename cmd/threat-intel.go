package cmd

import (
	"cybot/services"
	"fmt"
	"os"
	"sync"

	"github.com/olekukonko/tablewriter" // Import tablewriter
	"github.com/spf13/cobra"
)

var threatCmd = &cobra.Command{
	Use:   "threat-intel",
	Short: "Aggregate and correlate threat intelligence",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: threat-intel <hash|ip|domain>")
			os.Exit(1)
		}

		input := args[0]

		var wg sync.WaitGroup
		results := make(chan interface{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			vtResponse, err := services.ThreatIntel(input) // Get the structured response

			// Send the response to the channel (if it's not nil)
			if err != nil {
				fmt.Println(err)
			}
			if vtResponse != nil {
				results <- vtResponse
			}

		}()

		go func() {
			wg.Wait()
			close(results)
		}()

		// Create the table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Source", "Maliciousness", "Reputation", "Harmless", "Malicious", "Suspicious", "Scan Date", "Last Analysis Date"}) // Set headers

		for result := range results {
			vtResponse, ok := result.(*services.VirusTotalResponse)
			if ok && vtResponse != nil {
				// Format dates if necessary
				lastAnalysisDate := formatLastAnalysisDate(vtResponse.Data.Attributes.LastAnalysisDate)

				// Append data to the table
				table.Append([]string{
					"VirusTotal",
					fmt.Sprintf("%d", vtResponse.Data.Attributes.Maliciousness),
					fmt.Sprintf("%d", vtResponse.Data.Attributes.Reputation),
					fmt.Sprintf("%d", vtResponse.Data.Attributes.LastAnalysisStats.Harmless),
					fmt.Sprintf("%d", vtResponse.Data.Attributes.LastAnalysisStats.Malicious),
					fmt.Sprintf("%d", vtResponse.Data.Attributes.LastAnalysisStats.Suspicious),
					vtResponse.Data.Attributes.ScanDate,
					lastAnalysisDate,
				})

				// Contextual insights
				context := "File is not associated with known malicious IPs or domains. It's a new or less-known threat."
				if len(vtResponse.Data.Attributes.ScanEngines) < 5 {
					context = "This file was detected by a limited number of AV engines, indicating it could be a new or targeted threat."
				}
				if len(vtResponse.Data.Attributes.ScanEngines) > 5 && context != "This file was detected by a limited number of AV engines, indicating it could be a new or targeted threat." {
					context = "The file has been detected by multiple AV engines, increasing the likelihood of it being a known threat."
				}

				// Print context
				fmt.Println("Context:", context)

				// Actionable recommendations
				recommendation := "Further analysis required. Consider quarantining the file for deeper inspection."
				if len(vtResponse.Data.Attributes.ScanEngines) > 5 && context != "The file has been detected by multiple AV engines, increasing the likelihood of it being a known threat." {
					recommendation = "File detected as safe. However, continue monitoring for potential future detections."
				}

				// Print recommendations
				fmt.Println("Recommendation:", recommendation)

				// Print scan engine results
				for _, engine := range vtResponse.Data.Attributes.ScanEngines {
					fmt.Printf("Engine: %s, Result: %s\n", engine.EngineName, engine.Result)
				}
			} else {
				fmt.Println("Error: Invalid or nil response received.")
			}
		}

		table.Render() // Render the table
	},
}

func init() {
	rootCmd.AddCommand(threatCmd)
}

// Helper function to format LastAnalysisDate
func formatLastAnalysisDate(date interface{}) string {
	switch v := date.(type) {
	case string:
		return v
	case float64: // Assuming it's a Unix timestamp in float64 format
		// Convert to int64 and format
		return fmt.Sprintf("%d", int64(v))
	default:
		return "Unknown"
	}
}
