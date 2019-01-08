// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
)

func init() {
	register("exec",
		"this checks for use of exec Command",
		execCheck,
		callExpr)
}

func execCheck(f *File, node ast.Node) {
	// todo: dry this out.  The function name extraction needs to be moved
	// as similar code appears elsewhere;
	if call, ok := node.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			// SelectorExpr has two fields
			// X and Sel
			// X (through reflection) was found to be an Ident
			// Idents have a field Name also
			if id, ok := (fun.X).(*ast.Ident); ok {
				if(id.Name == "exec") && (fun.Sel.Name == "Command") {
					callStr := f.ASTString(call);
					f.Reportf(node.Pos(), "audit use of os/exec package: %s", callStr);
				}
			}
		}
	}
	return;
}

