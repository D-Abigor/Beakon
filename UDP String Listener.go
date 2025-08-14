package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	var input string
	listenPort := 52121
	bufferSize := 1024

	fmt.Println("listen port[\033[32m52121\033[0m]:")
	fmt.Scanln(&input)
	if input != "" {
		_, err := fmt.Sscanf(input, "%d", &listenPort)
		if err != nil {
			log.Fatalf("listen port input error %s", err)
		}
	}

	fmt.Println("buffer size[\033[32m1024\033[0m]:")
	fmt.Scanln(&input)
	if input != "" {
		_, err := fmt.Sscanf(input, "%d", &bufferSize)
		if err != nil {
			log.Fatalf("buffer size input error %s", err)
		}
	}

	listenAddr := net.UDPAddr{
		IP:   net.IPv4zero,
		Port: listenPort,
	}

	route, err := net.ListenUDP("udp4", &listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP port: %v: %v", listenAddr, err)
	}

	defer route.Close()

	log.Printf("Listening on UDP %d", listenAddr.Port)

	const maxUDPSize = 65535
	buffer := make([]byte, bufferSize)
	for {
		length, src, err := route.ReadFromUDP(buffer)
		if length == len(buffer) && len(buffer) < maxUDPSize {
			buffer = make([]byte, len(buffer)*2, maxUDPSize)
			log.Printf("\033[33m Buffer overflow: \033[0m buffer size has been increased to %d", len(buffer)*2)
			continue
		}
		if err != nil {
			log.Fatalf("\033[31mFatal Error\033[0m: could not read from UDP %s, %s", listenAddr.Port, err)
		}

		mssg := string(buffer[:length])
		log.Printf("\033[32m Recieved transmission from %v\033[0m \n message:  %s", src, mssg)
	}
}
