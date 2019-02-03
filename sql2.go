// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	//"fmt"
	"regexp"
)

func init() {
	register("sqlBackup",
		"this is a backup test for the SQL injection test",
		sql2Check,
		funcDecl)
}

// May want to consider the possibility that someone gives db or Conn an alias
// Consider also checking if the sql package was imported
var expressions = []string{
	"^((Conn.|db.))*(Exec)|(Query)$",
	}

var regexps 	[]*regexp.Regexp;

func loadRegexps() {
	for _, expression := range expressions {
		rexp := regexp.MustCompile(expression);
		regexps = append(regexps, rexp);	
	}	
}

func isSQLCall(funcName string) bool {
	// move this somewhere else to be faster and only performed once
	loadRegexps();
	
	for _, re := range regexps {
		if matches := re.MatchString(funcName); matches {
			return true;
		}
	}
	return false;
}

func checkTainted(expr *ast.Expr, tainted *map[string]bool) bool {
	switch exp := (*expr).(type) {
	case *ast.BinaryExpr:
		x := exp.X;
		y := exp.Y;
		// yes this is recursive but how deeply nested to people actually write assignments?
		return (checkTainted(&x, tainted) || checkTainted(&y, tainted));
	case *ast.BasicLit:
		// if it is a constant, it is not tainted
		return false;
	case *ast.Ident:
		if(exp.Obj != nil) {
			varName := exp.Obj.Name;
			if _, isTainted := (*tainted)[varName]; isTainted {
				return true;
			}
		} 
	}
	return false;
}

func addTainted(exprs []ast.Expr, tainted *map[string]bool) {
	for _, expr := range exprs {
		if ident, ok := expr.(*ast.Ident); ok {
			(*tainted)[ident.Obj.Name] = true;
		}
	}
}

func sql2Check(f *File, node ast.Node) {
	// only run if the other SQL failed
	// todo: consider replacing the entirety of the other checker
	if(!SQLCheckFailed) {
		return;
	}

	tainted := make(map[string]bool);
	if fun, ok := node.(*ast.FuncDecl); ok {
		// perform a sanity check
		if(fun.Body == nil || len(fun.Body.List) < 1) {
			return;
		}

		// get input parameter names and types
		for _, field := range fun.Type.Params.List {
			// add params to tainted list
			if t, ok := field.Type.(*ast.Ident); ok {
				if(t.Name == "string") {
					for _, name := range field.Names {
						tainted[name.Name] = true;
					}
				}
			}
		}


		// get assignment statements and check if variables in assignments are tainted
		for _, statement := range fun.Body.List {
			if assign, ok := statement.(*ast.AssignStmt); ok {
				for _, expr := range assign.Rhs {
					switch exp := expr.(type) {
					case *ast.BinaryExpr, *ast.Ident:
						isTainted := checkTainted(&exp, &tainted);
						if(isTainted) {
							// add in variable name from lhs
							// right now it is just adding in all the Lhs
							addTainted(assign.Lhs, &tainted);
						}
					case *ast.BasicLit:
					case *ast.CallExpr:
					default:
					// do nothing for now
					// basic literals are safe
					// call expressions might need to be reconsidered
					}
				}
				// look for things that look like SQL query calls
				// also for variables to map
				for _, expr := range assign.Rhs {
					if call, ok := expr.(*ast.CallExpr); ok {
						funcName, err := getFullFuncName(call);
						if err != nil {
							// not sure what to do here
						}
						if(isSQLCall(funcName)) {
							// check if arguments are tainted
							for _, arg := range call.Args {
								// todo: should checkTainted just use value arguments not references?
								isTainted := checkTainted(&arg, &tainted);
								if(isTainted) {
									x := f.ASTString(expr);
									f.Reportf(expr.Pos(), "audit tainted input to SQL query, %s", x);
								}
							}
						}
					}
				}
			}

			// get things that are _probably_ sql exec calls
			if exprStmt, ok := statement.(*ast.ExprStmt); ok {
				if call, ok := exprStmt.X.(*ast.CallExpr); ok {
					var funcName string;
					funcName, err := getFullFuncName(call);
					if err != nil {
						funcName = getFuncName(call);	
					}
					if(isSQLCall(funcName)) {
					// extract parameters of the call
					// check if they are tainted
					}
				}
			}
		}

	}
	return;
}

