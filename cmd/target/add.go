/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
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
		input := filehelper.ReadStdin()
		inputGID, _ := cmd.Flags().GetString("gid")
		if inputGID != "" {
			groups = append(groups, inputGID)
		}

		if input != nil {
			targets := []Target{}

			for _, line := range input {
				targets = append(targets, Target{Address: line, Description: "", Type: "default", Criticality: 30})
			}
			makeRequest(targets, groups)
		} else {
			fmt.Println("Please provide addresses to add")
		}
	},
}

func makeRequest(t []Target, groups []string) {
	postBody := PostBody{Targets: t, Groups: groups}
	requestJson, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets/add"), bytes.NewBuffer(requestJson))
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
	AddCmd.Flags().StringVarP(&gid, "gid", "g", "", "Group ID (To assign the targets to the group)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
