/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package targetGroup

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

type PostBody struct {
	Name string `json:"name"`
}

type ResponseBody struct {
	Name    string `json:"name"`
	GroupID string `json:"group_id"`
}

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			for _, targetGroupName := range input {
				postBody := PostBody{Name: targetGroupName}
				requestJson, _ := json.Marshal(postBody)
				req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/target_groups"), bytes.NewBuffer(requestJson))
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
				if resp.StatusCode == http.StatusCreated {
					var responseBody ResponseBody
					decoder := json.NewDecoder(resp.Body)
					if err := decoder.Decode(&responseBody); err != nil {
						fmt.Println("Error decoding response body:", err)
						return
					}
					fmt.Println(chalk.Green, chalk.Bold.TextStyle(responseBody.Name+" "+responseBody.GroupID), chalk.Reset)
				}
			}
		} else {
			fmt.Println("Please provide input")
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
