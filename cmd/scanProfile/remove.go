/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scanProfile

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

// removeCmd represents the remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the given scanProfile",
	Long: `Removes the given scanProfile. Takes the ids line by line from stdin. Example:

cat scanProfileids.txt | acucli scanProfile remove`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			makeDeleteRequest(input)
		}
	},
}

func makeDeleteRequest(ids []string) {
	var allDeleted = true
	for _, id := range ids {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", viper.GetString("URL"), "/scanning_profiles/", id), nil)
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

		if resp.StatusCode == 404 {
			fmt.Printf("%sSpecified ScanProfile couldn't found%s\n", chalk.Red, chalk.Reset)
			allDeleted = false
			break
		}

		defer resp.Body.Close()
	}
	if allDeleted {
		fmt.Println(chalk.Red, chalk.Bold.TextStyle("ScanProfile Removed:"), chalk.Reset)
		for _, id := range ids {
			fmt.Printf("%s%s%s\n", chalk.Red, id, chalk.Reset)
		}
	}
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
