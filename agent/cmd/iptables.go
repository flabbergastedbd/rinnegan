package cmd

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/spf13/cobra"
	"github.com/tunnelshade/rinnegan/agent/log"
	"strings"
)

var chains = [3]string{"PREROUTING", "OUTPUT", "POSTROUTING"}

var iptablesCmd = &cobra.Command{
	Use:   "iptables",
	Short: "Interact with iptables",
	Long:  "Interact with iptables",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Interact with iptables")
	},
}

var listIptablesCmd = &cobra.Command{
	Use:   "list",
	Short: "List with iptables",
	Long:  "List with iptables",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Listing with iptables")
		ipt, err := iptables.New()
		if err != nil {
			log.Fatal("Issue using iptables: %s", err.Error())
		}
		for _, c := range chains[:] {
			rules, err := ipt.ListWithCounters("nat", c)
			if err != nil {
				log.Warn("Failed to get iptable rules for " + c)
				continue
			}
			log.Info("Chain: " + c + "\n\n" + strings.Join(rules, "\n") + "\n\n")
		}
	},
}

var incomingCmd = &cobra.Command{
	Use:   "incoming",
	Short: "Handle incoming traffic rules",
	Long:  "Handle incoming traffic rules with iptables",
	Args:  cobra.MinimumNArgs(1),
}

var delIncomingCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove incoming traffic redirects",
	Long:  "Remove incoming traffic redirects with iptables, call with <proto> <local-ip> <local-port> <remote_ip:port>",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Interact with iptables")
		ipt, err := iptables.New()
		if err != nil {
			log.Fatal("Issue using iptables: %s", err.Error())
		}
		ipt.Delete("nat", "POSTROUTING", "-j", "MASQUERADE")
		if err != nil {
			log.Warn("Issue deleting iptables for masquerade: %s", err.Error())
		}
		err = ipt.Delete("nat", "PREROUTING", "-p", args[0], "-d", args[1], "--dport", args[2], "-j", "DNAT", "--to-destination", args[3])
		if err != nil {
			log.Fatal("Failed to add incoming redirect")
		}
	},
}

var addIncomingCmd = &cobra.Command{
	Use:   "add",
	Short: "Redirect incoming traffic",
	Long:  "Redirect incoming traffic with iptables, call with <proto> <local-ip> <local-port> <remote_ip:port>",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Interact with iptables")
		ipt, err := iptables.New()
		if err != nil {
			log.Fatal("Issue using iptables: %s", err.Error())
		}
		err = ipt.AppendUnique("nat", "POSTROUTING", "-j", "MASQUERADE")
		if err != nil {
			log.Warn("Issue adding iptables for masquerade: %s", err.Error())
		}
		err = ipt.AppendUnique("nat", "PREROUTING", "-p", args[0], "-d", args[1], "--dport", args[2], "-j", "DNAT", "--to-destination", args[3])
		if err != nil {
			log.Fatal("Failed to add incoming redirect: %s", err.Error())
		}
	},
}

var outgoingCmd = &cobra.Command{
	Use:   "outgoing",
	Short: "Handle outgoing traffic rules",
	Long:  "Handle redirect outgoing traffic rules with iptables",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Interact with iptables")
	},
}

var addOutgoingCmd = &cobra.Command{
	Use:   "add",
	Short: "Redirect outgoing traffic",
	Long:  "Redirect outgoing traffic with iptables, call with <proto> <remote-ip> <remote-port> <redirect:port>",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Interact with iptables")
		ipt, err := iptables.New()
		if err != nil {
			log.Fatal("Issue using iptables: %s", err.Error())
		}
		err = ipt.AppendUnique("nat", "OUTPUT", "-p", args[0], "-d", args[1], "--dport", args[2], "-j", "DNAT", "--to-destination", args[3])
		if err != nil {
			log.Fatal("Failed to add incoming redirect")
		}
	},
}

var delOutgoingCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove rule redirecting outgoing traffic",
	Long:  "Remove rule redirecting outgoing traffic with iptables, call with <proto> <remote-ip> <remote-port> <redirect:port>",
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Interact with iptables")
		ipt, err := iptables.New()
		if err != nil {
			log.Fatal("Issue using iptables: %s", err.Error())
		}
		err = ipt.Delete("nat", "OUTPUT", "-p", args[0], "-d", args[1], "--dport", args[2], "-j", "DNAT", "--to-destination", args[3])
		if err != nil {
			log.Fatal("Failed to add incoming redirect")
		}
	},
}

func init() {
	outgoingCmd.AddCommand(addOutgoingCmd)
	outgoingCmd.AddCommand(delOutgoingCmd)
	incomingCmd.AddCommand(addIncomingCmd)
	incomingCmd.AddCommand(delIncomingCmd)
	iptablesCmd.AddCommand(listIptablesCmd)
	iptablesCmd.AddCommand(incomingCmd)
	iptablesCmd.AddCommand(outgoingCmd)
	rootCmd.AddCommand(iptablesCmd)
}
