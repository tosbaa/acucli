/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package target

import (
	"fmt"

	"github.com/spf13/cobra"
)

// targetCmd represents the target command
var TargetCmd = &cobra.Command{
	Use:   "target",
	Short: "A brief description of your command target",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For examplefddffddfdfdfdf:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("target called")
	},
}

func init() {
	TargetCmd.AddCommand(ListCmd)
	TargetCmd.AddCommand(AddCmd)
	TargetCmd.AddCommand(RemoveCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
