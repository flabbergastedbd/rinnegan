package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tunnelshade/rinnegan/agent/daemon"
	"github.com/tunnelshade/rinnegan/agent/log"
	"net/url"
)

var influxdb string

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Control daemon",
	Long:  "Control rinnegan daemon",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Switching to daemon operations")
	},
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start daemon",
	Long:  "Start rinnegan daemon",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting the daemon")
		if len(influxdb) == 0 {
			log.Fatal("No influxdb url provided")
		}
		daemon.New(influxdb).Start()
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop daemon",
	Long:  "Stop rinnegan daemon",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Stopping the daemon")
		log.Info(daemon.HTTPPost("/daemon/stop", url.Values{}))
	},
}

func init() {
	daemonStartCmd.PersistentFlags().StringVarP(&influxdb, "influxdb", "i", "", "Url for influxdb")
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	rootCmd.AddCommand(daemonCmd)
}
