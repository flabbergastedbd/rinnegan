package module

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"
	log "github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/utils"
)

const agentInsertStatement = `agent,hostname=%s,interface="%s" hostname="%s",interface="%s",ip="%s" %d`

type alive struct {
	*base
}

func (m *alive) Start(v url.Values) {
	// Force call clean as we are done here
	defer m.die()
	//Set name
	m.Name = "alive"

loop:
	for {
		select {
		case <-m.shutdown:
			log.Info("Going to shutdown alive agent")
			break loop
		default:
			// Get interfaces
			interfaces, err := net.Interfaces()
			if err != nil {
				log.Warn("Failed to get interfaces, will try in next loop")
				continue
			}
			for _, iface := range interfaces {
				var addrs []string
				for _, addr := range iface.Addrs {
					if addr.Addr != "" {
						addrs = append(addrs, addr.Addr)
					}
				}
				m.send(fmt.Sprintf(agentInsertStatement, utils.Hostname, iface.Name, utils.Hostname, iface.Name, strings.Join(addrs, ","), time.Now().UnixNano()))
			}
		}
		// Only run once for now
		time.Sleep(15 * time.Second)
		break
	}
	// Only returns when either shutdown get an integer or command exits
}

//NewStrace creates new ps module that will gather info
func NewAlive(m *base) *alive {
	return &alive{
		m,
	}
}
