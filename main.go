package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
)

// Proxy structure
type Proxy struct {
	listenAddr string
	targetAddr string
}

// NewProxy creates a new Proxy instance
func NewProxy(listenAddr, targetAddr string) *Proxy {
	return &Proxy{
		listenAddr: listenAddr,
		targetAddr: targetAddr,
	}
}

// Run starts the proxy server
func (p *Proxy) Run() error {
	ln, err := net.Listen("tcp", p.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Proxy listening on %s\n", p.listenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go p.handleClient(conn)
	}
}

// handleClient handles client connections
func (p *Proxy) handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	// Read the first bytes to determine if this is a TLS connection
	buf := make([]byte, 256)
	n, err := clientConn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	// Detect if the connection is using TLS
	isTLS := buf[0] == 0x16 // TLS handshake starts with 0x16

	// Reset the connection so that we can re-read the initial bytes later
	_, err = clientConn.Write(buf[:n])
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}

	// If it's a TLS connection, handle it accordingly
	if isTLS {
		p.handleTLS(clientConn)
		return
	}

	// For non-TLS connections, we can handle them differently if needed
	// For simplicity, we'll just forward them to the target address
	targetConn, err := net.Dial("tcp", p.targetAddr)
	if err != nil {
		fmt.Println("Error connecting to target:", err)
		return
	}
	defer targetConn.Close()

	// Forward data between client and target
	var wg sync.WaitGroup
	wg.Add(2)
	go p.pipe(targetConn, clientConn, &wg)
	go p.pipe(clientConn, targetConn, &wg)
	wg.Wait()
}

// handleTLS handles TLS connections
func (p *Proxy) handleTLS(clientConn net.Conn) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	tlsClient := tls.Server(clientConn, tlsConfig)
	err := tlsClient.Handshake()
	if err != nil {
		fmt.Println("TLS handshake error:", err)
		return
	}
	defer tlsClient.Close()

	// Extract SNI from the handshake
	serverName := tlsClient.ConnectionState().ServerName
	fmt.Printf("Received SNI: %s\n", serverName)

	// Connect to the actual target based on SNI
	targetConn, err := net.Dial("tcp", p.targetAddr)
	if err != nil {
		fmt.Println("Error connecting to target:", err)
		return
	}
	defer targetConn.Close()

	// Forward data between client and target
	var wg sync.WaitGroup
	wg.Add(2)
	go p.pipe(targetConn, tlsClient, &wg)
	go p.pipe(tlsClient, targetConn, &wg)
	wg.Wait()
}

// pipe copies data from src to dst
func (p *Proxy) pipe(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(dst, src)
	if err != nil && err != io.EOF {
		fmt.Println("Error copying data:", err)
	}
}

func main() {
	// Create a new proxy instance
	proxy := NewProxy(":443", "localhost:8443")

	// Start the proxy server
	if err := proxy.Run(); err != nil {
		fmt.Println("Error running proxy:", err)
	}
}
