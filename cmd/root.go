package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:   "xgo",
	Short: "xgo demo",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("welcome to use xgo, use `xgo -h` for help")
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
