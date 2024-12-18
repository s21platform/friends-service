// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDBRepo is a mock of DBRepo interface.
type MockDBRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDBRepoMockRecorder
}

// MockDBRepoMockRecorder is the mock recorder for MockDBRepo.
type MockDBRepoMockRecorder struct {
	mock *MockDBRepo
}

// NewMockDBRepo creates a new mock instance.
func NewMockDBRepo(ctrl *gomock.Controller) *MockDBRepo {
	mock := &MockDBRepo{ctrl: ctrl}
	mock.recorder = &MockDBRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBRepo) EXPECT() *MockDBRepoMockRecorder {
	return m.recorder
}

// GetCountFriends mocks base method.
func (m *MockDBRepo) GetCountFriends(uuid string) (int64, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountFriends", uuid)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCountFriends indicates an expected call of GetCountFriends.
func (mr *MockDBRepoMockRecorder) GetCountFriends(uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountFriends", reflect.TypeOf((*MockDBRepo)(nil).GetCountFriends), uuid)
}

// GetPeerFollows mocks base method.
func (m *MockDBRepo) GetPeerFollows(initiator string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeerFollows", initiator)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeerFollows indicates an expected call of GetPeerFollows.
func (mr *MockDBRepoMockRecorder) GetPeerFollows(initiator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeerFollows", reflect.TypeOf((*MockDBRepo)(nil).GetPeerFollows), initiator)
}

// GetWhoFollowsPeer mocks base method.
func (m *MockDBRepo) GetWhoFollowsPeer(initiator string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWhoFollowsPeer", initiator)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWhoFollowsPeer indicates an expected call of GetWhoFollowsPeer.
func (mr *MockDBRepoMockRecorder) GetWhoFollowsPeer(initiator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWhoFollowsPeer", reflect.TypeOf((*MockDBRepo)(nil).GetWhoFollowsPeer), initiator)
}

// IsRowFriendExist mocks base method.
func (m *MockDBRepo) IsRowFriendExist(peer1, peer2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRowFriendExist", peer1, peer2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRowFriendExist indicates an expected call of IsRowFriendExist.
func (mr *MockDBRepoMockRecorder) IsRowFriendExist(peer1, peer2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRowFriendExist", reflect.TypeOf((*MockDBRepo)(nil).IsRowFriendExist), peer1, peer2)
}

// RemoveFriends mocks base method.
func (m *MockDBRepo) RemoveFriends(peer1, peer2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFriends", peer1, peer2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveFriends indicates an expected call of RemoveFriends.
func (mr *MockDBRepoMockRecorder) RemoveFriends(peer1, peer2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFriends", reflect.TypeOf((*MockDBRepo)(nil).RemoveFriends), peer1, peer2)
}

// RemoveSubscribe mocks base method.
func (m *MockDBRepo) RemoveSubscribe(peer1, peer2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSubscribe", peer1, peer2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSubscribe indicates an expected call of RemoveSubscribe.
func (mr *MockDBRepoMockRecorder) RemoveSubscribe(peer1, peer2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSubscribe", reflect.TypeOf((*MockDBRepo)(nil).RemoveSubscribe), peer1, peer2)
}

// SetFriend mocks base method.
func (m *MockDBRepo) SetFriend(peer1, peer2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFriend", peer1, peer2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetFriend indicates an expected call of SetFriend.
func (mr *MockDBRepoMockRecorder) SetFriend(peer1, peer2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFriend", reflect.TypeOf((*MockDBRepo)(nil).SetFriend), peer1, peer2)
}

// SetInvitePeer mocks base method.
func (m *MockDBRepo) SetInvitePeer(initiator, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInvitePeer", initiator, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInvitePeer indicates an expected call of SetInvitePeer.
func (mr *MockDBRepoMockRecorder) SetInvitePeer(initiator, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInvitePeer", reflect.TypeOf((*MockDBRepo)(nil).SetInvitePeer), initiator, email)
}
