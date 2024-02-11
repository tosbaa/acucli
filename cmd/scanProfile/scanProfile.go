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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
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

var id string

// scanProfileCmd represents the scanProfile command
var ScanProfileCmd = &cobra.Command{
	Use:   "scanProfile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ = cmd.Flags().GetString("id")
		if cmd.Flags().Changed("export") {
			output, _ := cmd.Flags().GetString("output")
			exportScanProfile(id, output)
		} else {
			responseCode, respBody := GetScanProfileRequest(id)
			if responseCode == 200 {
				filehelper.PrintStructFields(respBody)
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
		jsonData, err := json.Marshal(scanProfile)
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
