// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	"strings"
)

func init() {
	register("bind",
		"this test checks for network listeners bound to all interfaces",
		bindCheck,
		callExpr)
}

func bindCheck(f *File, node ast.Node) {
	var names []string
	var callName string
	if call, ok := node.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			// SelectorExpr has two fields
			// X and Sel
			// X (through reflection) was found to be an Ident
			// Ident's have a field Name also.
			if id, ok := (fun.X).(*ast.Ident); ok {
				names = append(names, id.Name);
				names = append(names, fun.Sel.Name);
				callName = strings.Join(names, "/");
				if(callName == "net/Listen" || callName == "tls/Listen") {
					if len(call.Args) > 1 { // just a check to be sure
						if basicLit, ok := call.Args[1].(*ast.BasicLit); ok {
							if(strings.Contains(basicLit.Value, "0.0.0.0")) {
								callStr := f.ASTString(call);
								f.Reportf(node.Pos(), "audit binding network listener to all interfaces: %s", callStr);
							}
						}
					}
				}	
			}
		}
	}

	return;
}
