package utils

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
)

func GetPort(listenAddr string) (int, error) {
	index := strings.LastIndex(listenAddr, ":")
	if index == -1 {
		return 0, fmt.Errorf(`unable to get port based on the address`)
	}
	if result, err := strconv.ParseInt(listenAddr[index+1:], 10, 64); err != nil {
		return 0, fmt.Errorf(`unable to get port based on the address: %w`, err)
	} else {
		return int(result), nil
	}
}

func GetIP(isV4 bool) (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	interfaceName := "eth0"
	switch runtime.GOOS {
	case "darwin":
		interfaceName = "en0"
	}
	for _, _interface := range interfaces {
		if _interface.Name != interfaceName {
			continue
		}
		addrs, err := _interface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				v4 := ipnet.IP.To4()
				if isV4 {
					if v4 != nil {
						return v4, nil
					}
					continue
				}
				if v6 := ipnet.IP.To16(); v6 != nil {
					return v6, nil
				}
			}
		}
	}
	return nil, fmt.Errorf(`no available IP address was found`)
}

func GetAllIP(isV4 bool) ([]net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := make([]net.IP, 0)
	for _, _interface := range interfaces {
		addrs, err := _interface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				v4 := ipnet.IP.To4()
				if isV4 {
					if v4 != nil {
						result = append(result, v4)
					}
					continue
				}
				if v6 := ipnet.IP.To16(); v6 != nil {
					result = append(result, v4)
				}
			}
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf(`no available IP address was found`)
	}
	return result, nil
}
