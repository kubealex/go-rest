package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)
	isReady := &atomic.Value{}
	isReady.Store(false)
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/stress", handleStress)
	r.HandleFunc("/healthz", handleHealthz)
	r.HandleFunc("/readyz", handleReadyz(isReady))
	r.HandleFunc("/unready", handleUnready(isReady))
	log.Printf("Webserver serving on port %s", "8080")
	done := make(chan bool)
	go func() {
		done <- getReady(isReady)
	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", "8080"), r))
	done <- true
}
