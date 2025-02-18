package services

import (
	"fmt"
	"os/exec"
	"strings"
)

func ScanAttacks() {
	cmd := exec.Command("netstat", "-ano")
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ESTABLISHED") || strings.Contains(line, "LISTENING") {
			fmt.Println(line)
		}
	}
}
