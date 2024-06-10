package src

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
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

	api_endpoint := os.Getenv("API_ENDPOINT")
	if api_endpoint == "" {
		api_endpoint = "http://app/api"
	}

	b, err := CheckDomainAndIp(api_endpoint, clientHello.ServerName, strings.Split(clientConn.RemoteAddr().String(), ":")[0])
	if err != nil {
		log.Println("Error checking domain:", err)
		return
	}

	if !b {
		log.Println("Access denied for domain:", clientHello.ServerName)
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

func HandleHTTPConnection() {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			target, err := url.Parse("http://" + req.Host)
			if err != nil {
				log.Println("Error parsing target URL:", err)
				return
			}

			b, err := CheckDomainAndIp("http://192.168.1.192:8000/api", req.Host, req.RemoteAddr)
			if err != nil {
				log.Println("Error checking domain:", err)
				return
			}

			if !b {
				log.Println("Access denied for domain:", req.Host)
				return
			}

			// Update the request to point to the target URL
			req.URL = target
			req.Host = target.Host
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Proxy server listening on port 80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
