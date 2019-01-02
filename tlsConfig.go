// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	"strconv"
)

// TLS is an acronym, therefore it should be in caps. It doesn't matter if it is exportable.
func init() {
	register("TLSConfig",
		"this is a check for insecure TLS configuration",
		iTLSConfigCheck,
		compositeLit)
}

const (
        VersionSSL30 = 0x0300
        VersionTLS10 = 0x0301
        VersionTLS11 = 0x0302
        VersionTLS12 = 0x0303
)

var secureCiphers = []string {
        "TLS_RSA_WITH_AES_128_CBC_SHA",
        "TLS_RSA_WITH_AES_256_CBC_SHA",
        "TLS_RSA_WITH_AES_128_CBC_SHA256",
        "TLS_RSA_WITH_AES_128_GCM_SHA256",
        "TLS_RSA_WITH_AES_256_GCM_SHA384",
        "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
        "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
        "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
        "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
        "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
        "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
        "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
        "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
        "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
        "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
        "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
        "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
	"0x002f", "0x0035", "0x003c", "0x009c", "0x009d", "0xc009", "0xc00a", "0xc013", "0xc014", "0xc023",
	"0xc027", "0xc02f", "0xc02b", "0xc030", "0xc02c", "0xcca8", "0xcca9",
}

func sliceContains(s string, slice []string) bool {
	for _, str := range secureCiphers {
		if s == str {
			return true;
		}
	}
	return false;
}

// checkTLSConfVal checks a key value expression from within a tls.Config composite element
// the idea is to check the values of these things for bad ones
func checkTLSConfVal(f *File, keyValueExpr *ast.KeyValueExpr) {
	// we are converting this down to identifiers where it is a 'basic' element type just in case
	if ident, ok := keyValueExpr.Key.(*ast.Ident); ok {
		switch ident.Name {

		case "InsecureSkipVerify":
			if val, ok := keyValueExpr.Value.(*ast.Ident); ok {
				if (val.Name == "true") {
					f.Reportf(keyValueExpr.Pos(), "InsecureSkipVerify is enabled, %s", f.ASTString(keyValueExpr));
				}
			} else {
				// value is not a basic identifier, so simple boolean values can't be checked
				f.Reportf(keyValueExpr.Pos(), "Audit use of InsecureSkipVerify, %s", f.ASTString(keyValueExpr));
			}
		case "PreferServerCipherSuites":
			if val, ok := keyValueExpr.Value.(*ast.Ident); ok {
				if val.Name == "false" {
					f.Reportf(keyValueExpr.Pos(), "PreferServerCipherSuites set to false, %s", f.ASTString(keyValueExpr));
				}
			} else {
				// can't be shown to be true; some sort of weird expression instead of simple true or false
				f.Reportf(keyValueExpr.Pos(), "Audit use of PreferServerCipherSuites, %s", f.ASTString(keyValueExpr));
			}
		case "MinVersion":
			if val, ok := keyValueExpr.Value.(*ast.BasicLit); ok {
				i, err := strconv.Atoi(val.Value);
				if err == nil {
					if ((int16)(i) < VersionTLS10) {
						f.Reportf(keyValueExpr.Pos(), "TLS minimum version is outdated, %s", f.ASTString(keyValueExpr));
					}
				}
			}
		case "MaxVersion":
			if val, ok := keyValueExpr.Value.(*ast.BasicLit); ok {
				i, err := strconv.Atoi(val.Value);
				if err == nil {
					if ((int16)(i) < VersionTLS11) {
						// todo: maybe reword this issue?
						f.Reportf(keyValueExpr.Pos(), "TLS maximum version is weak, %s", f.ASTString(keyValueExpr));
					}
				}
			}
		case "CipherSuites":
			if val, ok := keyValueExpr.Value.(*ast.CompositeLit); ok {
				for _, elt := range val.Elts {
					if cipherLit, ok := elt.(*ast.BasicLit); ok {
						if !sliceContains(cipherLit.Value, secureCiphers) {
							f.Reportf(cipherLit.Pos(), "Weak cipher, %s, is in use", f.ASTString(cipherLit));
						}
					}
				}	
			}
		}
	}
}

// iTLSConfigCheck checks TLS configuration structs.
// i stands for insecure. It also prevents this from being exportable.
func iTLSConfigCheck(f *File, node ast.Node) {
	if compLit, ok := node.(*ast.CompositeLit); ok && compLit.Type != nil {
		// elt stands for element, as in, composite element
		for _, elt := range compLit.Elts {
			if keyValueExpr, ok := elt.(*ast.KeyValueExpr); ok {
				// more succinct to call another function here
				checkTLSConfVal(f, keyValueExpr);
			}
		}
	}
}

