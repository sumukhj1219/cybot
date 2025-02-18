package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LastAnalysisStats struct {
	Harmless   int `json:"harmless"`
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
}

type ScanEngine struct {
	EngineName string `json:"engine_name"`
	Category   string `json:"category"`
	Result     string `json:"result"`
}

type FileContext struct {
	IsTargeted        bool
	IsNewThreat       bool
	AssociatedIPs     []string
	AssociatedDomains []string
}

type VirusTotalResponse struct {
	Data struct {
		Attributes struct {
			Maliciousness     int               `json:"maliciousness"`
			Reputation        int               `json:"reputation"`
			LastAnalysisStats LastAnalysisStats `json:"last_analysis_stats"`
			LastAnalysisDate  interface{}       `json:"last_analysis_date"` // Change to interface{}
			ScanDate          string            `json:"scan_date"`
			ScanEngines       []ScanEngine      `json:"scan_engines"`
		} `json:"attributes"`
	} `json:"data"`
}

func formatLastAnalysisDate(date interface{}) string {
	switch v := date.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%d", int64(v))
	default:
		return "Unknown"
	}
}

func ThreatIntel(hash string) (*VirusTotalResponse, error) {
	// Retrieve the API key from the environment variable
	apiKey := ""
	if apiKey == "" {
		return nil, fmt.Errorf("VIRUSTOTAL_API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", hash)

	// Create a new request with the API key in the headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Add("x-apikey", apiKey)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request to VirusTotal: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VirusTotal API returned non-200 status: %d", resp.StatusCode)
	}

	// Decode the JSON response
	var vtResponse VirusTotalResponse
	if err := json.NewDecoder(resp.Body).Decode(&vtResponse); err != nil {
		return nil, fmt.Errorf("Error decoding JSON response: %v", err)
	}

	// Return the response for further processing
	return &vtResponse, nil
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
