/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tosbaa/acucli/helpers/filehelper"
	"github.com/tosbaa/acucli/helpers/httpclient"
	"github.com/ttacon/chalk"
)

type ConfigResponseBody struct {
	Description       string `json:"description"`
	LimitCrawlerScope bool   `json:"limit_crawler_scope"`
	Login             struct {
		Kind string `json:"kind"`
	} `json:"login"`
	Sensor         bool `json:"sensor"`
	SSHCredentials struct {
		Kind string `json:"kind"`
	} `json:"ssh_credentials"`
	Proxy struct {
		Enabled bool `json:"enabled"`
	} `json:"proxy"`
	Authentication struct {
		Enabled bool `json:"enabled"`
	} `json:"authentication"`
	ClientCertificatePassword string `json:"client_certificate_password"`
	ScanSpeed                 string `json:"scan_speed"`
	CaseSensitive             string `json:"case_sensitive"`
	Technologies              string `json:"technologies"`
	CustomHeaders             string `json:"custom_headers"`
	CustomCookies             string `json:"custom_cookies"`
	ExcludedPaths             string `json:"excluded_paths"`
	UserAgent                 string `json:"user_agent"`
	Debug                     bool   `json:"debug"`
}

// getConfigCmd represents the getConfig command
var GetConfigCmd = &cobra.Command{
	Use:   "getConfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		input := filehelper.ReadStdin()
		if input != nil {
			responseCode, respBody := getConfigRequest(input[0])
			if responseCode == 200 {
				filehelper.PrintStructFields(respBody)

			} else {
				fmt.Fprintf(os.Stderr, "%sTarget not found%s\n", chalk.Red, chalk.Reset)
			}

		}
	},
}

func getConfigRequest(i string) (int, ConfigResponseBody) {
	var respBody ConfigResponseBody
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/targets/%s/configuration", viper.GetString("URL"), i), nil)
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
		fmt.Println("Error parsing JSON:", err)
		return 404, respBody
	}
	if resp.StatusCode == 404 {
		return 404, respBody
	} else {
		return 200, respBody
	}
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
