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

func Test_CreateGroup(t *testing.T) {
	type createGroupResponse struct {
		StatusCode int
		Content    []byte
	}
	type args struct {
		name       string
		roles      []string
		attributes ak.GroupAttributes
	}
	tests := []struct {
		name                string
		createGroupResponse createGroupResponse
		args                args
		mockAkRepo          func(mockIdpRepo *mock_ak.MockAuthentikRepository)
		want                *ak.Group
		wantErr             bool
	}{
		{
			name: "create group error",
			args: args{
				name:  "example-group",
				roles: []string{"example-role"},
				attributes: ak.GroupAttributes{
					Tenant: "example-tenant",
				},
			},
			wantErr: true,
			createGroupResponse: createGroupResponse{
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
				mockAkRepo.EXPECT().CreateGroup("example-group",
					[]string{"example-role"},
					groupAttributes{
						Tenant: "example-tenant",
					},
				).Return(nil, errors.New("{\"details\":\"no permissions\",\"code\":\"403\""))
			},
		},
		{
			name: "create group success",
			args: args{
				name:  "example-group",
				roles: []string{"example-role"},
				attributes: ak.GroupAttributes{
					Tenant: "example-tenant",
				},
			},
			wantErr: false,
			want: &ak.Group{
				PK:   "random-uuid",
				Name: "example-group",
			},
			createGroupResponse: createGroupResponse{
				StatusCode: http.StatusCreated,
				Content: (func() []byte {
					marshal, _ := json.Marshal(getGroupResponse{groupObj{
						PK:         "random-uuid",
						Name:       "example-group",
						Attributes: groupAttributes{},
					}})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().CreateGroup("example-group",
					[]string{"example-role"},
					groupAttributes{
						Tenant: "example-tenant",
					}).Return(ak.Group{
					PK:   "random-uuid",
					Name: "example-group",
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.createGroupResponse.StatusCode)
				_, _ = w.Write(tt.createGroupResponse.Content)
			})
			testServer := httptest.NewServer(h)
			defer testServer.Close()

			r := New(testServer.URL, "test")

			got, err := r.CreateGroup(tt.args.name, tt.args.roles, tt.args.attributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateUser(t *testing.T) {
	type createUserResponse struct {
		StatusCode int
		Content    []byte
	}
	type args struct {
		user ak.User
	}
	tests := []struct {
		name               string
		createUserResponse createUserResponse
		args               args
		mockAkRepo         func(mockIdpRepo *mock_ak.MockAuthentikRepository)
		want               *ak.User
		wantErr            bool
	}{
		{
			name: "create user token error",
			args: args{
				user: ak.User{
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@example.com",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				},
			},
			wantErr: true,
			createUserResponse: createUserResponse{
				StatusCode: http.StatusForbidden,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Details string `json:"detail"`
						Code    string `json:"code"`
					}{
						Details: "no permissions",
						Code:    "403",
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().CreateUser(ak.User{
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@example.com",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				}).Return(nil, errors.New("{\"details\":\"no permissions\",\"code\":\"403\""))
			},
		},
		{
			name: "create user payload error",
			args: args{
				user: ak.User{
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@invalidDomain",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				},
			},
			wantErr: true,
			createUserResponse: createUserResponse{
				StatusCode: http.StatusBadRequest,
				Content: (func() []byte {
					marshal, _ := json.Marshal(struct {
						Details string `json:"details"`
						Code    string `json:"code"`
					}{
						Details: "invalid email",
						Code:    "400",
					})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().CreateUser(ak.User{
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@invalidDomain",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				}).Return(nil, errors.New("{\"details\":\"invalid email\",\"code\":\"400\""))
			},
		},
		{
			name: "create user success",
			args: args{
				user: ak.User{
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@example.com",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				},
			},
			wantErr: false,
			want: &ak.User{
				PK:         999,
				Username:   "johndoe",
				Name:       "John Doe",
				Email:      "John.Doe@example.com",
				Path:       "users",
				IsActive:   true,
				Attributes: ak.UserAttributes{},
			},
			createUserResponse: createUserResponse{
				StatusCode: http.StatusCreated,
				Content: (func() []byte {
					marshal, _ := json.Marshal(getUserResponse{userObj{
						PK:         999,
						Username:   "johndoe",
						Name:       "John Doe",
						Email:      "John.Doe@example.com",
						Path:       "users",
						IsActive:   true,
						Attributes: userAttributes{},
					}})
					return marshal
				})(),
			},
			mockAkRepo: func(mockAkRepo *mock_ak.MockAuthentikRepository) {
				mockAkRepo.EXPECT().CreateUser(ak.User{
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@example.com",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				}).Return(ak.User{
					PK:         999,
					Username:   "johndoe",
					Name:       "John Doe",
					Email:      "John.Doe@example.com",
					Path:       "users",
					IsActive:   true,
					Attributes: ak.UserAttributes{},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.createUserResponse.StatusCode)
				_, _ = w.Write(tt.createUserResponse.Content)
			})
			testServer := httptest.NewServer(h)
			defer testServer.Close()

			r := New(testServer.URL, "test")

			got, err := r.CreateUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
