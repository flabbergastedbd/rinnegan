package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tunnelshade/rinnegan/agent/daemon"
	"github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/utils"
	"net/url"
	"os/exec"
	"strings"
)

var name string

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Contorl modules",
	Long:  "Control rinnegan modules",
	Args:  cobra.ExactArgs(1),
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List modules",
	Long:  "List currently running rinnegan modules",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Listing all current modules")
		log.Info(daemon.HTTPGet("/module"))
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run modules",
	Long:  "Run rinnegan module with specified type and arguments",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running a module")
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop module",
	Long:  "Stop rinnegan module by name which can be obtained by running list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Switching to module operations")
		log.Info(daemon.HTTPDelete("/module/" + args[0]))
	},
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Collects all running processes along with sockets opened for listening",
	Long:  "Collects all running processes along with sockets opened for listening",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running ps module")
		log.Info(daemon.HTTPPost("/module/ps", url.Values{}))
	},
}

var straceCmd = &cobra.Command{
	Use:   "strace PID TRACE_TYPE",
	Short: "strace module",
	Long:  "strace module that captures syscalls for a pid, TRACE_TYPE is directly passed to strace -e option",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running strace module")
		data := url.Values{}
		data.Set("pid", args[0])
		data.Set("tracerType", args[1])
		for k, v := range data {
			log.Debug("%s: %s", k, strings.Join(v, ","))
		}
		log.Info(daemon.HTTPPost("/module/strace", data))
	},
}

var fridaCmd = &cobra.Command{
	Use:   "frida PID SCRIPT_NAME",
	Short: "frida module used to run various scripts",
	Long:  "frida module used to run various scripts. For list of available scripts have a look at build/frida",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Running frida module")
		data := url.Values{}
		data.Set("pid", args[0])
		data.Set("scriptName", args[1])
		for k, v := range data {
			log.Debug("%s: %s", k, strings.Join(v, ","))
		}
		log.Info(daemon.HTTPPost("/module/frida", data))
	},
}

func init() {
	if _, err := exec.LookPath("strace"); err != nil {
		log.Warn("Strace module not available as it is not found in $PATH")
	} else if utils.ReadFile("/proc/sys/kernel/yama/ptrace_scope") != "0" {
		log.Warn("ptrace_scope != 0, please set it for strace to work")
	} else {
		runCmd.AddCommand(straceCmd)
	}

	if _, err := exec.LookPath("frida"); err != nil {
		log.Warn("Frida module not available as it is not found in $PATH")
	} else {
		runCmd.AddCommand(fridaCmd)
	}
	runCmd.AddCommand(psCmd)
	moduleCmd.AddCommand(runCmd)
	moduleCmd.AddCommand(listCmd)
	moduleCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(moduleCmd)
}
