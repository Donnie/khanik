package vidur

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"khanik/surang"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
)

var surangs = make(map[string]*surang.Surang)

func StartDaemon() {
	cntxt := &daemon.Context{
		PidFileName: "vidur.pid",
		PidFilePerm: 0644,
		LogFileName: "vidur.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Printf("Error starting daemon: %s\n", err)
		return
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	fmt.Println("Daemon started")
	runSurangManager()
}

func StopDaemon() {
	pidBytes, err := ioutil.ReadFile("vidur.pid")
	if err != nil {
		fmt.Println("Error reading PID file:", err)
		return
	}
	pid := strings.TrimSpace(string(pidBytes))
	cmd := exec.Command("kill", pid)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error stopping daemon:", err)
	} else {
		fmt.Println("Daemon stopped")
	}
	destroySurangs()
}

func ListSurangs() {
	loadSurangsFromConfig()
	for name, sur := range surangs {
		if surang.IsRunning(sur) {
			fmt.Printf("%s: Running on port %d\n", name, sur.Port)
		} else {
			fmt.Printf("%s: Not running (configured for port %d)\n", name, sur.Port)
		}
	}
}

func runSurangManager() {
	loadSurangsFromConfig()
	for {
		manageSurangs()
		time.Sleep(300 * time.Second)
	}
}

func loadSurangsFromConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}
	surangsConfig := viper.GetStringMap("surangs")
	for name, config := range surangsConfig {
		surangConfig := config.(map[string]interface{})
		sur := &surang.Surang{
			Name:     name,
			Command:  surangConfig["command"].(string),
			ExpectIP: surangConfig["expect_ip"].(string),
			Port:     surangConfig["port"].(int),
		}
		surangs[name] = sur
	}
}

func manageSurangs() {
	for _, sur := range surangs {
		fmt.Printf("Checking on port %d.", sur.Port)
		if !surang.IsRunning(sur) {
			surang.Start(sur)
		} else if !surang.Check(sur) {
			surang.Restart(sur)
		}
		fmt.Printf("All good!\n")
	}
}

func destroySurangs() {
	loadSurangsFromConfig()
	for _, sur := range surangs {
		if surang.IsRunning(sur) {
			surang.StopSurang(sur)
		}
	}
}
