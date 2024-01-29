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

type RemovePostBody struct {
	TargetIDList []string `json:"target_id_list"`
}

// RemoveCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			makeDeleteRequest(input)
		}

	},
}

func makeDeleteRequest(ids []string) {
	postBody := RemovePostBody{TargetIDList: ids}
	requestJson, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/targets/delete"), bytes.NewBuffer(requestJson))
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
		fmt.Println(chalk.Red, chalk.Bold.TextStyle("Target Removed:"), chalk.Reset)
		for _, id := range ids {
			fmt.Printf("%s%s%s\n", chalk.Red, id, chalk.Reset)
		}
	} else {
		fmt.Println(resp)
	}
	defer resp.Body.Close()
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RemoveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// RemoveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
