/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new <appName>",
	Short:   "new an app in ../../apps",
	Example: "template new report",
	RunE:    newAPP,
}

func newAPP(cmd *cobra.Command, args []string) (err error) {
	fmt.Println("TODO new")
	fmt.Println(args)
	fmt.Println(name)
	return nil
}
func init() {
	newCmd.Flags().StringVarP(&name, "name", "n", "", "APP's name")
	newCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(newCmd)
}
