package cmd

import "github.com/spf13/cobra"

var webCommand = &cobra.Command{
	Use:   "web",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init()  {
	rootCmd.AddCommand(webCommand)
}
