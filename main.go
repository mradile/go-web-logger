package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	dataPath     = ""
	instanceFile = ""
	instanceID   = ""
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func pong(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(instanceID))
	log.Infof("ping -> pong -> %s", instanceID)
}

func kill(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "killing %s", instanceID)
	log.Fatalf("killing %s", instanceID)
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3000"
	}

	dataPath = os.Getenv("DATA_PATH")
	if dataPath == "" {
		log.Fatal("DATA_PATH must not be empty")
	}

	instanceFile = fmt.Sprintf("%s/instance", dataPath)
	inFi, err := os.OpenFile(instanceFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("could not open file %s: %s", instanceFile, err)
	}
	defer inFi.Close()
	rawID, err := ioutil.ReadAll(inFi)
	if err != nil {
		log.Fatalf("could not open file %s: %s", instanceFile, err)
	}
	if string(rawID) == "" {
		instanceID = randString(10)
		if _, err := inFi.Write([]byte(instanceID)); err != nil {
			log.Fatalf("could not write file %s: %s", instanceFile, err)
		}
	} else {
		instanceID = string(rawID)
		inFi.Close()
	}

	http.HandleFunc("/ping", pong)
	http.HandleFunc("/kill", kill)

	log.Infof("starting http server, listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Infof("stopping http server: %s", err)
	}
}

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	c := strings.Split(chars, "")
	b := make([]string, n)
	for i := range b {
		b[i] = c[rand.Intn(len(c))]
	}
	return strings.Join(b, "")
}
