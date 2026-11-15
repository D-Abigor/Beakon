package main

import (
	"fmt"
	"net"
	"syscall"
)

const broadcastPort = 52121

func getIPAndSubnet() (net.IP, net.IPMask, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ip := ipnet.IP.To4(); ip != nil && len(ipnet.Mask) == net.IPv4len {
				return ip, ipnet.Mask, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("could not find IPMask")
}

func calcBroadcastIP() (net.IP, net.IP, error) {
	ip, ipmask, err := getIPAndSubnet()
	if err != nil {
		return nil, nil, err
	}
	broadcast := make(net.IP, net.IPv4len)
	for i := range broadcast {
		broadcast[i] = ip[i] | ^ipmask[i]
	}
	return ip, broadcast, nil
}

func main() {
	ip, broadcastIP, err := calcBroadcastIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	raw, err := conn.SyscallConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	var socketErr error
	if err := raw.Control(func(fd uintptr) {
		socketErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	}); err != nil {
		fmt.Println(err)
		return
	}
	if socketErr != nil {
		fmt.Println(socketErr)
		return
	}

	_, err = conn.WriteToUDP([]byte(ip.String()), &net.UDPAddr{IP: broadcastIP, Port: broadcastPort})
	if err != nil {
		fmt.Println(err)
	}
}
