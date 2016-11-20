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
	"strconv"
	"sync"
	"time"
)

var requestBytes map[string]int64
var requestLock sync.Mutex

type Backend struct {
	net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

var backendQueue chan *Backend

func init() {
	requestBytes = make(map[string]int64)
	backendQueue = make(chan *Backend, 10)
}

func getBackend() (*Backend, error) {
	select {
	case be := <- backendQueue:
		return be, nil
	case <- time.After(100 * time.Millisecond):
		be, err := net.Dial("tcp" , "127.0.0.1:8081")
		if err != nil {
			return nil, err
		}
		return &Backend {
			Conn: be,
			Reader: bufio.NewReader(be),
			Reader: bufio.NewWriter(be),
		}, nil
	}
}

func queueBackend(be *Backend) {
	select {
	case backendQueue <- be:
	case <- time.After(1 * time.Second):
		be.Close()
	}
}

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

		be, err := getBackend()
		if err != nil {
			return
		}
		if err := req.Write(be.Writer); err == nil {
			be.Writer.Flush()
		}
		go queueBackend(be)
		if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
			be_reader := bufio.NewReader(be)
			// 5. Send the request to the backend.
			if err := req.Write(be); err == nil {
				// 6. Read the response from the backend.
				if resp, err := http.ReadResponse(be_reader, req); err == nil {
					bytes := updateStats(req, resp)
					resp.Header.Set("X-Bytes", strconv.FormatInt(bytes, 10))
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

func updateStats(req *http.Request, resp *http.Response) int64 {
	requestLock.Lock()
	defer requestLock.Unlock()
	bytes := requestBytes[req.URL.Path] + resp.ContentLength
	requestBytes[req.URL.Path] = bytes
	return bytes
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

//Apachie Bench (AB chec how it works benchmarking)
//EPhimeral port exhaustion
//Edge storage Dropbox
//Dockers channels accross Network.
//Go Torch. Check how it works