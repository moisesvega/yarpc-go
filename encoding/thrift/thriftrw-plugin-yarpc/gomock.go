// Copyright (c) 2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"path/filepath"

	"go.uber.org/thriftrw/plugin"
)

const gomockTemplate = `
// Code generated by thriftrw-plugin-yarpc
// @generated

<$pkgname := printf "%stest" (lower .Name)>
package <$pkgname>

<$gomock := import "github.com/golang/mock/gomock">

// MockClient implements a gomock-compatible mock client for service
// <.Name>.
type MockClient struct {
	ctrl *<$gomock>.Controller
	recorder *_MockClientRecorder
}

var _ <import .ClientPackagePath>.Interface = (*MockClient)(nil)

type _MockClientRecorder struct {
	mock *MockClient
}

// Build a new mock client for service <.Name>.
//
// 	mockCtrl := gomock.NewController(t)
// 	client := <$pkgname>.NewMockClient(mockCtrl)
//
// Use EXPECT() to set expectations on the mock.
func NewMockClient(ctrl *<$gomock>.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &_MockClientRecorder{mock}
	return mock
}

// EXPECT returns an object that allows you to define an expectation on the
// <.Name> mock client.
func (m *MockClient) EXPECT() *_MockClientRecorder {
	return m.recorder
}

<range .AllFunctions>
<$context := import "context">
<$yarpc   := import "go.uber.org/yarpc">

// <.Name> responds to a <.Name> call based on the mock expectations. This
// call will fail if the mock does not expect this call. Use EXPECT to expect
// a call to this function.
//
// 	client.EXPECT().<.Name>(gomock.Any(), ...).Return(...)
// 	... := client.<.Name>(...)
func (m *MockClient) <.Name>(
	ctx <$context>.Context, <range .Arguments>
	_<.Name> <formatType .Type>,<end>
	opts ...<$yarpc>.CallOption,
) <if .OneWay> (ack <$yarpc>.Ack, err error) {
  <else>       (<if .ReturnType>success <formatType .ReturnType>,<end> err error) {
  <end>
	args := []interface{}{ctx,<range .Arguments> _<.Name>,<end>}
	for _, o := range opts {
		args = append(args, o)
	}
	i := 0
	ret := m.ctrl.Call(m, "<.Name>", args...)
	<if .OneWay>          ack,     _ = ret[i].(<$yarpc>.Ack); i++
	<else if .ReturnType> success, _ = ret[i].(<formatType .ReturnType>); i++
	<end>                 err,     _ = ret[i].(error)
	return
}

func (mr *_MockClientRecorder) <.Name>(
	ctx interface{}, <range .Arguments>
	_<.Name> interface{},<end>
	opts ...interface{},
) *gomock.Call {
	args := append([]interface{}{ctx,<range .Arguments> _<.Name>,<end>}, opts...)
	return mr.mock.ctrl.RecordCall(mr.mock, "<.Name>", args...)
}
<end>
`

func gomockGenerator(data *serviceTemplateData, files map[string][]byte) (err error) {
	packageName := filepath.Base(data.TestPackagePath())
	// kv.thrift => .../kv/keyvaluetest/client.go
	path := filepath.Join(data.Module.Directory, packageName, "client.go")
	files[path], err = plugin.GoFileFromTemplate(path, gomockTemplate, data, templateOptions...)
	return
}
