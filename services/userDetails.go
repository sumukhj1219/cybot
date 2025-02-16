package services

import (
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
)

type User struct {
	LocalIP  string
	PublicIP string
	OS       string
	Arch     string
	Hostname string
}

func UserDetails() (*User, error) {
	var u User

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.New("❌ Error in fetching local IP")
	}
	for _, add := range addrs {
		if ipNet, ok := add.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				u.LocalIP = ipNet.IP.String()
			}
		}
	}

	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return nil, errors.New("❌ Error in fetching public IP")
	}
	defer resp.Body.Close()

	ip, _ := io.ReadAll(resp.Body)
	u.PublicIP = string(ip)

	u.OS = runtime.GOOS
	u.Arch = runtime.GOARCH
	u.Hostname, _ = os.Hostname()

	return &u, nil
}
