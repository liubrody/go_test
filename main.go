package main

import ( 
	"fmt"
	"net"
)

func GetLocalAddr() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}


	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}	
		//fmt.Println(addr.Network(), addr.String())
	}
}

func main() {


	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP : net.IPv4zero, Port: 2050})
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Printf("Local : <%s> \n", listener.LocalAddr().String())

	data := make([]byte, 1500)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}

		fmt.Printf("<%s> %d\n", remoteAddr, n)
	}
}