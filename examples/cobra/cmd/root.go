package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	name    string
	subName string
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "My root command",
		Long:  `My root command long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nInside rootCmd Run with name: %v\n", name)
		},
	}

	subCmd = &cobra.Command{
		Use:   "sub",
		Short: "My subcommand",
		Long:  `My subcommand long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nInside subCmd Run with subname: %v\n", subName)
		},
	}
)

func init() {
	RootCmd.Flags().StringVarP(&name, "name", "n", "", "Usage")
	subCmd.Flags().StringVarP(&subName, "sub", "s", "", "Usage")
	RootCmd.AddCommand(subCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
