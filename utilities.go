package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/opencontainers/selinux/go-selinux"
)

// Utilities
func jsonBuilder() string {
	openDemo := os.Getenv("SESSION_NAME")
	uid := strconv.Itoa(os.Getuid())
	gid := strconv.Itoa(os.Getgid())
	pid := strconv.Itoa(os.Getpid())
	pidLabel, _ := selinux.PidLabel(os.Getpid())
	osVersion := runtime.GOOS
	hostname, _ := os.Hostname()
	secretFile := readFile()

	jsonString := map[string]string{"emeaOpenDemoSession": varChecker(openDemo), "hostname": hostname, "userId": uid, "groupId": gid, "osVersion": osVersion, "pid": pid, "pidLabel:": pidLabel, "secretFileContent": secretFile}
	jsonResult, _ := json.MarshalIndent(jsonString, "", "   ")
	log.Println("JSON payload prepared, sending response...")
	return string(jsonResult)
}

func readFile() string {
	content, err := ioutil.ReadFile("/configdir/myconfigfile")
	if err != nil {
		return "THE SPECIFIED FILE WAS NOT FOUND"
	} else {
		return string(content)
	}
}

func varChecker(variable string) string {
	log.Println("Running variable check")
	if variable == "" {
		return "VARIABLE " + variable + "NOT SET"
	}
	return variable
}

func stressGenerator() {
	log.Println("Generating stress load")

	f, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	for i := 0; i < n; i++ {
		go func() {
			for {
				fmt.Fprintf(f, ".")
			}
		}()
	}
	log.Println("Sleeping before the response")
	time.Sleep(10 * time.Second)
}

func getReady(isReady *atomic.Value) bool {
	log.Println("Waiting 25 seconds to simulate a slow starting app")
	time.Sleep(25 * time.Second)
	isReady.Store(true)
	log.Println("App should now be ready!")
	return true
}
