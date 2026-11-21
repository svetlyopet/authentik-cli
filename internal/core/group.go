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

func CreateGroup(name string, roles []string, attributes ak.GroupAttributes) (*ak.Group, error) {
	group, err := ak.Repo.CreateGroup(name, roles, attributes)
	if err != nil {
		return nil, err
	}

	logger.LogObjectChange(constants.ObjectTypeGroup, constants.ActionCreated, name)

	return group, nil
}

func GetGroupDetails(name string) (*Group, error) {
	groupPreflight, err := ak.Repo.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	group, err := ak.Repo.GetGroup(groupPreflight.PK)
	if err != nil {
		return nil, err
	}

	return mapToGetGroupDetails(group), nil
}

func GetGroupByName(name string) (group *ak.Group, err error) {
	group, err = ak.Repo.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func DeleteGroup(name string) (err error) {
	group, err := GetGroupByName(name)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if !errors.As(err, &notExistsError) {
			return err
		}
		return nil
	}

	err = ak.Repo.DeleteGroup(group.PK)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeGroup, constants.ActionDeleted, name)

	return nil
}

func AddUserToGroup(userPK int, groupPK, groupName string) (err error) {
	err = ak.Repo.AddUserToGroup(userPK, groupPK)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeGroup, constants.ActionChanged, groupName)

	return nil
}

// GroupChanges lists the changes requested for a group. A nil field is left
// unchanged.
type GroupChanges struct {
	Name *string
}

// UpdateGroup applies the requested changes to the group with the given name.
// Fields that were not requested, or already hold the requested value, are
// not sent to authentik and keep their current value.
func UpdateGroup(name string, changes GroupChanges) error {
	group, err := GetGroupByName(name)
	if err != nil {
		return err
	}

	patch := ak.GroupPatch{
		Name: diff.Changed(group.Name, changes.Name),
	}

	// Tenants find their admin group by name, so renaming it would orphan the tenant.
	if patch.Name != nil && group.GroupAttributes.Tenant != "" {
		return fmt.Errorf("group %s is managed by tenant %s and cannot be renamed", name, group.GroupAttributes.Tenant)
	}

	if patch.IsEmpty() {
		logger.LogObjectChange(constants.ObjectTypeGroup, constants.ActionUnchanged, name)
		return nil
	}

	updated, err := ak.Repo.UpdateGroup(group.PK, patch)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeGroup, constants.ActionChanged, updated.Name)

	return nil
}
