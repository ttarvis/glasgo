// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	"go/token"
	"fmt"
	"encoding/hex"
	"encoding/base64"
)

func init() {
	register("hardcoded",
		"this is a test to look for suspected hardcoded credentials",
		hardcodedCheck,
		assignStmt, genDecl)
}

// isHighEntropy checks for string encoding to
// properly decode it and then measures its entropy
func isHighEntropy(s *string) bool {
	// entropyN is normalised entropy
	var entropy, entropyN float64;

	// if there is an error decoding, it could
	// be that the string matched the regular expression
	// for that encoding type but does not properly decode.
	// So, these reverse the normal situation.
	// if there isn't an error, it won't fallthrough;
	switch {
	case isHex(*s):
		buf, err := hex.DecodeString(*s);
		if err == nil {
			entropy, entropyN = H(buf);
			break;
		}
		fallthrough;
	case isBase64(*s):
		buf, err := base64.StdEncoding.DecodeString(*s);
		if err == nil {
			entropy, entropyN = H(buf);
			break;
		}
		fallthrough;
	default:
		buf := []byte(*s);
		entropy, entropyN = H(buf);	
	}		
	// these values, 2.5 and .98 are set through experiment
	// consider changing these or being more rigorous here
	if((entropy > 2.5) && (entropyN > .98)) {
		return true;
	}
	return false;
}

/*
func checkAssign(node *ast.AssignStmt) {

}
*/

// checkGenDecl starts with a GenDecl, checks if it is a const or a var
// then goes through its specs, then converts them to ValueSpecs
// then takes the Exprs from the ValueSpec and converts those to 
// BasicLits and then finally takes the actual string values for testing
// todo: could add check on the variable names first but this is debatable
func checkGenDecl(f *File, node *ast.GenDecl) {
	if (node.Tok == token.CONST || node.Tok == token.VAR) {
		for _, spec := range node.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				for _, expr := range valSpec.Values {
					if basicLit, ok := expr.(*ast.BasicLit); ok {
						suspectVal := basicLit.Value
						// strip quotes from suspectVal
						suspectVal = suspectVal[1: len(suspectVal) -1];
						// now check suspectVal 
						if isHighEntropy(&suspectVal) {
							f.Reportf(expr.Pos(), "Possible credential found: %s", suspectVal);
							return;
						}
						return;
					}
				}
			}
		}
	}
}

func hardcodedCheck(f *File, node ast.Node) {
	switch t := node.(type) {
	case *ast.AssignStmt:
		fmt.Println("assign", t.Rhs[0]);	
	case *ast.GenDecl:
		checkGenDecl(f, t);
	}
	return;
}

