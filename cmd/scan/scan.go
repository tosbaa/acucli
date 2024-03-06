/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type postBody struct {
	TargetID    string   `json:"target_id"`
	ProfileID   string   `json:"profile_id"`
	Schedule    Schedule `json:"schedule"`
	Incremental bool     `json:"incremental"`
}

type Schedule struct {
	Disable       bool    `json:"disable"`
	TimeSensitive bool    `json:"time_sensitive"`
	StartDate     *string `json:"start_date"`
}

var scanProfileId string

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Command to start the scan",
	Long: `Command to start the scan, takes target ids from stdin and id of scan profile in flag, example:

cat target_ids.txt | acucli scan --scanProfileID=47973ea9-018b-4294-9903-bb1cf3b1e886`,
	Run: func(cmd *cobra.Command, args []string) {
		targets := filehelper.ReadStdin()
		scanProfileID, _ := cmd.Flags().GetString("scanProfileID")
		var responseCode int
		for _, target := range targets {
			responseCode = startScan(target, scanProfileID)
			if responseCode == 200 {
				fmt.Printf("%sScan started: %s%s\n", chalk.Green, target, chalk.Reset)
			} else {
				fmt.Printf("%sError occured while starting scan%s\n", chalk.Red, chalk.Reset)
			}
		}
	},
}

func startScan(targetID string, scanProfileID string) int {
	postBody := postBody{ProfileID: scanProfileID, Incremental: false, Schedule: Schedule{Disable: false, TimeSensitive: false, StartDate: nil}}

	postBody.TargetID = targetID
	requestJson, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/scans"), bytes.NewBuffer(requestJson))
	if err != nil {
		fmt.Println("Error creating request:", err)
		fmt.Println(err)
		return 500
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		fmt.Println(err)
		return 500
	}

	if resp.StatusCode != 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
			fmt.Println(err)
		}

		// Convert the body to type string
		fmt.Println(string(body))
		return resp.StatusCode
	} else {
		return 200
	}

}

func init() {
	ScanCmd.Flags().StringVarP(&scanProfileId, "scanProfileID", "", "", "scanProfile ID")
	ScanCmd.MarkFlagRequired("scanProfileID")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
