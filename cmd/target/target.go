/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type responseBody struct {
	Address                  string `json:"address"`
	Agents                   []any  `json:"agents"`
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
}

var id string

func (r responseBody) String() string {
	var builder strings.Builder

	builder.WriteString("Response Body:\n")
	builder.WriteString("+-----------------------------+-----------------------------------------------------+\n")
	builder.WriteString("| Field                       | Value                                               |\n")
	builder.WriteString("+-----------------------------+-----------------------------------------------------+\n")

	addRow(&builder, "Address", r.Address)
	addRow(&builder, "Agents", fmt.Sprintf("%v", r.Agents))
	addRow(&builder, "ContinuousMode", fmt.Sprintf("%t", r.ContinuousMode))
	addRow(&builder, "Criticality", fmt.Sprintf("%d", r.Criticality))
	addRow(&builder, "DefaultScanningProfileID", r.DefaultScanningProfileID)
	// ... Add rows for other fields

	builder.WriteString("+-----------------------------+-----------------------------------------------------+\n")

	return builder.String()
}

func addRow(builder *strings.Builder, field, value string) {
	builder.WriteString(fmt.Sprintf("| %-28s | %-50s |\n", field, value))
}

// targetCmd represents the target command
var TargetCmd = &cobra.Command{
	Use:   "target",
	Short: "Endpoint for target operations",
	Long:  `You can do operations with subcommands to the targets or you can retrieve target information by giving id flag to the target command`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ = cmd.Flags().GetString("id")
		responseCode, respBody := GetTargetRequest(id)
		if responseCode == 200 {
			fmt.Println(respBody)
		} else {
			fmt.Fprintf(os.Stderr, "%sTarget not found%s\n", chalk.Red, chalk.Reset)
		}

	},
}

func GetTargetRequest(id string) (int, responseBody) {
	var respBody responseBody
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", httpclient.BASE_URL, "/targets/", id), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 404, respBody
	}

	// Perform the request using the custom client
	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 404, respBody
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 404, respBody
	}
	if resp.StatusCode == 404 {
		return 404, respBody
	} else {
		return 200, respBody
	}
}

func init() {
	TargetCmd.Flags().StringVarP(&id, "id", "", "", "Target ID")
	TargetCmd.MarkFlagRequired("id")

	TargetCmd.AddCommand(ListCmd)
	TargetCmd.AddCommand(AddCmd)
	TargetCmd.AddCommand(RemoveCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
