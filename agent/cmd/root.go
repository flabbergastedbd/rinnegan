package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tunnelshade/rinnegan/agent/log"
)

var verbose bool = false

var rootCmd = &cobra.Command{
	Short: "Rinnegan agent CLI",
	Long:  "Rinnegan agent that handles daemon and interactions with it on targets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Do stuff
		log.Info("Rinnegan agent")
	},
}

func init() {
	cobra.OnInitialize(setLogLevel)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func setLogLevel() {
	if verbose == true {
		log.EnableDebug()
	}
}

func Execute() {
	rootCmd.Execute()
}
