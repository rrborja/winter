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

package metadata

import (
	"strings"
	"text/scanner"
)

type state uint8

const (
	Begin state = 1 << iota
	DeclaratorSymbol
	PathExpression
	QuerySymbol
	VariableTerm
	End
)

type expect uint8

const (
	NoneExpected expect = 1 << iota
	ExpectIdentifier
	ExpectPath
	ExpectArguments
)

func ParseMetadata(line string) (meta *Metadata, err error) {

	currentState := Begin
	expectedState := NoneExpected

	line = strings.TrimLeft(strings.TrimSpace(line), "/")

	var s scanner.Scanner
	s.Init(strings.NewReader(line))

	var tok rune

	if s.Scan() == '>' {
		currentState = DeclaratorSymbol

		for currentState < End && err == nil {
			tok = s.Scan()

			operator := tok

			switch operator {
			case '/':
				if meta.Type != Route {
					err = NewError("Syntax error. Use of '/' is only valid for Path declaration")
				} else if currentState == QuerySymbol {
					err = NewError("Syntax error. Use only identifiers in the query statement")
				} else {
					expectedState = ExpectIdentifier
				}
			case '?':
				currentState = QuerySymbol
			case ':':
				expectedState = ExpectIdentifier
				switch currentState {
				case DeclaratorSymbol:
					if meta != nil { //
						err = NewError("Bug error. Expected none declared identifiers.")
					}

					meta = NewMetadata(Variable)

					currentState = VariableTerm
				case PathExpression:
					if s.Scan() != scanner.Ident {
						err = NewError("Syntax error. Expected an identifier or variable in your route declaration")
					} else {
						identifier := s.TokenText()
						meta.Info.(*RouteInfo).ConcatenatePathVariable(identifier)
						expectedState = NoneExpected
					}
				case VariableTerm:
					if meta.Type != MultiVariable {
						meta.Info = NewMultiVariableInfo(meta.Info.(*VariableInfo).Name)
						meta.Type = MultiVariable
					}
				case QuerySymbol:
					expectedState = ExpectArguments
				default:
					err = NewError("Syntax error using ':'")
				}
			case scanner.Ident:
				if expectedState == ExpectIdentifier {
					expectedState = NoneExpected
				}
				switch currentState {
				case VariableTerm:
					if meta.Type == MultiVariable {
						multiVariableInfo := meta.Info.(*MultiVariableInfo)
						multiVariableInfo.addMultiVariableInfo(s.TokenText())
					} else {
						meta.Info = NewVariableInfo(s.TokenText())
					}
				case DeclaratorSymbol:
					keyword := s.TokenText()

					httpMethod := ToHttpMethod(keyword)

					switch httpMethod.(type) {
					case Get:
					case Post:
					case Put:
					case Delete:
					default:
						err = NewError("Expected an Http Method before a path")
					}

					expectedState = ExpectPath
					currentState = PathExpression

					routeInfo := NewRouteInfo(httpMethod)

					meta = NewMetadata(Route, routeInfo)

					//nextKeyword := s.TokenText()
				case PathExpression:
					meta.Info.(*RouteInfo).ConcatenatePath(s.TokenText())
				case PathExpression | VariableTerm:
					meta.Info.(*RouteInfo).ConcatenatePathVariable(s.TokenText())
				case QuerySymbol:
					expectedState = ExpectArguments | NoneExpected
					meta.Info.(*RouteInfo).AddArgument(s.TokenText())
				}
			case scanner.EOF:
				currentState = End
			}

		}
	} else {
		return
	}

	if err == nil && expectedState&ExpectIdentifier > 0 {
		if expectedState == ExpectIdentifier {
			err = NewError("Syntax error. Expected an identifier after ':'")
		} else {
			err = NewError("Syntax error.")
		}
	}

	return
}
