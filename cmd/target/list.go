/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tosbaa/acucli/helpers/httpclient"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TargetList struct {
	Targets []struct {
		Address                  string `json:"address"`
		Agents                   any    `json:"agents"`
		ContinuousMode           bool   `json:"continuous_mode"`
		Criticality              int    `json:"criticality"`
		DefaultScanningProfileID string `json:"default_scanning_profile_id"`
		DeletedAt                any    `json:"deleted_at"`
		Description              string `json:"description"`
		Fqdn                     string `json:"fqdn"`
		FqdnHash                 string `json:"fqdn_hash"`
		FqdnStatus               string `json:"fqdn_status"`
		FqdnTmHash               string `json:"fqdn_tm_hash"`
		IssueTrackerID           any    `json:"issue_tracker_id"`
		LastScanDate             string `json:"last_scan_date"`
		LastScanID               string `json:"last_scan_id"`
		LastScanSessionID        string `json:"last_scan_session_id"`
		LastScanSessionStatus    string `json:"last_scan_session_status"`
		ManualIntervention       bool   `json:"manual_intervention"`
		SeverityCounts           struct {
			Critical int `json:"critical"`
			High     int `json:"high"`
			Info     int `json:"info"`
			Low      int `json:"low"`
			Medium   int `json:"medium"`
		} `json:"severity_counts"`
		TargetID     string `json:"target_id"`
		Threat       int    `json:"threat"`
		Type         any    `json:"type"`
		Verification any    `json:"verification"`
	} `json:"targets"`
	Pagination struct {
		Count      int    `json:"count"`
		CursorHash string `json:"cursor_hash"`
		Cursors    []any  `json:"cursors"`
		Sort       any    `json:"sort"`
	} `json:"pagination"`
}

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the targets",
	Long:  `Lists all the targets with their name and their corresponding id to use it for other commands`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create an HTTP GET request using the custom client
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets"), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Perform the request using the custom client
		resp, err := httpclient.MyHTTPClient.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var targetList TargetList
		err = json.Unmarshal(body, &targetList)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		for _, target := range targetList.Targets {
			fmt.Printf("%s\t%s\n", target.Address, target.TargetID)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
