/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scanProfile

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type ScanProfile struct {
	Checks    []string `json:"checks"`
	Custom    bool     `json:"custom"`
	Name      string   `json:"name"`
	ProfileID string   `json:"profile_id"`
	SortOrder int      `json:"sort_order"`
}

func (sp ScanProfile) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Name: %s\n", sp.Name))
	sb.WriteString(fmt.Sprintf("ProfileID: %s\n", sp.ProfileID))
	sb.WriteString(fmt.Sprintf("Custom: %v\n", sp.Custom))
	sb.WriteString(fmt.Sprintf("SortOrder: %d\n", sp.SortOrder))
	sb.WriteString("Checks:\n")
	for _, check := range sp.Checks {
		sb.WriteString(fmt.Sprintf("  - %s\n", check))
	}

	return sb.String()
}

var id string

// scanProfileCmd represents the scanProfile command
var ScanProfileCmd = &cobra.Command{
	Use:   "scanProfile",
	Short: "Manage Scan Profiles",
	Long: `Without flags with just id, it gives the summary of the Scan Profile, can also export with --export flag. Example:
	
	acucli scanProfile --id=9ef9d9ee-510e-47da-8549-43ea826d1cdc : Gives the summary of the Scan Profile
	acucli scanProfile --id=9ef9d9ee-510e-47da-8549-43ea826d1cdc --export : Exports the Scan Profile JSON
	acucli scanProfile --id=03aa8950-a289-48f3-a0f9-8416b8c5a8d5 --export --output=/home/user : Exports the Scan Profile to defined output folder

	`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ = cmd.Flags().GetString("id")
		if cmd.Flags().Changed("export") {
			output, _ := cmd.Flags().GetString("output")
			exportScanProfile(id, output)
		} else {
			responseCode, respBody := GetScanProfileRequest(id)
			if responseCode == 200 {
				fmt.Println(respBody)
			} else {
				fmt.Fprintf(os.Stderr, "%sScan Profile not found%s\n", chalk.Red, chalk.Reset)
			}
		}

	},
}

func exportScanProfile(id string, path string) {
	respCode, scanProfile := GetScanProfileRequest(id)
	var writePath string
	if respCode == 404 {
		fmt.Fprintf(os.Stderr, "Scan Profile not found\n")
	} else {
		jsonData, err := json.MarshalIndent(scanProfile, "", "  ")
		if err != nil {
			log.Fatalf("Error serializing struct to JSON: %v", err)
		} else {
			filename := fmt.Sprintf("%s.json", scanProfile.Name)
			// Ensure filename is valid and not empty
			if filename == ".json" {
				filename = "default.json" // Fallback filename
			}

			if path == "" {
				workingDir, err := os.Getwd()
				if err != nil {
					log.Fatalf("Error getting current working directory: %v", err)
				}
				writePath = filepath.Join(workingDir, filename)
			} else {
				// Check if the given path is a directory
				if stat, err := os.Stat(path); err == nil && stat.IsDir() {
					writePath = filepath.Join(path, filename)
				} else {
					writePath = path
				}
			}
			err = os.WriteFile(writePath, jsonData, 0644)
			if err != nil {
				log.Fatalf("Error writing JSON to file: %v", err)
			}
		}
	}
}

func GetScanProfileRequest(id string) (int, ScanProfile) {
	var respBody ScanProfile
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", viper.GetString("URL"), "/scanning_profiles/", id), nil)
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
	ScanProfileCmd.Flags().StringVarP(&id, "id", "", "", "Scan Profile ID")
	ScanProfileCmd.MarkFlagRequired("id")
	ScanProfileCmd.Flags().BoolP("export", "e", false, "Enable export")
	ScanProfileCmd.Flags().StringP("output", "o", ".", "Output directory")

	viper.BindPFlag("export", ScanProfileCmd.Flags().Lookup("export"))
	viper.BindPFlag("output", ScanProfileCmd.Flags().Lookup("output"))

	ScanProfileCmd.AddCommand(ListCmd)
	ScanProfileCmd.AddCommand(AddCmd)
	ScanProfileCmd.AddCommand(RemoveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanProfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanProfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
