package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	PKTTYPE_DISCOVERY       = 1
	PKTTYPE_DISCOVERY_REPLY = 2
)

func main() {

	addr := net.UDPAddr{
		Port: 62001,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		return
	}

	// Stop at SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		_ = conn.Close()
	}()

	newbuf := make([]byte, 5+1+1)
	copy(newbuf, "INGV\000")
	newbuf[5] = PKTTYPE_DISCOVERY

	broadcastDestination := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 62001}
	_, err = conn.WriteToUDP(newbuf, &broadcastDestination)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot send UDP reply to discovery:", err)
		_ = conn.Close()
		return
	}
	fmt.Println("Discovery packet sent")

	var buf [1024]byte

	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			break
		}

		if rlen > 5 && bytes.Compare(buf[:5], []byte("INGV\000")) == 0 {
			switch buf[5] {
			case PKTTYPE_DISCOVERY_REPLY:
				deviceId := fmt.Sprintf("%x", buf[6:6+6])
				version := string(buf[6+6 : 6+6+4])
				model := strings.Trim(string(buf[6+6+4:6+6+4+8]), "\000")
				fmt.Println("Found:", deviceId, "- model:", model, "- version:", version, "- IP:", remote.IP)
			}
		}
	}

}
