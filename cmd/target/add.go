/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type Target struct {
	Address     string `json:"address"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Criticality int    `json:"criticality"`
}

type PostBody struct {
	Targets []Target `json:"targets"`
	Groups  []string `json:"groups"`
}

var gid string
var target string

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add targets",
	Long: `Adding targets
	 you can give the URL or file path that includes targets line by line`,
	Run: func(cmd *cobra.Command, args []string) {
		groups := []string{}
		inputTarget, _ := cmd.Flags().GetString("target")
		inputGID, _ := cmd.Flags().GetString("gid")
		if inputGID != "" {
			groups = append(groups, inputGID)
		}

		targets := []Target{}
		// Check if the input is a file
		if isFile, filePath := isFilePath(inputTarget); isFile {
			// Read file contents
			contents, err := readFile(filePath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading file:")
				return
			}

			// Print array of file contents

			for _, line := range contents {
				targets = append(targets, Target{Address: line, Description: "", Type: "default", Criticality: 30})
			}
			makeRequest(targets, groups)
		} else {
			targets = append(targets, Target{Address: target, Description: "", Type: "default", Criticality: 30})
			makeRequest(targets, groups)
		}
	},
}

// Check if the input is a file path
func isFilePath(input string) (bool, string) {
	// Check if the file exists in the current directory
	if _, err := os.Stat(input); err == nil {
		absPath, _ := filepath.Abs(input)
		return true, absPath
	}
	return false, ""
}

// Read file contents and return as an array of strings
func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var contents []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}

func makeRequest(t []Target, groups []string) {
	postBody := PostBody{Targets: t, Groups: groups}
	requestJson, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", httpclient.BASE_URL, "/targets/add"), bytes.NewBuffer(requestJson))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.MyHTTPClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response struct {
			Targets []struct {
				Address     string `json:"address"`
				Criticality int    `json:"criticality"`
				Description string `json:"description"`
				FQDN        string `json:"fqdn"`
				Type        string `json:"type"`
				Domain      string `json:"domain"`
				TargetID    string `json:"target_id"`
				TargetType  string `json:"target_type"`
			} `json:"targets"`
		}

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&response); err != nil {
			fmt.Println("Error decoding response body:", err)
			return
		}

		fmt.Println(chalk.Green, chalk.Bold.TextStyle("Successfully Added Targets:"), chalk.Reset)
		for _, target := range response.Targets {
			fmt.Printf("%s%s\t%s%s\n", chalk.Green, target.Address, target.TargetID, chalk.Reset)
		}
	} else {
		fmt.Println("Error:", resp.Status)
	}
}

func init() {
	AddCmd.Flags().StringVarP(&target, "target", "t", "", "Target (Target to Add)")
	AddCmd.Flags().StringVarP(&gid, "gid", "g", "", "Group ID (To assign the targets to the group)")

	AddCmd.MarkFlagRequired("target")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
