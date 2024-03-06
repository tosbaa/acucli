/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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

type postBody struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

// addTargetsCmd represents the addTargets command
var AddTargetsCmd = &cobra.Command{
	Use:   "addTargets",
	Short: "Add targets to a target group",
	Long: `Add targets from stdin and the target group via id flag. Example:
	echo targets.txt| acucli targetGroup addTargets --id 0637a8b0-900d-44e8-9a04-edef6ac25e23 : Add targets from file adds the defined target group id
		`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ = cmd.Flags().GetString("id")
		input := filehelper.ReadStdin()
		if input != nil {
			pBody := postBody{}
			pBody.Add = input
			pBody.Remove = []string{}
			addTargets(pBody, id)
		}
	},
}

func addTargets(pBody postBody, id string) {
	requestJson, _ := json.Marshal(pBody)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/target_groups/%s/targets", viper.GetString("URL"), id), bytes.NewBuffer(requestJson))
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
		fmt.Println(chalk.Green, chalk.Bold.TextStyle("Successfully Added Targets:"), chalk.Reset)
		for _, target := range pBody.Add {
			fmt.Printf("%s%s%s\n", chalk.Green, target, chalk.Reset)
		}
	} else {
		fmt.Printf("%s%s", "failed", resp.Status)
	}
	defer resp.Body.Close()
}

func init() {
	AddTargetsCmd.Flags().StringVarP(&id, "id", "", "", "Group Target ID")
	AddTargetsCmd.MarkFlagRequired("id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTargetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTargetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
