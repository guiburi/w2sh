package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	name    string
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "My root command",
		Long:  `My root command long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd Run with args: %v\n", args)
			fmt.Printf("Inside rootCmd Run with args: %v\n", name)
		},
	}

	subCmd = &cobra.Command{
		Use:   "sub",
		Short: "My subcommand",
		Long:  `My subcommand long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd Run with args: %v\n", args)
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Usage")
	RootCmd.AddCommand(subCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
