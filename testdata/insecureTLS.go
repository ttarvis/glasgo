package main

import (
	"crypto/tls"
)

func insecureTLS() {
	config := &tls.Config{
		InsecureSkipVerify: 		true,
		PreferServerCipherSuites: 	false,
		MinVersion: 			0,
		MaxVersion:			1,
		CipherSuites:			[]uint16{0x000a, 0x0005, 0xc007, 0xc030},
	};

	tls.Dial("tcp", "mail.google.com:443", config);
}
