// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/tianrandailove/peitho/pkg/k8s (interfaces: K8sService)

// Package k8s is a generated GoMock package.
package k8s

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockK8sService is a mock of K8sService interface.
type MockK8sService struct {
	ctrl     *gomock.Controller
	recorder *MockK8sServiceMockRecorder
}

func (m *MockK8sService) CreateChaincodeDeploymentWithPuller(ctx context.Context, name string, image string, env []string, cmd []string, pullerImag string, pullerCMD []string) error {
	return nil
}

// MockK8sServiceMockRecorder is the mock recorder for MockK8sService.
type MockK8sServiceMockRecorder struct {
	mock *MockK8sService
}

// NewMockK8sService creates a new mock instance.
func NewMockK8sService(ctrl *gomock.Controller) *MockK8sService {
	mock := &MockK8sService{ctrl: ctrl}
	mock.recorder = &MockK8sServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockK8sService) EXPECT() *MockK8sServiceMockRecorder {
	return m.recorder
}

// CreateChaincodeDeployment mocks base method.
func (m *MockK8sService) CreateChaincodeDeployment(arg0 context.Context, arg1, arg2 string, arg3, arg4 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChaincodeDeployment", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateChaincodeDeployment indicates an expected call of CreateChaincodeDeployment.
func (mr *MockK8sServiceMockRecorder) CreateChaincodeDeployment(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChaincodeDeployment", reflect.TypeOf((*MockK8sService)(nil).CreateChaincodeDeployment), arg0, arg1, arg2, arg3, arg4)
}

// CreateConfigMap mocks base method.
func (m *MockK8sService) CreateConfigMap(arg0 context.Context, arg1 string, arg2 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateConfigMap", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateConfigMap indicates an expected call of CreateConfigMap.
func (mr *MockK8sServiceMockRecorder) CreateConfigMap(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateConfigMap", reflect.TypeOf((*MockK8sService)(nil).CreateConfigMap), arg0, arg1, arg2)
}

// DeleteChaincodeDeployment mocks base method.
func (m *MockK8sService) DeleteChaincodeDeployment(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChaincodeDeployment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChaincodeDeployment indicates an expected call of DeleteChaincodeDeployment.
func (mr *MockK8sServiceMockRecorder) DeleteChaincodeDeployment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChaincodeDeployment", reflect.TypeOf((*MockK8sService)(nil).DeleteChaincodeDeployment), arg0, arg1)
}

// DeleteConfigMapDeployment mocks base method.
func (m *MockK8sService) DeleteConfigMapDeployment(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConfigMapDeployment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConfigMapDeployment indicates an expected call of DeleteConfigMapDeployment.
func (mr *MockK8sServiceMockRecorder) DeleteConfigMapDeployment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConfigMapDeployment", reflect.TypeOf((*MockK8sService)(nil).DeleteConfigMapDeployment), arg0, arg1)
}

// QueryDeploymentStatus mocks base method.
func (m *MockK8sService) QueryDeploymentStatus(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryDeploymentStatus", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryDeploymentStatus indicates an expected call of QueryDeploymentStatus.
func (mr *MockK8sServiceMockRecorder) QueryDeploymentStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryDeploymentStatus", reflect.TypeOf((*MockK8sService)(nil).QueryDeploymentStatus), arg0, arg1)
}

// UpdateDeployment mocks base method.
func (m *MockK8sService) UpdateDeployment(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDeployment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDeployment indicates an expected call of UpdateDeployment.
func (mr *MockK8sServiceMockRecorder) UpdateDeployment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeployment", reflect.TypeOf((*MockK8sService)(nil).UpdateDeployment), arg0, arg1)
}
func (mr *MockK8sServiceMockRecorder) CreateChaincodeDeploymentWithPuller(
	ctx context.Context,
	name string,
	image string,
	env []string,
	cmd []string,
	pullerImag string,
	pullerCMD []string,
) error {
	return nil
}
