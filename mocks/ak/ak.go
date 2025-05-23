// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ak/ak.go
//
// Generated by this command:
//
//	mockgen -destination=mocks/ak/ak.go -package=mock_ak -source=internal/ak/ak.go
//

// Package mock_ak is a generated GoMock package.
package mock_ak

import (
	reflect "reflect"

	ak "github.com/svetlyopet/authentik-cli/internal/ak"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthentikRepository is a mock of AuthentikRepository interface.
type MockAuthentikRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthentikRepositoryMockRecorder
	isgomock struct{}
}

// MockAuthentikRepositoryMockRecorder is the mock recorder for MockAuthentikRepository.
type MockAuthentikRepositoryMockRecorder struct {
	mock *MockAuthentikRepository
}

// NewMockAuthentikRepository creates a new mock instance.
func NewMockAuthentikRepository(ctrl *gomock.Controller) *MockAuthentikRepository {
	mock := &MockAuthentikRepository{ctrl: ctrl}
	mock.recorder = &MockAuthentikRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthentikRepository) EXPECT() *MockAuthentikRepositoryMockRecorder {
	return m.recorder
}

// AddUserToGroup mocks base method.
func (m *MockAuthentikRepository) AddUserToGroup(userPK int, uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToGroup", userPK, uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToGroup indicates an expected call of AddUserToGroup.
func (mr *MockAuthentikRepositoryMockRecorder) AddUserToGroup(userPK, uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToGroup", reflect.TypeOf((*MockAuthentikRepository)(nil).AddUserToGroup), userPK, uuid)
}

// AssignViewPermissionsToTenantRole mocks base method.
func (m *MockAuthentikRepository) AssignViewPermissionsToTenantRole(rolePK string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignViewPermissionsToTenantRole", rolePK)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignViewPermissionsToTenantRole indicates an expected call of AssignViewPermissionsToTenantRole.
func (mr *MockAuthentikRepositoryMockRecorder) AssignViewPermissionsToTenantRole(rolePK any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignViewPermissionsToTenantRole", reflect.TypeOf((*MockAuthentikRepository)(nil).AssignViewPermissionsToTenantRole), rolePK)
}

// CreateApplication mocks base method.
func (m *MockAuthentikRepository) CreateApplication(name, slug string, providerPK int) (*ak.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApplication", name, slug, providerPK)
	ret0, _ := ret[0].(*ak.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateApplication indicates an expected call of CreateApplication.
func (mr *MockAuthentikRepositoryMockRecorder) CreateApplication(name, slug, providerPK any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApplication", reflect.TypeOf((*MockAuthentikRepository)(nil).CreateApplication), name, slug, providerPK)
}

// CreateGroup mocks base method.
func (m *MockAuthentikRepository) CreateGroup(name string, roles []string, attributes ak.GroupAttributes) (*ak.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", name, roles, attributes)
	ret0, _ := ret[0].(*ak.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockAuthentikRepositoryMockRecorder) CreateGroup(name, roles, attributes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockAuthentikRepository)(nil).CreateGroup), name, roles, attributes)
}

// CreateOidcProvider mocks base method.
func (m *MockAuthentikRepository) CreateOidcProvider(provider ak.OidcProvider) (*ak.OidcProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOidcProvider", provider)
	ret0, _ := ret[0].(*ak.OidcProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOidcProvider indicates an expected call of CreateOidcProvider.
func (mr *MockAuthentikRepositoryMockRecorder) CreateOidcProvider(provider any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOidcProvider", reflect.TypeOf((*MockAuthentikRepository)(nil).CreateOidcProvider), provider)
}

// CreateRole mocks base method.
func (m *MockAuthentikRepository) CreateRole(name string) (*ak.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", name)
	ret0, _ := ret[0].(*ak.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockAuthentikRepositoryMockRecorder) CreateRole(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockAuthentikRepository)(nil).CreateRole), name)
}

// CreateUser mocks base method.
func (m *MockAuthentikRepository) CreateUser(usr ak.User) (*ak.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", usr)
	ret0, _ := ret[0].(*ak.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthentikRepositoryMockRecorder) CreateUser(usr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthentikRepository)(nil).CreateUser), usr)
}

// DeleteApplication mocks base method.
func (m *MockAuthentikRepository) DeleteApplication(slug string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteApplication", slug)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApplication indicates an expected call of DeleteApplication.
func (mr *MockAuthentikRepositoryMockRecorder) DeleteApplication(slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApplication", reflect.TypeOf((*MockAuthentikRepository)(nil).DeleteApplication), slug)
}

// DeleteGroup mocks base method.
func (m *MockAuthentikRepository) DeleteGroup(uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockAuthentikRepositoryMockRecorder) DeleteGroup(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockAuthentikRepository)(nil).DeleteGroup), uuid)
}

