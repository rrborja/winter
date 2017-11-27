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
	"fmt"
	"strings"
)

type HttpMethod interface{}

type Get struct{}
type Post struct{}
type Put struct{}
type Delete struct{}

type IncompatibleMethod struct{}

type PathList []interface{}

func (pathList PathList) String() string {
	var relativePath string
	for _, path := range pathList {
		relativePath += fmt.Sprintf("%c%v", '/', path)
	}
	return relativePath
}

func ToStringOfHttpMethod(method HttpMethod) string {
	switch method.(type) {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	default:
		return "UNDEFINED"
	}
}

type RouteInfo struct {
	Method  HttpMethod
	Path    PathList
	Mapping map[string]*interface{}
	Query   []string
}

type Entry struct {
	key   string
	value interface{}
}

func (entry Entry) String() string {
	return fmt.Sprintf("%c%s", ':', entry.Text())
}

func (entry Entry) Text() string {
	return entry.key
}

func NewRouteInfo(method HttpMethod) RouteInfo {
	return RouteInfo{method, nil, nil, nil}
}

func (routeInfo *RouteInfo) ConcatenatePath(path string) {
	routeInfo.Path.add(path)
}

func (routeInfo *RouteInfo) ConcatenatePathVariable(variable string) {
	if routeInfo.Mapping == nil {
		routeInfo.Mapping = make(map[string]*interface{}, 1)
	}
	entry := Entry{variable, nil}
	routeInfo.Path.add(entry)
	routeInfo.Mapping[variable] = nil
}

func (routeInfo *RouteInfo) AddArgument(argument string) {
	if routeInfo.Query == nil {
		routeInfo.Query = []string{argument}
	} else {
		routeInfo.Query = append(routeInfo.Query, argument)
	}
}

func (pathList *PathList) add(object interface{}) {
	*pathList = append(*pathList, object)
}

func ToHttpMethod(method string) HttpMethod {
	switch strings.ToLower(method) {
	case "get":
		return Get{}
	case "post":
		return Post{}
	case "put":
		return Put{}
	case "delete":
		return Delete{}
	default:
		return IncompatibleMethod{}
	}
}
