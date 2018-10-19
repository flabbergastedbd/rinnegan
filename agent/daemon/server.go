package daemon

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	log "github.com/tunnelshade/rinnegan/agent/log"
	"github.com/tunnelshade/rinnegan/agent/module"
)

type Daemon struct {
	running  bool
	modules  []module.Module
	wg       *sync.WaitGroup
	dbURL    string
	listener net.Listener
}

func (d *Daemon) Start() {
	ln, err := net.Listen("unix", Socket)
	if err != nil {
		log.Warn(err.Error())
		log.Fatal("Listening on socket failed")
	}
	d.listener = ln
	r := mux.NewRouter()
	r.HandleFunc("/daemon/stop", d.stop).Methods("POST")
	r.HandleFunc("/module", d.listModules).Methods("GET")
	r.HandleFunc("/module/{moduleType}", d.addModule).Methods("POST")
	r.HandleFunc("/module/{moduleName}", d.removeModule).Methods("DELETE")
	d.wg.Add(1)
	http.Serve(ln, r)
}

func (d *Daemon) listModules(w http.ResponseWriter, r *http.Request) {
	log.Debug("Listing Modules")
	var runningModules []module.Module
	for _, m := range d.modules {
		if m.IsRunning() {
			runningModules = append(runningModules, m)
			io.WriteString(w, m.GetName()+"\n")
		}
	}
	d.modules = runningModules
}

func (d *Daemon) removeModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("Removing module %s", vars["moduleName"])
	for i, m := range d.modules {
		if m.GetName() == vars["moduleName"] {
			m.Stop()
			d.modules = append(d.modules[:i], d.modules[i+1:]...)
			break
		}
	}
}

func (d *Daemon) addModule(w http.ResponseWriter, r *http.Request) {
	log.Debug("Adding a new module")
	vars := mux.Vars(r)
	moduleType := vars["moduleType"]
	m := module.Add(moduleType, d.dbURL, d.wg)
	d.modules = append(d.modules, m)
	//Add number to group at last
	d.wg.Add(1)
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	for k, v := range r.Form {
		log.Debug("%s: %s", k, strings.Join(v, ","))
	}
	go m.Start(r.Form)
}

func (d *Daemon) stop(w http.ResponseWriter, r *http.Request) {
	log.Info("Stopping daemon")
	log.Debug("Creating waitgroup")
	d.running = false
	for _, m := range d.modules {
		m.Stop()
		//wg.Add(1)
	}
	log.Debug("Waiting for all module to die")
	d.wg.Done()
	d.wg.Wait()
	log.Info("Closing listener")
	d.listener.Close()
	log.Info("Removing socket")
	os.Remove(Socket)
}

func New(dbURL string) *Daemon {
	return &Daemon{
		running: true,
		dbURL:   dbURL,
		modules: make([]module.Module, 0, 20),
		wg:      &sync.WaitGroup{},
	}
}
