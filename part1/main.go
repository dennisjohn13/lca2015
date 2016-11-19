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
	if ln, err := net.Listen("tcp",":8080") ; err == nil {
		for {
			// Accept Connection
			if conn , err := ln.Accept(); err == nil {

			}
		}
	}
}
