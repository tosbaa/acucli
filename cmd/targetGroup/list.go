/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package targetGroup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tosbaa/acucli/helpers/httpclient"
)

type responseBody struct {
	TargetGroups []struct {
		GroupID     string `json:"group_id"`
		Name        string `json:"name"`
		TargetCount int    `json:"target_count"`
		Description string `json:"description"`
		VulnCount   struct {
			Critical int `json:"critical"`
			High     int `json:"high"`
			Medium   int `json:"medium"`
			Low      int `json:"low"`
			Info     int `json:"info"`
		} `json:"vuln_count"`
	} `json:"groups"`
	Pagination struct {
		Count      int      `json:"count"`
		Cursors    []string `json:"cursors"`
		CursorHash string   `json:"cursor_hash"`
		Sort       string   `json:"sort"`
	} `json:"pagination"`
}

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the target groups",
	Long:  `Lists all the target groups with their name and their corresponding id to use it for other commands`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create an HTTP GET request using the custom client
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", httpclient.BASE_URL, "/target_groups"), nil)
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

		var responseBody responseBody
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
		for _, targetGroup := range responseBody.TargetGroups {
			fmt.Printf("'%s'\t%s\n", targetGroup.Name, targetGroup.GroupID)
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
