// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

func init() {
	register("sql",
		"this test checks for non constant sql query strings",
		sqlCheck,
		fileNode)
}

type sqlPackage struct {
	packageName	string
	argNames	[]string
	enabled		bool
	pkg		*types.Package
}

type SQLQuery struct {
	Func		*types.Func
	SSA		*ssa.Function
	ArgCount	int
	Param		int
}

var isCheckedSQL bool
var SQLCheckFailed bool
	
var sqlPackages = []sqlPackage{
	{
		packageName:	"database/sql",
		argNames:	[]string{"query"},
	},
}

// this may be called more than once, never do the same thing twice
// todo: dry, this will likely be used in other places
func getPkgImports(lp *loader.Program) map[string]*types.Package {
	pkgs := make(map[string]*types.Package);
	for _, pkg := range lp.AllPackages {
		if pkg.Importable {
			pkgs[pkg.Pkg.Path()] = pkg.Pkg;
		}
	} 
	return pkgs;
}

func GetQueries(sqlPackages sqlPackage, sqlPkg *types.Package, ssa *ssa.Program) []*SQLQuery {
	methods := make([]*SQLQuery, 0);
	scope := sqlPkg.Scope()
	for _, name := range scope.Names() {
		o := scope.Lookup(name);
		if !o.Exported() {
			continue;
		}
		if _, ok := o.(*types.TypeName); !ok {
			continue;
		}
		n := o.Type().(*types.Named);
		for i := 0; i < n.NumMethods(); i++ {
			m := n.Method(i);
			if !m.Exported() {
				continue;
			}
			sig := m.Type().(*types.Signature);
			if num, ok := FuncHasQuery(sqlPackages, sig); ok {
				methods = append(methods, &SQLQuery{
					Func:		m,
					SSA:		ssa.FuncValue(m),
					ArgCount:	sig.Params().Len(),
					Param:		num,
				});
			}
		}
	}
	return methods;
}

func FuncHasQuery(sqlPackages sqlPackage, sig *types.Signature) (int, bool) {
	params := sig.Params();
	for i := 0; i < params.Len(); i++ {
		v := params.At(i);
		for _, paramName := range sqlPackages.argNames {
			if v.Name() == paramName {
				// i is offset
				return i, true;
			}
		}
	}
	return 0, false
}

func GetNonConstantCalls(cGraph *callgraph.Graph, queries []*SQLQuery) []ssa.CallInstruction {
	cGraph.DeleteSyntheticNodes();

	suspected := make([]ssa.CallInstruction, 0);
	for _, m := range queries {
		node := cGraph.CreateNode(m.SSA);
		for _, edge := range node.In {

			isInternalSQLPkg := false
			for _, pkg := range sqlPackages {
				if pkg.packageName == edge.Caller.Func.Pkg.Pkg.Path() {
					isInternalSQLPkg = true;
					break
				}
			}
			if isInternalSQLPkg {
				continue
			}

			cc := edge.Site.Common()
			args := cc.Args;

			v := args[m.Param];

			if _, ok := v.(*ssa.Const); !ok {
				if inter, ok := v.(*ssa.MakeInterface); ok && types.IsInterface(v.(*ssa.MakeInterface).Type()) {
					if inter.X.Referrers() == nil || inter.X.Type() != types.Typ[types.String] {
						continue;
					}
				}
				suspected = append(suspected, edge.Site);
			}	
		}
	}

	return suspected;
}

func sqlCheck(f *File, node ast.Node) {
	if isCheckedSQL {
		return;
	}

	if (f.pkg.lp == nil) || (f.pkg.ssaProg == nil) || (f.pkg.cGraph == nil) {
		// skip this test
		// run a different one
		// mark the check as complete for now so this is not run again
		isCheckedSQL 	= true;
		SQLCheckFailed 	= true;
		warnf("unable to complete primary check for potential SQL injection");
		return;
	}

	imports := getPkgImports(f.pkg.lp);
	if len(imports) == 0 {
		return;
	}

	isSqlImported := false;
	for i := range sqlPackages {
		if _, ok := imports[sqlPackages[i].packageName]; ok {
			sqlPackages[i].enabled = true;
			isSqlImported = true;
			sqlPackages[i].pkg = imports[sqlPackages[i].packageName];
		}
	}
	if !isSqlImported {
		// maybe make a mention of not finding any SQL in use?
		// todo: see above
		return;
	}

	queries := make([]*SQLQuery, 0);

	for i := range sqlPackages {
		if sqlPackages[i].enabled {
			queries = append(queries, GetQueries(sqlPackages[i], sqlPackages[i].pkg, f.pkg.ssaProg)...);
		}
	}

	suspected := GetNonConstantCalls(f.pkg.cGraph, queries);

	for _, suspectCall := range suspected {
		f.Reportf(suspectCall.Pos(), "audit use of non-constant query: %s", suspectCall);
	}

	isCheckedSQL = true;

	return;
}

