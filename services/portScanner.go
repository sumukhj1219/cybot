package services

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func PortScanner(target string, startPort, endPort int, protocol string) {
	var wg sync.WaitGroup
	openPorts := make(chan int, endPort-startPort+1)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			if scanPort(protocol, target, port) {
				openPorts <- port
			}
		}(port)
	}

	go func() {
		wg.Wait()
		close(openPorts)
	}()

	for port := range openPorts {
		fmt.Printf("✅ Port %d is OPEN\n", port)
	}
}

func scanPort(protocol, target string, port int) bool {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout(protocol, address, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
