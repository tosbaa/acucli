/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scanProfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

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
		var scanProfile ScanProfile
		data := filehelper.ReadStdin()
		var stringBuilder strings.Builder
		for _, str := range data {
			stringBuilder.WriteString(str) // Add each string to the builder.
		}
		combinedString := stringBuilder.String()

		// Convert the combined string to a byte slice.
		byteSlice := []byte(combinedString)
		err := json.Unmarshal(byteSlice, &scanProfile)
		if err == nil {
			makeRequest(scanProfile)
		}
	},
}

func makeRequest(scanProfile ScanProfile) {
	requestJson, _ := json.Marshal(scanProfile)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", viper.GetString("URL"), "/scanning_profiles"), bytes.NewBuffer(requestJson))
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
		fmt.Println(chalk.Green, chalk.Bold.TextStyle("Successfully Added ScanProfile"), chalk.Reset)

	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Print(string(body))
	}
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
