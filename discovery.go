package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Device struct {
	DeviceID string
	Model    string
	Version  string
	RemoteIP net.IP
}

var discoveryStopSignal = make(chan os.Signal, 1)

func stopDiscovery() {
	discoveryStopSignal <- syscall.SIGINT
}

func discovery(timeout time.Duration) []Device {
	var err error
	var ret = []Device{}

	addr := net.UDPAddr{
		Port: 62001,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		return nil
	}

	// Stop at SIGINT
	signal.Notify(discoveryStopSignal, syscall.SIGINT)

	go func() {
		<-discoveryStopSignal
		fmt.Println("\nStopping...")
		_ = conn.Close()
	}()

	if timeout > 0 {
		t := time.NewTimer(timeout)
		go func() {
			<-t.C
			discoveryStopSignal <- syscall.SIGINT
		}()
	}

	newbuf := make([]byte, 5+1+1)
	copy(newbuf, "INGV\000")
	newbuf[5] = PKTTYPE_DISCOVERY

	broadcastDestination := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 62001}
	_, err = conn.WriteToUDP(newbuf, &broadcastDestination)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot send UDP reply to discovery:", err)
		_ = conn.Close()
		return nil
	}
	fmt.Println("Discovery packet sent, press CTRL-C to close")

	var buf [1024]byte

	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			break
		}

		if rlen > 5 && bytes.Compare(buf[:5], []byte("INGV\000")) == 0 {
			switch buf[5] {
			case PKTTYPE_DISCOVERY_REPLY:
				dev := Device{
					DeviceID: fmt.Sprintf("%x", buf[6:6+6]),
					Model:    strings.Trim(string(buf[6+6+4:6+6+4+8]), "\000"),
					Version:  string(buf[6+6 : 6+6+4]),
					RemoteIP: remote.IP,
				}
				ret = append(ret, dev)
				fmt.Println("Found:", dev.DeviceID, "- model:", dev.Model, "- version:", dev.Version, "- IP:", dev.RemoteIP)
			}
		}
	}
	return ret
}
