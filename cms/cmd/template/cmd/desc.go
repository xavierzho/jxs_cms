/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// descCmd represents the desc command
var descCmd = &cobra.Command{
	Use:     "desc <moduleName>",
	Short:   "describe what the app's module does",
	Example: "template desc api",
	RunE:    modelInfo,
}

func modelInfo(cmd *cobra.Command, args []string) (err error) {
	fmt.Println("TODO desc")
	fmt.Println(args)
	return nil
}

func init() {
	descCmd.Flags().StringP("module", "m", "", "APP's module name")
	descCmd.MarkFlagRequired("module")

	rootCmd.AddCommand(descCmd)
}
