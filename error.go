// Copyright 2018 Terence Tarvis.  All rights reserved.

package main

import (
	"go/ast"
	"go/types"
	"strings"
)

func init() {
	register("error",
		"this tests to see if any errors were ignored",
		errorCheck,
		assignStmt,
		exprStmt)
}

// isPrint checks to see if the call is a print statement
// this is used because people normally don't care about
// print statement errors
// todo: this could be more rigorous.  What if users invent
// their own print statement that could be damaging if errors
// are ignored? Consider also getting the full name and checking
// if the print statement is from the fmt package.
// Or maybe do exact matches for print statement methods in fmt.
// i.e. Println;
func isPrint(call *ast.CallExpr) bool {
	name := getFuncName(call);
	name = strings.ToLower(name);

	return strings.Contains(name, "print");
}

func returnsError(f *File, call *ast.CallExpr) int {
	if typeValue := f.pkg.info.TypeOf(call); typeValue != nil {
		switch t := typeValue.(type) {
		case *types.Tuple:
			for i := 0; i < t.Len(); i++ {
				variable := t.At(i)
				if variable != nil && variable.Type().String() == "error" {
					return i;
				}
			}
		case *types.Named:
			if t.String() == "error"{
				return 0;
			}
		}	
	}
	return -1;
}
 
// Possibly check if anything returns an error before running the test
// however, this may take roughly the same amount of effort as
// just running the test in the first place.
//
func errorCheck(f *File, node ast.Node) {
	switch stmt := node.(type) {
	case *ast.AssignStmt:
		for _, rhs := range stmt.Rhs {
			if call, ok := rhs.(*ast.CallExpr); ok {
				index := returnsError(f, call)
				if index < 0 {
					continue
				}
				// ignore print calls unless verbose
				if isPrint(call) {
					if(!(*verbose)) {
						continue;
					}
				} 
				lhs := stmt.Lhs[index]
				if id, ok := lhs.(*ast.Ident); ok && id.Name == "_" {
					// todo real reporting
					re := f.ASTString(rhs);
					le := f.ASTString(lhs);
					f.Reportf(stmt.Pos(), "error ignored %s %s", le, re);
				}
			}
		}
	case *ast.ExprStmt:
		if expr, ok := stmt.X.(*ast.CallExpr); ok {
			pos := returnsError(f, expr);
			if pos >= 0 {
				// todo: real reporting
				// ignore print statements unless verbose
				if isPrint(expr) {
					if(!(*verbose)) {
						return;
					}
				}
				x := f.ASTString(expr);
				f.Reportf(stmt.Pos(), "error ignored %s", x);
			}
		}
	}
}
