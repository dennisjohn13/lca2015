/*
   main.go

   This is the stub.

   This code is released to the public domain. Originally prepared for the LCA 2015 conference
   by Mark Smith <mark@qq.is>.
*/

package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

func main() {
	// Code goes here.
	// Listen for connection forever.
	if ln, err := net.Listen("tcp", ":8080"); err == nil {
		//Infinite For loop  and accept connection over and over.
		for {
			// Accept Connection from client
			if conn, err := ln.Accept(); err == nil {
				reader := bufio.NewReader(conn)
				if req, err := http.ReadRequest(reader); err == nil {
					if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
						be_reader := bufio.NewReader(be)
						if err := req.Write(be); err == nil {
							if resp, err := http.ReadResponse(be_reader, req); err == nil {
								resp.Close = true
								if err := resp.Write(conn); err == nil {
									log.Printf("%s: %d", req.URL.Path, resp.StatusCode)
								}
								conn.Close()
							}
						}
					}
				}
			}
		}
	}
}
