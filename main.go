package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	version = "0.1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if !isAccessGrantedByIP(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Write([]byte(Config.Global.Version))
}

func main() {
	DefaultConfig()
	err := ReadConfig("conf/config.yaml")
	if err != nil {
		logf("warn", "read", http.StatusOK, "Unable to read configuration file: %s\n", err)
	}

	// set version by executing clamscan --version command in shell
	cmd := exec.Command("clamscan", "--version")
	out, err := cmd.Output()
	if err != nil {
		logf("crit", "start", http.StatusOK, "Unable to set version string: %s", err)
	}
	Config.Global.Version = fmt.Sprintf("%s (%s)\n", version, strings.TrimSuffix(string(out), "\n"))

	logf("info", "start", http.StatusOK, "Starting AVScan API")

	r := mux.NewRouter()

	r.HandleFunc("/", RenderIndex).Methods("GET")
	r.HandleFunc("/api/v1/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/v1/config", configHandler).Methods("GET")
	r.HandleFunc("/api/v1/version", versionHandler).Methods("GET")
	r.HandleFunc("/api/v1/scan", scanHandler).Methods("POST", "PUT")

	log.Fatal(http.ListenAndServe(Config.Server.Listen+":"+strconv.Itoa(Config.Server.Port), r))
}
