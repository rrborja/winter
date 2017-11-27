// Copyright 2017 Ritchie Borja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package _test

import (
	//"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/rrborja/winter/metadata"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

type GoFileRegistry struct {
	Controller        string
	ControllerMethods []*ControllerMethodDescriptor
}

type InterpreterMemory struct {
	fset    ast.CommentMap
	current *ControllerMethodDescriptor
	stored  *GoFileRegistry
}

type ControllerMethodDescriptor struct {
	name          string
	Info          *metadata.Metadata
	Variables     map[string]*metadata.Metadata
	VariableTypes map[string]string
}

func (gfr *GoFileRegistry) add(cmd *ControllerMethodDescriptor) (err error) {
	if gfr.ControllerMethods == nil {
		gfr.ControllerMethods = make([]*ControllerMethodDescriptor, 0)
	}
	gfr.ControllerMethods = append(gfr.ControllerMethods, cmd)
	return
}

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

func (mdr *InterpreterMemory) Interpret(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.CommentGroup:
	case *ast.Comment:
		if mdr.current.Info == nil {
			var err error
			mdr.current.Info, err = metadata.ParseMetadata(n.Text)
			if err != nil {
				panic(err)
			}
		}
	case *ast.FuncDecl:
		method := new(ControllerMethodDescriptor)
		method.name = n.Name.Name
		method.Variables = make(map[string]*metadata.Metadata)
		mdr.current = method
	case *ast.FuncType:
		mdr.current.Variables = make(map[string]*metadata.Metadata, n.Params.NumFields())
		mdr.current.VariableTypes = make(map[string]string, n.Params.NumFields())
		for _, f := range n.Params.List {
			for _, name := range f.Names {
				for _, decField := range mdr.fset.Filter(f) {
					variable, err := metadata.ParseMetadata(decField[0].List[0].Text)
					if err != nil {
						panic(err)
					} else {
						if variable != nil {
							mdr.current.Variables[name.Name] = variable
							mdr.current.VariableTypes[name.Name] = f.Type.(*ast.Ident).Name
						}
					}
				}

			}
		}
	case *ast.BlockStmt:
		mdr.stored.add(mdr.current)
		mdr.current = nil
		return nil
	}
	return VisitorFunc(mdr.Interpret)
}

func (mdr InterpreterMemory) String() string {

	gfr := mdr.stored
	name := gfr.Controller
	methods := gfr.ControllerMethods

	methodInfo := make([]string, len(methods))
	vars := make([][]string, len(methods))

	for i, method := range methods {

		methodSignature := method.Info
		variables := method.Variables

		routeInfo := methodSignature.Info.(*metadata.RouteInfo)
		httpMethod := routeInfo.Method
		httpPath := routeInfo.Path

		methodInfo[i] = fmt.Sprintf("\tHandler:\t%s\n\t  Route:\t%s %s", method.name, metadata.ToStringOfHttpMethod(httpMethod), httpPath)

		vars[i] = make([]string, len(variables))

		j := 0
		for name, variable := range variables {
			variableInfo := variable.Info.(*metadata.VariableInfo)
			variableType := method.VariableTypes[name]
			if j < 1 {
				vars[i][j] = fmt.Sprintf("\tMapping:\t%s %s -> :%s", name, variableType, variableInfo.Name)
			} else {
				vars[i][j] = fmt.Sprintf("\t\t\t%s %s -> :%s", name, variableType, variableInfo.Name)
			}
			j++
		}

		methodInfo[i] = fmt.Sprintf("%s",
			strings.Join([]string{methodInfo[i], strings.Join(vars[i], "\n")}, "\n"))
	}

	return fmt.Sprintf("Controller name: %s\n\n%s", name, strings.Join(methodInfo, "\n\n"))
}

func TestParseComments(t *testing.T) {
	src := `
	package main

	type Login winter.Controller
	type Salaag winter.Controller

	//kkk
	//> GET /customers/:id
	func (login *Login) GetCustomer(
		id uint32, //> :id
		token string, //> :token
		coordinates string, //> :coordinates
	) (response winter.Response, err winter.Error) {
		return
	}

	//> GET /emails/:email
	func (login *Login) GetEmail(
		email string, //> :email
	) (response winter.Response, err winter.Error) {
		return
	}
	`
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "parser_test.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	mdr := new(InterpreterMemory)
	mdr.fset = ast.NewCommentMap(fset, f, f.Comments)
	mdr.stored = new(GoFileRegistry)

	ast.Print(fset, f)
	ast.Walk(VisitorFunc(mdr.Interpret), f)

	fmt.Println(mdr.String())

	assert.Fail(t, "none")
}
