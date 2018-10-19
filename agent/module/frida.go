package module

import (
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/utils"
)

type frida struct {
	*base
}

func getScriptPath(name string) string {
	return filepath.Join(utils.GetBinDir(), "frida", name+".js")
}

func (m *frida) Start(v url.Values) {
	// Force call clean as we are done here
	defer m.die()
	if len(v["pid"]) == 0 || len(v["scriptName"]) == 0 {
		log.Warn("Cannot proceed with empty pid or script name, aborting")
		return
	}
	path := getScriptPath(v["scriptName"][0])
	if _, err := os.Stat(path); err != nil {
		log.Warn("Cannot proceed with a non existant script: %s", path)
		return
	}
	//Set name
	m.Name = "frida_" + v["scriptName"][0] + "_" + v["pid"][0]

	cmd := exec.Command("frida", "-l", path, "-p", v["pid"][0])
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Warn("Cannot run frida script without stdin for clean close")
	}

	cmd.Start()
	//Wait for shutdown
	<-m.shutdown

	io.WriteString(stdin, "q\n")
}

//NewFrida creates new ps module that will gather info
func NewFrida(m *base) *frida {
	return &frida{
		m,
	}
}
