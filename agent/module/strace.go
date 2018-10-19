package module

import (
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/utils"
)

const syscallInsertStatement = `syscall,pid=%s,hostname=%s,syscall=%s hostname="%s",pid=%s,syscall="%s",%s,return=%s %s000`

var pidRegex *regexp.Regexp = regexp.MustCompile(`^\[pid\s+([0-9]+)\]`)
var timestampRegex *regexp.Regexp = regexp.MustCompile(`(\d{9,}\.\d{6}) `)
var syscallRegex *regexp.Regexp = regexp.MustCompile(` ([a-za-z_]+)\(`)
var argsRegex *regexp.Regexp = regexp.MustCompile(`[\( ]([a-z_\|]+|\d+|".*?\"|\{.*\})[,\)]`)
var returnRegex *regexp.Regexp = regexp.MustCompile(` = (-?\d+)`)

func getFirstSubmatch(matches []string) string {
	match := ""
	if len(matches) > 1 {
		match = matches[1]
	}
	return match
}

type strace struct {
	*base
}

func (m *strace) parseStraceOutputLine(line string) {
	pid := getFirstSubmatch(pidRegex.FindStringSubmatch(line))
	if len(pid) == 0 {
		log.Debug("Cannot parse line: %s", line)
		return
	}
	timestamp := getFirstSubmatch(timestampRegex.FindStringSubmatch(line))
	syscall := getFirstSubmatch(syscallRegex.FindStringSubmatch(line))
	returnValue := getFirstSubmatch(returnRegex.FindStringSubmatch(line))
	var args []string

	if len(syscall) == 0 || len(timestamp) == 0 || len(returnValue) == 0 {
		log.Debug("Cannot parse line: %s", line)
		return
	}

	for i, arg := range argsRegex.FindAllStringSubmatch(line, -1) {
		args = append(args, fmt.Sprintf(`arg%d="%s"`, i, strings.Replace(getFirstSubmatch(arg), `"`, `'`, -1)))
	}

	m.send(fmt.Sprintf(syscallInsertStatement, pid, utils.Hostname, syscall, utils.Hostname, pid, syscall, strings.Join(args, ","), returnValue, strings.Replace(timestamp, ".", "", 1)))
}

func (m *strace) Start(v url.Values) {
	// Force call clean as we are done here
	defer m.die()
	if len(v["pid"]) == 0 {
		log.Warn("Cannot proceed with empty pid, aborting")
		return
	}
	if len(v["tracerType"]) == 0 {
		log.Warn("Cannot proceed with empty tracerType, aborting")
		return
	}
	//Set name
	m.Name = "strace_" + v["tracerType"][0] + "_" + v["pid"][0]

	cmd := exec.Command("strace", "-ttt", "-s", "4096", "-v", "-yy", "-f", "-e", v["tracerType"][0], "-p", v["pid"][0])
	// Only returns when either shutdown get an integer or command exits
	cmdMonitor(cmd, m.shutdown, m.parseStraceOutputLine)
}

//NewStrace creates new ps module that will gather info
func NewStrace(m *base) *strace {
	return &strace{
		m,
	}
}
