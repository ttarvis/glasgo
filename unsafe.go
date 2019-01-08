// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
)

func init() {
	register("unsafe",
		"this checks for use of the unsafe package",
		unsafeCheck,
		callExpr)
}

func unsafeCheck(f *File, node ast.Node) {
	// todo: dry this out.  The function name extraction needs to be moved
	// as similar code appears elsewhere
	if call, ok := node.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
                        // SelectorExpr has two fields
                        // X and Sel
                        // X (through reflection) was found to be an Ident
                        // Ident's have a field Name also.
                        if id, ok := (fun.X).(*ast.Ident); ok {
				// todo: this test could be more rigorous but may be sufficient
				if(id.Name == "unsafe") {
					callStr := f.ASTString(call);
					f.Reportf(node.Pos(), "audit use of unsafe package: %s", callStr);
				}
			}
		}
	}
	return;
}

