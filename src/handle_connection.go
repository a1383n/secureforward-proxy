package src

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func HandleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	if err := clientConn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Print(err)
		return
	}

	clientHello, clientReader, err := peekClientHello(clientConn)
	if err != nil {
		log.Print(err)
		return
	}

	if err := clientConn.SetReadDeadline(time.Time{}); err != nil {
		log.Print(err)
		return
	}

	ipAddress, err := resolveIPAddress(clientHello.ServerName)
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("Connecting to %s on %s\n", clientHello.ServerName, ipAddress)
	backendConn, err := net.DialTimeout("tcp", net.JoinHostPort(ipAddress, "443"), 5*time.Second)
	if err != nil {
		log.Print(err)
		return
	}
	defer backendConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(clientConn, backendConn)
		clientConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()
	go func() {
		io.Copy(backendConn, clientReader)
		backendConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()

	wg.Wait()
}
