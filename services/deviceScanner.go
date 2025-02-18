package services

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
)

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no valid local IP found")
}

func pingDevice(ip string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	output, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "bytes from") {
		results <- ip
	}
}

func ScanNetwork() {
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Println("❌ Error fetching local IP:", err)
		return
	}

	subnet := localIP[:strings.LastIndex(localIP, ".")+1]
	fmt.Println("\n🔍 Scanning network:", subnet+"0/24")

	var wg sync.WaitGroup
	results := make(chan string, 256)

	for i := 1; i <= 254; i++ {
		wg.Add(1)
		ip := fmt.Sprintf("%s%d", subnet, i)
		go pingDevice(ip, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\n🖥️ Active Devices:")
	found := false
	for ip := range results {
		fmt.Println("✅", ip)
		found = true
	}

	if !found {
		fmt.Println("⚠️ No active devices found. Try running as administrator.")
	}
}

func ScanNetworkARP() {
	cmd := exec.Command("arp", "-a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Failed to execute ARP command:", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	fmt.Println("\n🖥️ Active Devices (ARP Scan):")
	found := false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			fmt.Println("✅", fields[0])
			found = true
		}
	}

	if !found {
		fmt.Println("⚠️ No active devices found using ARP.")
	}
}
