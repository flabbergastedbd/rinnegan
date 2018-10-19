package module

import (
	"fmt"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/utils"
)

const psInsertStatement = `process,pid=%d,hostname=%s hostname="%s",pid=%d,ppid=%d,user="%s",command="%s" %d`
const sockInsertStatement = `socket,pid=%d,hostname=%s hostname="%s",pid=%d,ip="%s",port=%d,protocol="%s",type="%s" %d`

type ps struct {
	*base
}

func (m *ps) Start(v url.Values) {
	m.Name = "ps"

	procs, err := process.Processes()
	if err != nil {
		log.Warn("Failed to get processes list: %s", err.Error())
	}
	for _, p := range procs {
		cmd, _ := p.Cmdline()
		cmd = strings.Replace(cmd, `"`, `'`, -1)
		ppid, _ := p.Ppid()
		user, _ := p.Username()
		m.send(fmt.Sprintf(psInsertStatement, p.Pid, utils.Hostname, utils.Hostname, p.Pid, ppid, user, cmd, time.Now().UnixNano()))
		connections, err := p.Connections()
		if err != nil {
			log.Warn("Failed to get connections for pid : %d", p.Pid)
			continue
		}
		for _, c := range connections {
			if c.Status == "LISTEN" {
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

				m.send(fmt.Sprintf(sockInsertStatement, p.Pid, utils.Hostname, utils.Hostname, p.Pid, c.Laddr.IP, c.Laddr.Port, proto, socketType, time.Now().UnixNano()))
			}
		}
	}
	// Force call clean as we are done here
	m.die()
}

//NewPs creates new ps module that will gather info
func NewPs(m *base) *ps {
	return &ps{
		m,
	}
}
