/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package targetGroup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type RemovePostBody struct {
	TargetGroupIDList []string `json:"group_id_list"`
}

var id string

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a target group",
	Long:  `You can remove a target group by giving its id as flag`,
	Run: func(cmd *cobra.Command, args []string) {
		var id_array = []string{}
		id, _ = cmd.Flags().GetString("id")
		// Check if the input is a file
		if isFile, filePath := filehelper.IsFilePath(id); isFile {
			// Read file contents
			contents, err := filehelper.ReadFile(filePath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading file:")
				return
			}

			// Print array of file contents

			for _, line := range contents {
				id_array = append(id_array, line)
			}
		} else {
			id_array = append(id_array, id)
		}
		makeDeleteRequest(id_array)

	},
}

func makeDeleteRequest(ids []string) {
	postBody := RemovePostBody{TargetGroupIDList: ids}
	requestJson, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", httpclient.BASE_URL, "/target_groups/delete"), bytes.NewBuffer(requestJson))
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
	if resp.StatusCode == 204 {
		fmt.Println(chalk.Red, chalk.Bold.TextStyle("Target Group Removed:"), chalk.Reset)
		for _, id := range ids {
			fmt.Printf("%s%s%s\n", chalk.Red, id, chalk.Reset)
		}
	} else {
		fmt.Println(resp)
	}
	defer resp.Body.Close()

}

func init() {
	RemoveCmd.Flags().StringVarP(&id, "id", "", "", "Target Group ID")
	RemoveCmd.MarkFlagRequired("id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RemoveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// RemoveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
