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

	child = &cobra.Command{
		Use:   "child",
		Short: "My childcommand",
		Long:  `My childcommand long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nInside child Run with subname: %v\n", subName)
		},
	}

	sibling = &cobra.Command{
		Use:   "sibling",
		Short: "My sibling",
		Long:  `My sibling long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nInside sibling Run with subname: %v\n", subName)
		},
	}

	grandchild = &cobra.Command{
		Use:   "grand",
		Short: "My grandchild",
		Long:  `My grandchild long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nInside grandchild Run with subname: %v\n", subName)
		},
	}
)

func init() {
	RootCmd.Flags().StringVarP(&name, "name", "n", "", "Usage")
	child.Flags().StringVarP(&subName, "sub", "s", "", "Usage")
	RootCmd.AddCommand(child)
	RootCmd.AddCommand(sibling)
	child.AddCommand(grandchild)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
