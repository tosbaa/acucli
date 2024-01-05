/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package targetGroup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tosbaa/acucli/cmd/target"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type idResponseBody struct {
	TargetIDList []string `json:"target_id_list"`
}

var id string

// targetGroupCmd represents the targetGroup command
var TargetGroupCmd = &cobra.Command{
	Use:   "targetGroup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		responseCode, respBody := GetTargetGroupRequest(id)
		if responseCode == 200 {
			for _, targetID := range respBody.TargetIDList {
				_, target := target.GetTargetRequest(targetID)
				fmt.Printf("%s\t%s\n", targetID, target.Address)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%sTargetGroup not found%s\n", chalk.Red, chalk.Reset)
		}

	},
}

func GetTargetGroupRequest(id string) (int, idResponseBody) {
	var respBody idResponseBody
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/target_groups/%s/targets", httpclient.BASE_URL, id), nil)
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
		fmt.Println("Error making request:", err)
		return 404, respBody
	} else {
		return 200, respBody

	}
}

func init() {
	TargetGroupCmd.Flags().StringVarP(&id, "id", "", "", "Target ID")
	TargetGroupCmd.MarkFlagRequired("id")

	TargetGroupCmd.AddCommand(RemoveCmd)
	TargetGroupCmd.AddCommand(AddCmd)
	TargetGroupCmd.AddCommand(ListCmd)
	TargetGroupCmd.AddCommand(AddTargetsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetGroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetGroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
