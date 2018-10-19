package module

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/tunnelshade/rinnegan/agent/log"
)

type parser func(string)

func cmdMonitor(cmd *exec.Cmd, shutdown chan int, parse parser) {
	//Create command and get stderr pipe
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Warn("Bailing out as stderr pipe cannot be opened")
		return
	}

	err = cmd.Start()
	if err != nil {
		log.Warn("Failed to run strace")
		return
	}
	stderrReader := bufio.NewReader(stderrPipe)
	var line string
	//Monitor shutdown channel but parse output
loop:
	for {
		select {
		case <-shutdown:
			log.Info("Going to shutdown command")
			cmd.Process.Kill()
			break loop
		default:
			line, err = stderrReader.ReadString(byte('\n'))
			parse(line)
			if err == io.EOF {
				log.Warn("Strace died, cleaning up")
				// Read all lines and send them and break
				break loop
			}
		}
	}
	cmd.Wait()
}
