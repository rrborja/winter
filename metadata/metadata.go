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

type Metadata struct {
	Type        MetadataType
	Info        interface{}
	ErrorReason string
}

type VariableInfo struct {
	Name string
}

type MultiVariableInfo struct {
	Names []string
}

type MetadataType int

const (
	Route MetadataType = iota
	Variable
	MultiVariable
)

func NewMetadata(metaDataType MetadataType, routeInfos ...RouteInfo) *Metadata {
	var routeInfo *RouteInfo
	switch len(routeInfos) {
	case 0:
	case 1:
		routeInfo = &routeInfos[0]
	default:
		panic("More than 1 argument for RouteInfo argument")
	}
	return &Metadata{metaDataType, routeInfo, ""}
}

func NewVariableInfo(name string) *VariableInfo {
	return &VariableInfo{name}
}

func NewMultiVariableInfo(names ...string) *MultiVariableInfo {
	return &MultiVariableInfo{names}
}

func (multiVariableInfo *MultiVariableInfo) addMultiVariableInfo(names string) {
	multiVariableInfo.Names = append(multiVariableInfo.Names, names)
}

func NewError(reason string) Metadata {
	meta := NewMetadata(-1)
	meta.ErrorReason = reason
	return *meta
}

func (metadata *Metadata) String() string {
	return "none"
}

func (metadata Metadata) Error() string {
	return metadata.ErrorReason
}
