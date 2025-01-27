package main

import (
	"github.com/pkg/browser"
	"github.com/shirou/gopsutil/process"
	"log"
	"sync"
	"time"
)

var browserMutex sync.Mutex
var appActive = false

func doMonitor(targetApp []string) {
	for {
		isAppRunning := false
		processes, err := process.Processes()
		if err != nil {
			continue
		}

		for _, proc := range processes {
			name, err := proc.Name()
			if err != nil {
				continue
			}

			for _, app := range targetApp {
				if name == app {
					isAppRunning = true
					handleTargetApp(proc)
				}
			}
		}

		if isAppRunning && !appActive {
			openBrowserOnce()
			appActive = true
		}

		if !isAppRunning && appActive {
			appActive = false
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func handleTargetApp(proc *process.Process) {
	if err := proc.Kill(); err != nil {
		log.Printf("Failed to kill process (PID: %d): %v", proc.Pid, err)
		return
	}
}

func openBrowserOnce() {
	browserMutex.Lock()
	defer browserMutex.Unlock()

	if err := browser.OpenURL(_Facebook); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
