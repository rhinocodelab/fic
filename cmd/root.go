/* github.com/rhinocodelab/fic/cmd/root.go */

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fic",
	Short: "A powerful tool to ensure file integrity on your system",
	Long: `A CLI application that helps to scan for changes,
and monitor file integrity in real-time.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'fic [command]' for options.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
