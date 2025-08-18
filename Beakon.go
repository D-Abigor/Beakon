package main

import (
	"fmt"
	"net"
)

func getIPAndSubnet() (net.IP, net.IPMask, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.To4() != nil {
				return ipnet.IP, ipnet.Mask, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("could not find IPMask")
}

func calcBroadcastIP() net.IP {
	ip, ipmask, _ := getIPAndSubnet()
	broadcast := make(net.IP, len(ip))

	for i := 0; i < len(ip); i++ {
		broadcast[i] = ip[i] + ^ipmask[i]
	}
	return broadcast
}

func main() {
	broadcastIP := calcBroadcastIP()
}
