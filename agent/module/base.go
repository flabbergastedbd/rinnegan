package module

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/tunnelshade/rinnegan/agent/log"
)

type Module interface {
	Start(url.Values)
	Stop()
	GetName() string
	IsRunning() bool
}

type base struct {
	Name       string
	sendBuffer []string
	dbURL      string
	shutdown   chan int
	wg         *sync.WaitGroup
	running    bool
}

func (m *base) IsRunning() bool {
	return m.running
}

func (m *base) send(data string) {
	if len(data) > 0 {
		m.sendBuffer = append(m.sendBuffer, data)
	}

	if len(m.sendBuffer) == 10 || len(data) == 0 {
		data := strings.Join(m.sendBuffer, "\n")
		log.Debug("Sending following data to db\n%s", data)
		resp, err := http.Post(m.dbURL+"/write?db=rinnegan", "text/plain", strings.NewReader(data))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if len(body) > 0 {
			log.Warn("Sending: %s", data)
		}
		m.sendBuffer = m.sendBuffer[:0]
	}
}

func (m *base) Stop() {
	m.shutdown <- 1
}

func (m *base) GetName() string {
	return m.Name
}

func (m *base) die() {
	//Empty the sendBuffer by doing an empty string send
	m.send("")
	m.running = false
	//Should be last step to let daemon know we are done
	m.wg.Done()
}

func Add(moduleType string, dbURL string, wg *sync.WaitGroup) Module {
	log.Debug("Adding new module: %s", moduleType)
	m := &base{
		dbURL:    dbURL,
		shutdown: make(chan int, 1),
		wg:       wg,
		running:  true,
	}
	var md Module
	switch moduleType {
	case "ps":
		md = NewPs(m)
	case "strace":
		md = NewStrace(m)
	case "frida":
		md = NewFrida(m)
	case "alive":
		md = NewAlive(m)
	case "netstat":
		md = NewNetstat(m)
	}
	return md
}