// DeleteProvider mocks base method.
func (m *MockAuthentikRepository) DeleteProvider(pk int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProvider", pk)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProvider indicates an expected call of DeleteProvider.
func (mr *MockAuthentikRepositoryMockRecorder) DeleteProvider(pk any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProvider", reflect.TypeOf((*MockAuthentikRepository)(nil).DeleteProvider), pk)
}

// DeleteRole mocks base method.
func (m *MockAuthentikRepository) DeleteRole(uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRole", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRole indicates an expected call of DeleteRole.
func (mr *MockAuthentikRepositoryMockRecorder) DeleteRole(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRole", reflect.TypeOf((*MockAuthentikRepository)(nil).DeleteRole), uuid)
}

// DeleteUser mocks base method.
func (m *MockAuthentikRepository) DeleteUser(userPK string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userPK)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockAuthentikRepositoryMockRecorder) DeleteUser(userPK any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockAuthentikRepository)(nil).DeleteUser), userPK)
}

// GetApplicationByName mocks base method.
func (m *MockAuthentikRepository) GetApplicationByName(name string) (*ak.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationByName", name)
	ret0, _ := ret[0].(*ak.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationByName indicates an expected call of GetApplicationByName.
func (mr *MockAuthentikRepositoryMockRecorder) GetApplicationByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationByName", reflect.TypeOf((*MockAuthentikRepository)(nil).GetApplicationByName), name)
}

// GetAuthentikTargetUrl mocks base method.
func (m *MockAuthentikRepository) GetAuthentikTargetUrl() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthentikTargetUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAuthentikTargetUrl indicates an expected call of GetAuthentikTargetUrl.
func (mr *MockAuthentikRepositoryMockRecorder) GetAuthentikTargetUrl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthentikTargetUrl", reflect.TypeOf((*MockAuthentikRepository)(nil).GetAuthentikTargetUrl))
}

// GetFlows mocks base method.
func (m *MockAuthentikRepository) GetFlows() ([]ak.Flow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlows")
	ret0, _ := ret[0].([]ak.Flow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlows indicates an expected call of GetFlows.
func (mr *MockAuthentikRepositoryMockRecorder) GetFlows() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlows", reflect.TypeOf((*MockAuthentikRepository)(nil).GetFlows))
}

// GetGroup mocks base method.
func (m *MockAuthentikRepository) GetGroup(uuid string) (*ak.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", uuid)
	ret0, _ := ret[0].(*ak.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockAuthentikRepositoryMockRecorder) GetGroup(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockAuthentikRepository)(nil).GetGroup), uuid)
}

// GetGroupByName mocks base method.
func (m *MockAuthentikRepository) GetGroupByName(name string) (*ak.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByName", name)
	ret0, _ := ret[0].(*ak.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByName indicates an expected call of GetGroupByName.
func (mr *MockAuthentikRepositoryMockRecorder) GetGroupByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByName", reflect.TypeOf((*MockAuthentikRepository)(nil).GetGroupByName), name)
}

// GetOidcProvider mocks base method.
func (m *MockAuthentikRepository) GetOidcProvider(pk int) (*ak.OidcProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOidcProvider", pk)
	ret0, _ := ret[0].(*ak.OidcProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOidcProvider indicates an expected call of GetOidcProvider.
func (mr *MockAuthentikRepositoryMockRecorder) GetOidcProvider(pk any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOidcProvider", reflect.TypeOf((*MockAuthentikRepository)(nil).GetOidcProvider), pk)
}

// GetRoleByName mocks base method.
func (m *MockAuthentikRepository) GetRoleByName(name string) (*ak.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoleByName", name)
	ret0, _ := ret[0].(*ak.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleByName indicates an expected call of GetRoleByName.
func (mr *MockAuthentikRepositoryMockRecorder) GetRoleByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleByName", reflect.TypeOf((*MockAuthentikRepository)(nil).GetRoleByName), name)
}

// GetUserByUsername mocks base method.
func (m *MockAuthentikRepository) GetUserByUsername(username string) (*ak.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", username)
	ret0, _ := ret[0].(*ak.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockAuthentikRepositoryMockRecorder) GetUserByUsername(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockAuthentikRepository)(nil).GetUserByUsername), username)
}
