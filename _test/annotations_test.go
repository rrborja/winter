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
	"github.com/rrborja/winter/metadata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnnotationWithDeclarator(t *testing.T) {
	meta, _ := metadata.ParseMetadata(">GET /customers")
	assert.NotNil(t, meta)
}

func TestAnnotationWithIdentifier(t *testing.T) {
	meta, err := metadata.ParseMetadata(">                :ID")
	defer assert.NoError(t, err)

	variableInfo := meta.Info.(*metadata.VariableInfo)
	identifier := variableInfo.Name

	assert.Equal(t, "ID", identifier)
}

func TestExtractVariableInfo(t *testing.T) {
	//meta, err := metadata.ParseMetadata("> :customerId")
}

func TestMultiVariableInfo(t *testing.T) {
	meta, err := metadata.ParseMetadata("> :name :address")
	defer assert.NoError(t, err)

	variableInfo := meta.Info.(*metadata.MultiVariableInfo)
	assert.IsType(t, *variableInfo, metadata.MultiVariableInfo{})
}

func TestBlankSingleVariable(t *testing.T) {
	_, err := metadata.ParseMetadata("> :")
	assert.EqualError(t, err, "Syntax error. Expected an identifier after ':'")
}

func TestLastBlankSingleMultiVariable(t *testing.T) {
	_, err := metadata.ParseMetadata("> :phone :")
	assert.EqualError(t, err, "Syntax error. Expected an identifier after ':'")
}

func TestExpectHttpMethod(t *testing.T) {
	_, err := metadata.ParseMetadata(">  customer/:id")
	assert.EqualError(t, err, "Expected an Http Method before a path")
}

func TestValidRouteSyntax(t *testing.T) {
	_, err := metadata.ParseMetadata("> GET /customer/:id")
	assert.NoError(t, err)
}

func TestMetadataOfRoute(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> GET /customer")
	assert.Equal(t, "/customer", meta.Info.(*metadata.RouteInfo).Path.String())
}

func TestRouteParsedVariable(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> GET /customer/:id")
	_, hasKey := meta.Info.(*metadata.RouteInfo).Mapping["id"]
	assert.True(t, hasKey)
}

func TestRouteParsedGetHttpMethod(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> GET /customer/:id")
	method := meta.Info.(*metadata.RouteInfo).Method
	assert.Equal(t, metadata.Get{}, method)
}

func TestRouteParsedPostHttpMethod(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> POST /customer/:id")
	method := meta.Info.(*metadata.RouteInfo).Method
	assert.Equal(t, metadata.Post{}, method)
}

func TestRouteParsedPutHttpMethod(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> PUT /customer/:id")
	method := meta.Info.(*metadata.RouteInfo).Method
	assert.Equal(t, metadata.Put{}, method)
}

func TestRouteParsedDeleteHttpMethod(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> DELETE /customer/:id")
	method := meta.Info.(*metadata.RouteInfo).Method
	assert.Equal(t, metadata.Delete{}, method)
}

func TestRouteSyntaxWithQuery(t *testing.T) {
	_, err := metadata.ParseMetadata("> GET /customer/:id ? :token :customized")
	assert.NoError(t, err)
}

func TestRouteExtractionWithQuery(t *testing.T) {
	meta, _ := metadata.ParseMetadata("> GET /customer/:id ? :token :customized")
	queries := meta.Info.(*metadata.RouteInfo).Query
	assert.Equal(t, []string{"token", "customized"}, queries)
}
