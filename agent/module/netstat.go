package module

import (
	"fmt"
	"net/url"
	"strconv"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/utils"
)

const netstatInsertStatement = `netconn,pid=%d,hostname=%s hostname="%s",pid=%d,laddr="%s",lport=%d,raddr="%s",rport=%d,protocol="%s",type="%s" %d`

type netstat struct {
	*base
}

func (m *netstat) Start(v url.Values) {
	log.Debug("Starting netstat module")
	defer m.die()
	m.Name = "netstat"

	var pid int
	var err error
	var procs []*process.Process

	if len(v["pid"]) == 0 {
		log.Debug("No pids provided so tracing connections of all processes, not recommended")
		procs, err = process.Processes()
		if err != nil {
			log.Warn("Failed to get processes list: %s", err.Error())
			return
		}

	} else {
		m.Name = m.Name + "_" + v.Get("pid")
		pid, err = strconv.Atoi(v.Get("pid"))
		if err != nil {
			log.Warn("Unable to run netstat as no pid could be extracted from argument")
			return
		}
		pidProc, err := process.NewProcess(int32(pid))
		if err != nil {
			log.Warn("Unable to get process with pid %d", pid)
			return
		}
		procs = append(procs, pidProc)
	}

loop:
	for {
		select {
		case <-m.shutdown:
			log.Info("Going to shutdown netstat module")
			break loop
		default:
			for _, p := range procs {
				connections, err := p.Connections()
				if err != nil {
					log.Warn("Failed to get connections for pid : %d", p.Pid)
					continue
				}
				for _, c := range connections {
					if c.Status == "ESTABLISHED" {
						socketType, proto := "", ""

						switch c.Type {
						case syscall.SOCK_DGRAM:
							proto = "udp"
						case syscall.SOCK_STREAM:
							proto = "tcp"
						}

						switch c.Family {
						case syscall.AF_INET:
							socketType = "ipv4"
						case syscall.AF_INET6:
							socketType = "ipv6"
						case syscall.AF_UNIX:
							socketType = "unix"
						}

						m.send(fmt.Sprintf(netstatInsertStatement, p.Pid, utils.Hostname, utils.Hostname, p.Pid, c.Laddr.IP, c.Laddr.Port, c.Raddr.IP, c.Raddr.Port, proto, socketType, time.Now().UnixNano()))
					}
				}
			}
		}
		time.Sleep(time.Second * 2)
	}
	// Force call clean as we are done here
}

//NewNetstat creates new netstat module that will gather info
func NewNetstat(m *base) *netstat {
	return &netstat{
		m,
	}
}
