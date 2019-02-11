// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	"go/token"
	"fmt"
	"reflect"
)

func init() {
	register("intToStr",
		"check if integers are being converted to strings using string()",
		intToStrCheck,
		callExpr)
}

func intToStrCheck(f *File, node ast.Node) {
	formatString := "integer possibly converted improperly: %s";
	if stmt, ok := node.(*ast.CallExpr); ok {
	// technically, string() is not a function but a type conversion
		name := getFuncName(stmt);
		if(name == "string") {
			// length of args to string() is only 1
			if(len(stmt.Args) == 1) {
				switch arg := stmt.Args[0].(type) {
				case *ast.Ident:
					if t := f.pkg.info.TypeOf(arg); t != nil {
						// is this really the best way to check?
						if(t.String() == "int") {
							str := f.ASTString(stmt);
							f.Reportf(stmt.Pos(), formatString, str);
						}
					}
				case *ast.BasicLit:
					if(arg.Kind == token.INT) {
						str := f.ASTString(stmt);
						f.Reportf(stmt.Pos(), formatString, str);
					}
				case *ast.CallExpr:
					if t := f.pkg.info.TypeOf(arg); t != nil {
						if(t.String() == "int") {
							str := f.ASTString(stmt);
							f.Reportf(stmt.Pos(), formatString, str);
						}
					}
				default:
					// todo: figure out what the other cases here are
					// ast.CompositeLit, ast.SliceExpr, ast.IndexExpr
					// ast.StartExpr, ast.SelectorExpr
					// ast.TypeAssertExpr, ast.BinaryExpr
				}
			}
		}
	} else {
		warnf("something strange happened at %s, please report", f.loc(stmt.Pos()) );
	}
	return;
}
