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
