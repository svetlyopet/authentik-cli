package core

import (
	"errors"
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/diff"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateUser(username, name, surname, email, tenant string) (err error) {
	userPath := "users"
	userAttr := ak.UserAttributes{}

	if surname != "" {
		name = fmt.Sprintf("%s %s", name, surname)
	}

	tenantGroup := &ak.Group{}
	if tenant != "" {
		tenantGroup, err = GetGroupByName(fmt.Sprintf(constants.TenantAdminGroupNamePattern, tenant))
		if err != nil {
			return customErrors.NewNotExists("tenant not found")
		}
		userPath = tenant
		userAttr.UserType = "tenant-admin"
		userAttr.Tenant = tenant
	} else {
		userAttr.UserType = "global"
		userAttr.Tenant = "global"
	}

	usr := ak.User{
		Username:    username,
		Name:        name,
		Email:       email,
		Path:        userPath,
		IsActive:    true,
		IsSuperuser: false,
		Attributes:  userAttr,
	}

	user, err := ak.Repo.CreateUser(usr)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeUser, constants.ActionCreated, username)

	if tenant != "" {
		err = AddUserToGroup(user.PK, tenantGroup.PK, tenantGroup.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetUserDetails(username string) (*User, error) {
	user, err := ak.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return mapToGetUserDetails(user), nil
}

func DeleteUser(username string) (err error) {
	user, err := ak.Repo.GetUserByUsername(username)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
	}

	err = ak.Repo.DeleteUser(fmt.Sprintf("%d", user.PK))
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeUser, constants.ActionDeleted, username)

	return nil
}

// UserChanges lists the changes requested for a user. A nil field is left
// unchanged.
type UserChanges struct {
	Username *string
	Name     *string
	// Surname is appended to Name to form the display name, as on creation,
	// so it can only be set together with Name.
	Surname  *string
	Email    *string
	IsActive *bool
}

// UpdateUser applies the requested changes to the user with the given
// username. Fields that were not requested, or already hold the requested
// value, are not sent to authentik and keep their current value.
func UpdateUser(username string, changes UserChanges) error {
	if changes.Surname != nil && changes.Name == nil {
		return errors.New("a surname can only be set together with a name, as both form the display name")
	}

	user, err := ak.Repo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	var name *string
	if changes.Name != nil {
		fullName := *changes.Name
		if changes.Surname != nil && *changes.Surname != "" {
			fullName = fmt.Sprintf("%s %s", fullName, *changes.Surname)
		}
		name = &fullName
	}

	patch := ak.UserPatch{
		Username: diff.Changed(user.Username, changes.Username),
		Name:     diff.Changed(user.Name, name),
		Email:    diff.Changed(user.Email, changes.Email),
		IsActive: diff.Changed(user.IsActive, changes.IsActive),
	}
	if patch.IsEmpty() {
		logger.LogObjectChange(constants.ObjectTypeUser, constants.ActionUnchanged, username)
		return nil
	}

	updated, err := ak.Repo.UpdateUser(user.PK, patch)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeUser, constants.ActionChanged, updated.Username)

	return nil
}
