package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

var log = logrus.New()

var logFile = ""

func pong(w http.ResponseWriter, req *http.Request) {
	log.Info("ping -> pong")
}

func logcat(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Fprintf(writer, "could not open log file: %s", err)
	}
	fmt.Fprint(writer, string(data))
}

func killswitch(writer http.ResponseWriter, request *http.Request) {
	log.Fatal("kill")
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3000"
	}

	logFile = os.Getenv("LOG_FILE")
	if logFile == "" {
		log.Fatal("LOG_FILE must not be empty")
	}

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("could not open file %s: %s", logFile, err)
	}
	defer f.Close()

	log.Infof("openend log file %s", logFile)
	log.Out = f

	http.HandleFunc("/ping", pong)
	http.HandleFunc("/logcat", logcat)
	http.HandleFunc("/kill", killswitch)

	log.Infof("starting http server, listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Infof("stopping http server: %s", err)
	}
}
