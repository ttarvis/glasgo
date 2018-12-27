package main

import (
	"net"
	"crypto/tls"
)

func bindAll() {
	tcpListener, err := net.Listen("tcp", "0.0.0.0:80");
	tlsListener, err := tls.Listen("tcp", "0.0.0.0:443", nil);

	if err != nil {
		err = nil;
	}
	if(tcpListener == nil || tlsListener == nil) {
		tcpListener = nil;
	}
}
