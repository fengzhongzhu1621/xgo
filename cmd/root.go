package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xgo",
	Short: "xgo demo",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("welcome to use xgo, use `xgo -h` for help")
		Start()
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
