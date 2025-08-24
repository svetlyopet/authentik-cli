package authentik

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	mock_ak "github.com/svetlyopet/authentik-cli/mocks/ak"
)

func Test_CreateRole(t *testing.T) {
	type createRoleResponse struct {
		StatusCode int
		Content    []byte
	}
	type args struct {
		name string
	}
	tests := []struct {
		name               string
		createRoleResponse createRoleResponse
		args               args
		mockAkRepo         func(mockIdpRepo *mock_ak.MockAuthentikRepository)
		want               *ak.Role
		wantErr            bool
	}{
		{
			name: "create role error",
			args: args{
				name: "example-role",
			},
			wantErr: true,
			createRoleResponse: createRoleResponse{
				StatusCode: http.StatusBadRequest,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Details string `json:"detail"`
						Code    string `json:"code"`
					}{
						Details: "bad request",
						Code:    "400",
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().CreateRole("example-role").
					Return(nil, errors.New("{\"details\":\"bad request\",\"code\":\"400\""))
			},
		},
		{
			name: "create role success",
			args: args{
				name: "example-role",
			},
			wantErr: false,
			want: &ak.Role{
				PK:   "random-uuid",
				Name: "example-role",
			},
			createRoleResponse: createRoleResponse{
				StatusCode: http.StatusCreated,
				Content: (func() []byte {
					marshal, _ := json.Marshal(getRoleResponse{roleObj{
						PK:   "random-uuid",
						Name: "example-role",
					}})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().CreateRole("example-role").Return(ak.Role{
					PK:   "random-uuid",
					Name: "example-role",
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.createRoleResponse.StatusCode)
				_, _ = w.Write(tt.createRoleResponse.Content)
			})
			testServer := httptest.NewServer(h)
			defer testServer.Close()

			r := New(testServer.URL, "test")

			got, err := r.CreateRole(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetRoleByName(t *testing.T) {
	type getRoleByNameResponse struct {
		StatusCode int
		Content    []byte
	}
	type args struct {
		name string
	}
	tests := []struct {
		name                  string
		getRoleByNameResponse getRoleByNameResponse
		args                  args
		mockAkRepo            func(mockIdpRepo *mock_ak.MockAuthentikRepository)
		want                  *ak.Role
		wantErr               bool
	}{
		{
			name: "role not found",
			args: args{
				name: "missing-role",
			},
			wantErr: true, // <-- should be true!
			getRoleByNameResponse: getRoleByNameResponse{
				StatusCode: http.StatusOK,
				Content:    []byte(`{"results":[]}`),
			},
		},
		{
			name: "get role by name error",
			args: args{
				name: "example-role",
			},
			wantErr: true,
			getRoleByNameResponse: getRoleByNameResponse{
				StatusCode: http.StatusBadRequest,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Details string `json:"detail"`
						Code    string `json:"code"`
					}{
						Details: "bad request",
						Code:    "400",
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().GetRoleByName("example-role").
					Return(nil, errors.New("{\"details\":\"bad request\",\"code\":\"400\""))
			},
		},
		{
			name: "multiple roles found",
			args: args{
				name: "duplicate-role",
			},
			wantErr: true,
			getRoleByNameResponse: getRoleByNameResponse{
				StatusCode: http.StatusOK,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Results []roleObj `json:"results"`
					}{
						Results: []roleObj{
							{
								PK:   "uuid-1",
								Name: "duplicate-role",
							},
							{
								PK:   "uuid-2",
								Name: "duplicate-role",
							},
						},
					})
					return marshal
				})(),
			},
		},
		{
			name: "get role by name success",
			args: args{
				name: "example-role",
			},
			wantErr: false,
			want: &ak.Role{
				PK:   "random-uuid",
				Name: "example-role",
			},
			getRoleByNameResponse: getRoleByNameResponse{
				StatusCode: http.StatusOK,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Results []roleObj `json:"results"`
					}{
						Results: []roleObj{
							{
								PK:   "random-uuid",
								Name: "example-role",
							},
						},
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().GetRoleByName("example-role").Return(&ak.Role{
					PK:   "random-uuid",
					Name: "example-role",
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.getRoleByNameResponse.StatusCode)
				_, _ = w.Write(tt.getRoleByNameResponse.Content)
			})
			testServer := httptest.NewServer(h)
			defer testServer.Close()

			r := New(testServer.URL, "test")

			got, err := r.GetRoleByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoleByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoleByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DeleteRole(t *testing.T) {
	type deleteRoleResponse struct {
		StatusCode int
		Content    []byte
	}
	type args struct {
		id string
	}
	tests := []struct {
		name               string
		deleteRoleResponse deleteRoleResponse
		args               args
		mockAkRepo         func(mockIdpRepo *mock_ak.MockAuthentikRepository)
		wantErr            bool
	}{
		{
			name: "delete role error",
			args: args{
				id: "123",
			},
			wantErr: true,
			deleteRoleResponse: deleteRoleResponse{
				StatusCode: http.StatusBadRequest,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Details string `json:"detail"`
						Code    string `json:"code"`
					}{
						Details: "bad request",
						Code:    "400",
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().DeleteRole("123").
					Return(errors.New("{\"details\":\"bad request\",\"code\":\"400\""))
			},
		},
		{
			name: "delete role success",
			args: args{
				id: "123",
			},
			wantErr: false,
			deleteRoleResponse: deleteRoleResponse{
				StatusCode: http.StatusNoContent,
				Content:    []byte{},
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().DeleteRole("123").Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.deleteRoleResponse.StatusCode)
				_, _ = w.Write(tt.deleteRoleResponse.Content)
			})
			testServer := httptest.NewServer(h)
			defer testServer.Close()

			r := New(testServer.URL, "test")

			err := r.DeleteRole(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_AssignViewPermissionsToTenantRole(t *testing.T) {
	type assignPermissionsResponse struct {
		StatusCode int
		Content    []byte
	}
	type args struct {
		rolePK string
	}
	tests := []struct {
		name                      string
		assignPermissionsResponse assignPermissionsResponse
		args                      args
		mockAkRepo                func(mockIdpRepo *mock_ak.MockAuthentikRepository)
		wantErr                   bool
	}{
		{
			name: "assign permissions error",
			args: args{
				rolePK: "random-uuid",
			},
			wantErr: true,
			assignPermissionsResponse: assignPermissionsResponse{
				StatusCode: http.StatusBadRequest,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Details string `json:"detail"`
						Code    string `json:"code"`
					}{
						Details: "bad request",
						Code:    "400",
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().AssignViewPermissionsToTenantRole("random-uuid").
					Return(errors.New("{\"details\":\"bad request\",\"code\":\"400\""))
			},
		},
		{
			name: "assign permissions success",
			args: args{
				rolePK: "random-uuid",
			},
			wantErr: false,
			assignPermissionsResponse: assignPermissionsResponse{
				StatusCode: http.StatusOK,
				Content:    []byte{},
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().AssignViewPermissionsToTenantRole("random-uuid").Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.assignPermissionsResponse.StatusCode)
				_, _ = w.Write(tt.assignPermissionsResponse.Content)
			})
			testServer := httptest.NewServer(h)
			defer testServer.Close()

			r := New(testServer.URL, "test")

			err := r.AssignViewPermissionsToTenantRole(tt.args.rolePK)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignViewPermissionsToTenantRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
