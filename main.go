package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile)

	certPath := "server.crt"
	keyPath := "server.key"
	listen := ":443"
	prefix := "echo: "

	flag.StringVar(&certPath, "cert", "server.crt", "Path to TLS server certificate")
	flag.StringVar(&keyPath, "key", "server.key", "Path to TLS server key")
	flag.StringVar(&listen, "listen", ":443", "Address/Port to listen on (i.e. :443 or 127.0.0.1:443)")
	flag.StringVar(&prefix, "prefix", "echo: ", "Prefix all echo'ed data with this value")
	flag.Parse()

	certPathAbs, err := filepath.Abs(certPath)
	if err != nil {
		log.Fatalf("Failed to convert path %q to absolute path: %v\n", certPath, err)
	}
	certPath = certPathAbs
	keyPathAbs, err := filepath.Abs(keyPath)
	if err != nil {
		log.Fatalf("Failed to convert path %q to absolute path: %v\n", keyPath, err)
	}
	keyPath = keyPathAbs

	cer, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load X509 key pair - Certificate: %s, Key: %s - %v\n", certPath, keyPath, err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", listen, config)
	if err != nil {
		log.Fatalf("Failed to setup TLS listener %s: %v\n", listen, err)
	}
	defer ln.Close()

	log.Printf("Now accepting TLS connections for %s\n", listen)
	connNum := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept new connection: %v\n", err)
			continue
		}
		go handleConnection(conn, connNum, prefix)
		connNum += 1
	}
}

func handleConnection(conn net.Conn, connNum int, prefix string) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read from connection (%d): %v\n", connNum, err)
			}
			return
		}

		conn.Write([]byte(fmt.Sprintf("%s%s", prefix, msg)))
	}
}
