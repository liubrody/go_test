package main

import ( 
	"fmt"
	"net"
	"os"
	"github.com/cryptix/wav"
	"os/signal"
	"syscall"
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

var wrFile = wav.File {
	SampleRate: 		32000,
	Channels:		2,
	SignificantBits:	16,
}


func main() {
	quit := false
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println(s)
			quit = true
		default:
			fmt.Println("other", s)
		}
		}
	}()	

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP : net.IPv4zero, Port: 2050})
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Printf("Local : <%s> \n", listener.LocalAddr().String())


	data := make([]byte, 1500)
	var wavFile *os.File
	var wr *wav.Writer
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s frome<%s>",
				remoteAddr,  err)
		}
		//listener.WriteToUDP(data[:n], remoteAddr)

		if wr == nil {
			wavFile, err = os.OpenFile(
				"test.wav", 
				os.O_WRONLY | os.O_TRUNC | os.O_CREATE,
				0666,
			)
			if err != nil {
				fmt.Println(err)
				return 
			}
			defer wavFile.Close()
			
			if data[0] <= 32 {
				wrFile.Channels = uint16(data[0])
			}
			wr, err = wrFile.NewWriter(wavFile)
			if err != nil {
				fmt.Println(err)
				return 
			}
			defer wr.Close()
		}

		err = wr.WriteSample(data[8:n])
		if err != nil {
			fmt.Println(err)
		}

		if quit == true {
			break;
		}
	}

}
