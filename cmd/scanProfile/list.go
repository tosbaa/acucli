/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scanProfile

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/httpclient"
)

type ScanProfiles struct {
	ScanningProfiles []ScanProfile `json:"scanning_profiles"`
}

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", viper.GetString("URL"), "/scanning_profiles"), nil)
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

		var scanProfiles ScanProfiles
		err = json.Unmarshal(body, &scanProfiles)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		for _, scanProfile := range scanProfiles.ScanningProfiles {
			fmt.Printf("%s\t%s\n", scanProfile.Name, scanProfile.ProfileID)
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
