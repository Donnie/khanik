package vidur

import (
	"fmt"
	"khanik/surang"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
)

const (
	pidFile = "vidur.pid"
	logFile = "vidur.log"
)

var (
	surangs map[string]*surang.Surang
	mu      sync.RWMutex
)

// StartDaemon starts the surang manager daemon.
func StartDaemon() error {
	cntxt := &daemon.Context{
		PidFileName: pidFile,
		PidFilePerm: 0644,
		LogFileName: logFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}
	if d != nil {
		// Parent process exits.
		return nil
	}
	defer cntxt.Release()

	// Daemon process continues.
	return runSurangManager()
}

// StopDaemon stops the surang manager daemon.
func StopDaemon() error {
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("error reading PID file: %w", err)
	}
	pid := strings.TrimSpace(string(pidBytes))
	cmd := exec.Command("kill", pid)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error stopping daemon: %w", err)
	}
	if err := destroySurangs(); err != nil {
		return err
	}
	return nil
}

// RestartDaemon stops and then starts the surang manager daemon.
func RestartDaemon() error {
	// Stop the existing daemon
	if err := StopDaemon(); err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	// Start the daemon again
	if err := StartDaemon(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	return nil
}

// ListSurangs lists all configured surangs and their statuses.
func ListSurangs() error {
	if err := loadSurangsFromConfig(); err != nil {
		return err
	}
	mu.RLock()
	defer mu.RUnlock()
	for name, s := range surangs {
		running := s.IsRunning()
		status := "Not running"
		if running {
			status = "Running"
		}
		fmt.Printf("%s: %s on port %d\n", name, status, s.Port)
	}
	return nil
}

func runSurangManager() error {
	if err := loadSurangsFromConfig(); err != nil {
		return err
	}
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		if err := manageSurangs(); err != nil {
			fmt.Fprintf(os.Stderr, "Error managing surangs: %v\n", err)
		}
		<-ticker.C
	}
}

func loadSurangsFromConfig() error {
	mu.Lock()
	defer mu.Unlock()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	surangsConfig := viper.GetStringMap("surangs")
	surangs = make(map[string]*surang.Surang)
	for name, config := range surangsConfig {
		surangConfig := config.(map[string]interface{})
		port, ok := surangConfig["port"].(int)
		if !ok {
			port = int(surangConfig["port"].(float64))
		}
		s := &surang.Surang{
			Name:     name,
			Command:  surangConfig["command"].(string),
			ExpectIP: surangConfig["expect_ip"].(string),
			Port:     port,
		}
		surangs[name] = s
	}
	return nil
}

func manageSurangs() error {
	mu.RLock()
	defer mu.RUnlock()
	var wg sync.WaitGroup
	for _, s := range surangs {
		wg.Add(1)
		go func(s *surang.Surang) {
			defer wg.Done()
			if running := s.IsRunning(); !running {
				if err := s.Start(); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to start surang %s: %v\n", s.Name, err)
				} else {
					fmt.Printf("Started surang: %s on port %d\n", s.Name, s.Port)
				}
			} else if ok, err := s.Check(); !ok || err != nil {
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error checking surang %s: %v\n", s.Name, err)
				}
				if err := s.Restart(); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to restart surang %s: %v\n", s.Name, err)
				} else {
					fmt.Printf("Restarted surang: %s on port %d\n", s.Name, s.Port)
				}
			} else {
				fmt.Printf("Surang %s on port %d is running and healthy.\n", s.Name, s.Port)
			}
		}(s)
	}
	wg.Wait()
	return nil
}

func destroySurangs() error {
	if err := loadSurangsFromConfig(); err != nil {
		return err
	}
	mu.RLock()
	defer mu.RUnlock()
	var wg sync.WaitGroup
	for _, s := range surangs {
		wg.Add(1)
		go func(s *surang.Surang) {
			defer wg.Done()
			if err := s.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to stop surang %s: %v\n", s.Name, err)
			} else {
				fmt.Printf("Stopped surang: %s\n", s.Name)
			}
		}(s)
	}
	wg.Wait()
	return nil
}
