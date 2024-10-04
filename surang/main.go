package surang

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Surang struct {
	Name     string
	Command  string
	ExpectIP string
	Port     int
}

func IsRunning(surang *Surang) bool {
	cmd := exec.Command("pgrep", "-f", fmt.Sprintf("ssh -f -N -D %d %s", surang.Port, surang.Command))
	output, _ := cmd.Output()
	return len(output) > 0
}

func Start(surang *Surang) {
	cmd := exec.Command("ssh", "-f", "-N", "-D", strconv.Itoa(surang.Port), surang.Command)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting surang %s: %v\n", surang.Name, err)
	} else {
		fmt.Printf("Started surang: %s on port %d\n", surang.Name, surang.Port)
	}
}

func Restart(surang *Surang) {
	StopSurang(surang)
	Start(surang)
}

func StopSurang(surang *Surang) {
	cmd := exec.Command("pkill", "-f", fmt.Sprintf("ssh -f -N -D %d %s", surang.Port, surang.Command))
	cmd.Run()
	fmt.Printf("Stopped surang: %s\n", surang.Name)
}

func Check(surang *Surang) bool {
	cmd := exec.Command("curl", "-4", "-s", "--socks5", fmt.Sprintf("localhost:%d", surang.Port), "icanhazip.com")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error checking IP for surang %s: %v\n", surang.Name, err)
		return false
	}
	return strings.TrimSpace(string(output)) == surang.ExpectIP
}
