/*
   main.go

   This is the stub.

   This code is released to the public domain. Originally prepared for the LCA 2015 conference
   by Mark Smith <mark@qq.is>.
*/

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
)

func handelConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		req, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read request: %s", err)
			}
			return
		}
		// 4. Connect to the backend web server.
		if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
			be_reader := bufio.NewReader(be)
			// 5. Send the request to the backend.
			if err := req.Write(be); err == nil {
				// 6. Read the response from the backend.
				if resp, err := http.ReadResponse(be_reader, req); err == nil {
					// 7. Send the response to the client, making sure to close it.
					resp.Close = true
					if err := resp.Write(conn); err == nil {
						log.Printf("proxied %s: got %d", req.URL.Path, resp.StatusCode)
					}
					conn.Close()
					// Repeat back at 2: accept the next connection.
				}
			}
		}
	}
}

func main() {
	// Code goes here.
	// Listen for connection forever.
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen: %s", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Failed to accept Connection: %s", err)
		}
		go handelConnection(conn)

	}
}
