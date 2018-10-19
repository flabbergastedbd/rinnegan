package daemon

import (
	"github.com/tunnelshade/rinnegan/agent/utils"
	"path/filepath"
)

//Socket mostly unix kind where daemon listens for CRUD
//var Socket = filepath.Join(utils.GetBinDir(), "fifo", "daemon")
var socketDir = filepath.Join("/tmp/rinnegan", "fifo")
var Socket = filepath.Join(socketDir, "daemon")

func init() {
	utils.CreateIfNotExists(socketDir)
}
