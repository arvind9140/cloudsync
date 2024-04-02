package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var token *oauth2.Token

func loginRequired(cmd *cobra.Command, args []string) error {
    if token == nil {
        return fmt.Errorf("login required, please run 'cloudsync login' first")
    }
    return nil
}


var rootCmd = &cobra.Command{
    Use:   "cloudsync",
    Short: "My CLI application",
    Long:  `A brief description of your application`,
	PreRunE: loginRequired,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Welcome to cloudSync CLI application!")
        fmt.Println("Please log in to continue.")
		fmt.Println("for login run 'login'")
		

        return nil
    },
}



func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
