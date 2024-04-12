package main

import (
	"fmt"
	"github.com/joseluiszuflores/stori-challenge/internal/bootstrap"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	//nolint:exhaustruct
	var cmdEvaluate = &cobra.Command{
		Use:   "evaluate",
		Short: "Evaluates the file and sends an email to user previously added",
		Long: `Sent email to user with  the number of transactions grouped by month,
and the average credit and average debit amounts grouped by month. `,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			path, err := cmd.Flags().GetString("path")
			if err != nil || path == "" {
				//nolint:forbidigo
				fmt.Println("The flag path is empty. Please consider add the path of file csv")

				os.Exit(-1)
			}

			user, err := cmd.Flags().GetString("user")
			if err != nil || user == "" {
				//nolint:forbidigo
				fmt.Println("The flag user is necessary to sent the email correctly")

				os.Exit(-1)
			}

			if err := bootstrap.Run(path, user); err != nil {
				//nolint:forbidigo
				fmt.Printf("Something happening err:[%s]", err)

				os.Exit(-1)
			}
			//nolint:forbidigo
			fmt.Print("Success sending the email with the information")
		},
	}
	cmdEvaluate.Flags().String("path", "", "is the path of file csv")
	cmdEvaluate.Flags().String("user", "", "is user id")
	//nolint:exhaustruct
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdEvaluate)
	//nolint:errcheck
	rootCmd.Execute()
}
