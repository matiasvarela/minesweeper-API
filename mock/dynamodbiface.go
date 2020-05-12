// Code generated by MockGen. DO NOT EDIT.
// Source: dynamodbiface.go

// Package mock is a generated GoMock package.
package mock

import (
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDynamoDB is a mock of DynamoDB interface
type MockDynamoDB struct {
	ctrl     *gomock.Controller
	recorder *MockDynamoDBMockRecorder
}

// MockDynamoDBMockRecorder is the mock recorder for MockDynamoDB
type MockDynamoDBMockRecorder struct {
	mock *MockDynamoDB
}

// NewMockDynamoDB creates a new mock instance
func NewMockDynamoDB(ctrl *gomock.Controller) *MockDynamoDB {
	mock := &MockDynamoDB{ctrl: ctrl}
	mock.recorder = &MockDynamoDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDynamoDB) EXPECT() *MockDynamoDBMockRecorder {
	return m.recorder
}

// Query mocks base method
func (m *MockDynamoDB) Query(arg0 *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(*dynamodb.QueryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockDynamoDBMockRecorder) Query(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDynamoDB)(nil).Query), arg0)
}

// GetItem mocks base method
func (m *MockDynamoDB) GetItem(arg0 *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", arg0)
	ret0, _ := ret[0].(*dynamodb.GetItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem
func (mr *MockDynamoDBMockRecorder) GetItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockDynamoDB)(nil).GetItem), arg0)
}

// PutItem mocks base method
func (m *MockDynamoDB) PutItem(arg0 *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutItem", arg0)
	ret0, _ := ret[0].(*dynamodb.PutItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutItem indicates an expected call of PutItem
func (mr *MockDynamoDBMockRecorder) PutItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItem", reflect.TypeOf((*MockDynamoDB)(nil).PutItem), arg0)
}
