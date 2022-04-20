package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// Handlers

func handleRoot(w http.ResponseWriter, r *http.Request) {

	jsonResponse := jsonBuilder()
	fmt.Fprintln(w, string(jsonResponse))
	log.Println("Serving root path /")
}

func handleStress(w http.ResponseWriter, r *http.Request) {
	log.Println("Stress test call received, filling up the cores!")
	stressGenerator()
	log.Println("Stress cycle gone!")
	w.WriteHeader(http.StatusOK)
}

func handleUnready(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		isReady.Store(false)
		log.Println("App not ready, received unset command")
		time.Sleep(60 * time.Second)
		isReady.Store(true)
		log.Println("App is now back online!")
		w.WriteHeader(http.StatusOK)
	}
}

// Health Checks

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("Healthcheck hit!")
}

func handleReadyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			log.Println("Readyness hit failed!")
			return
		}
		log.Println("Readyness hit!")
		w.WriteHeader(http.StatusOK)
	}
}
