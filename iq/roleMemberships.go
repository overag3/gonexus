package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// Before 70
	restRoleMembersOrgDeprecated = "api/v2/organizations/%s/roleMembers"
	restRoleMembersAppDeprecated = "api/v2/applications/%s/roleMembers"

	// After 70
	restRoleMembersOrgGet    = "api/v2/roleMemberships/organization/%s"
	restRoleMembersAppGet    = "api/v2/roleMemberships/application/%s"
	restRoleMembersReposGet  = "api/v2/roleMemberships/repository_container"
	restRoleMembersGlobalGet = "api/v2/roleMemberships/global"

	restRoleMembersOrgUser         = "api/v2/roleMemberships/organization/%s/role/%s/user/%s"
	restRoleMembersOrgGroup        = "api/v2/roleMemberships/organization/%s/role/%s/group/%s"
	restRoleMembersAppUser         = "api/v2/roleMemberships/application/%s/role/%s/user/%s"
	restRoleMembersAppGroup        = "api/v2/roleMemberships/application/%s/role/%s/group/%s"
	restRoleMembersRepositoryUser  = "api/v2/roleMemberships/repository_container/role/%s/user/%s"
	restRoleMembersRepositoryGroup = "api/v2/roleMemberships/repository_container/role/%s/group/%s"
	restRoleMembersGlobalUser      = "api/v2/roleMemberships/global/role/%s/user/%s"
	restRoleMembersGlobalGroup     = "api/v2/roleMemberships/global/role/%s/group/%s"
)

// Constants to describe a Member Type
const (
	MemberTypeUser  = "USER"
	MemberTypeGroup = "GROUP"
)

type memberMappings struct {
	MemberMappings []MemberMapping `json:"memberMappings"`
}

// MemberMapping describes a list of Members against a Role
type MemberMapping struct {
	RoleID  string   `json:"roleId"`
	Members []Member `json:"members"`
}

// Member describes a member to map with a role
type Member struct {
	OwnerID         string `json:"ownerId,omitempty"`
	OwnerType       string `json:"ownerType,omitempty"`
	Type            string `json:"type"`
	UserOrGroupName string `json:"userOrGroupName"`
}

func hasRev70API(ctx context.Context, iq IQ) bool {
	api := fmt.Sprintf(restRoleMembersOrgGet, RootOrganization)
	request, _ := iq.NewRequest(ctx, "HEAD", api, nil)
	_, resp, _ := iq.Do(request)
	return resp.StatusCode != http.StatusNotFound
}

func newMapping(roleID, memberType, memberName string) MemberMapping {
	return MemberMapping{
		RoleID: roleID,
		Members: []Member{
			{
				Type:            memberType,
				UserOrGroupName: memberName,
			},
		},
	}
}

func newMappings(roleID, memberType, memberName string) memberMappings {
	return memberMappings{
		MemberMappings: []MemberMapping{newMapping(roleID, memberType, memberName)},
	}
}

func organizationAuthorizationsByID(ctx context.Context, iq IQ, orgID string) ([]MemberMapping, error) {
	var endpoint string
	if hasRev70API(ctx, iq) {
		endpoint = fmt.Sprintf(restRoleMembersOrgGet, orgID)
	} else {
		endpoint = fmt.Sprintf(restRoleMembersOrgDeprecated, orgID)
	}

	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve role mapping for organization %s: %v", orgID, err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)

	return mappings.MemberMappings, err
}

func organizationAuthorizationsByRoleID(ctx context.Context, iq IQ, roleID string) ([]MemberMapping, error) {
	orgs, err := GetAllOrganizationsContext(ctx, iq)
	if err != nil {
		return nil, fmt.Errorf("could not find organizations: %v", err)
	}

	mappings := make([]MemberMapping, 0)
	for _, org := range orgs {
		orgMaps, _ := organizationAuthorizationsByID(ctx, iq, org.ID)
		for _, m := range orgMaps {
			if m.RoleID == roleID {
				mappings = append(mappings, m)
			}
		}
	}

	return mappings, nil
}

func OrganizationAuthorizationsContext(ctx context.Context, iq IQ, name string) ([]MemberMapping, error) {
	org, err := GetOrganizationByNameContext(ctx, iq, name)
	if err != nil {
		return nil, fmt.Errorf("could not find organization with name %s: %v", name, err)
	}

	return organizationAuthorizationsByID(ctx, iq, org.ID)
}

// OrganizationAuthorizations returns the member mappings of an organization
func OrganizationAuthorizations(iq IQ, name string) ([]MemberMapping, error) {
	return OrganizationAuthorizationsContext(context.Background(), iq, name)
}

func OrganizationAuthorizationsByRoleContext(ctx context.Context, iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	return organizationAuthorizationsByRoleID(ctx, iq, role.ID)
}

// OrganizationAuthorizationsByRole returns the member mappings of all organizations which match the given role
func OrganizationAuthorizationsByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	return OrganizationAuthorizationsByRoleContext(context.Background(), iq, roleName)
}

func setOrganizationAuth(ctx context.Context, iq IQ, name, roleName, member, memberType string) error {
	org, err := GetOrganizationByNameContext(ctx, iq, name)
	if err != nil {
		return fmt.Errorf("could not find organization with name %s: %v", name, err)
	}

	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	var payload io.Reader
	if hasRev70API(ctx, iq) {
		switch memberType {
		case MemberTypeUser:
			endpoint = fmt.Sprintf(restRoleMembersOrgUser, org.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersOrgGroup, org.ID, role.ID, member)
		}
	} else {
		endpoint = fmt.Sprintf(restRoleMembersOrgDeprecated, org.ID)
		current, err := OrganizationAuthorizationsContext(ctx, iq, name)
		if err != nil && current == nil {
			current = make([]MemberMapping, 0)
		}
		current = append(current, newMapping(role.ID, memberType, member))

		buf, err := json.Marshal(memberMappings{MemberMappings: current})
		if err != nil {
			return fmt.Errorf("could not create mapping: %v", err)
		}
		payload = bytes.NewBuffer(buf)
	}

	_, _, err = iq.Put(ctx, endpoint, payload)
	if err != nil {
		return fmt.Errorf("could not update organization role mapping: %v", err)
	}

	return nil
}

func SetOrganizationUserContext(ctx context.Context, iq IQ, name, roleName, user string) error {
	return setOrganizationAuth(ctx, iq, name, roleName, user, MemberTypeUser)
}

// SetOrganizationUser sets the role and user that can have access to an organization
func SetOrganizationUser(iq IQ, name, roleName, user string) error {
	return SetOrganizationUserContext(context.Background(), iq, name, roleName, user)
}

func SetOrganizationGroupContext(ctx context.Context, iq IQ, name, roleName, group string) error {
	return setOrganizationAuth(ctx, iq, name, roleName, group, MemberTypeGroup)
}

// SetOrganizationGroup sets the role and group that can have access to an organization
func SetOrganizationGroup(iq IQ, name, roleName, group string) error {
	return SetOrganizationGroupContext(context.Background(), iq, name, roleName, group)
}

func applicationAuthorizationsByID(ctx context.Context, iq IQ, appID string) ([]MemberMapping, error) {
	var endpoint string
	if hasRev70API(ctx, iq) {
		endpoint = fmt.Sprintf(restRoleMembersAppGet, appID)
	} else {
		endpoint = fmt.Sprintf(restRoleMembersAppDeprecated, appID)
	}

	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve role mapping for application %s: %v", appID, err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)

	return mappings.MemberMappings, err
}

func applicationAuthorizationsByRoleID(ctx context.Context, iq IQ, roleID string) ([]MemberMapping, error) {
	apps, err := GetAllApplicationsContext(ctx, iq)
	if err != nil {
		return nil, fmt.Errorf("could not find applications: %v", err)
	}

	mappings := make([]MemberMapping, 0)
	for _, app := range apps {
		appMaps, _ := applicationAuthorizationsByID(ctx, iq, app.ID)
		for _, m := range appMaps {
			if m.RoleID == roleID {
				mappings = append(mappings, m)
			}
		}
	}

	return mappings, nil
}

func ApplicationAuthorizationsContext(ctx context.Context, iq IQ, name string) ([]MemberMapping, error) {
	app, err := GetApplicationByPublicIDContext(ctx, iq, name)
	if err != nil {
		return nil, fmt.Errorf("could not find application with name %s: %v", name, err)
	}

	return applicationAuthorizationsByID(ctx, iq, app.ID)
}

// ApplicationAuthorizations returns the member mappings of an application
func ApplicationAuthorizations(iq IQ, name string) ([]MemberMapping, error) {
	return ApplicationAuthorizationsContext(context.Background(), iq, name)
}

func ApplicationAuthorizationsByRoleContext(ctx context.Context, iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	return applicationAuthorizationsByRoleID(ctx, iq, role.ID)
}

// ApplicationAuthorizationsByRole returns the member mappings of all applications which match the given role
func ApplicationAuthorizationsByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	return ApplicationAuthorizationsByRoleContext(context.Background(), iq, roleName)
}

func setApplicationAuth(ctx context.Context, iq IQ, name, roleName, member, memberType string) error {
	app, err := GetApplicationByPublicIDContext(ctx, iq, name)
	if err != nil {
		return fmt.Errorf("could not find application with name %s: %v", name, err)
	}

	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	var payload io.Reader
	if hasRev70API(ctx, iq) {
		switch memberType {
		case MemberTypeUser:
			endpoint = fmt.Sprintf(restRoleMembersAppUser, app.ID, role.ID, member)
		case MemberTypeGroup:
			endpoint = fmt.Sprintf(restRoleMembersAppGroup, app.ID, role.ID, member)
		}
	} else {
		endpoint = fmt.Sprintf(restRoleMembersAppDeprecated, app.ID)
		current, err := ApplicationAuthorizationsContext(ctx, iq, name)
		if err != nil && current == nil {
			current = make([]MemberMapping, 0)
		}
		current = append(current, newMapping(role.ID, memberType, member))

		buf, err := json.Marshal(memberMappings{MemberMappings: current})
		if err != nil {
			return fmt.Errorf("could not create mapping: %v", err)
		}
		payload = bytes.NewBuffer(buf)
	}

	_, _, err = iq.Put(ctx, endpoint, payload)
	if err != nil {
		return fmt.Errorf("could not update organization role mapping: %v", err)
	}

	return nil
}

func SetApplicationUserContext(ctx context.Context, iq IQ, name, roleName, user string) error {
	return setApplicationAuth(ctx, iq, name, roleName, user, MemberTypeUser)
}

// SetApplicationUser sets the role and user that can have access to an application
func SetApplicationUser(iq IQ, name, roleName, user string) error {
	return SetApplicationUserContext(context.Background(), iq, name, roleName, user)
}

func SetApplicationGroupContext(ctx context.Context, iq IQ, name, roleName, group string) error {
	return setApplicationAuth(ctx, iq, name, roleName, group, MemberTypeGroup)
}

// SetApplicationGroup sets the role and group that can have access to an application
func SetApplicationGroup(iq IQ, name, roleName, group string) error {
	return SetApplicationGroupContext(context.Background(), iq, name, roleName, group)
}

func revokeLT70(ctx context.Context, iq IQ, authType, authName, roleName, memberType, memberName string) error {
	var err error
	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var (
		authID, baseEndpoint string
		mapping              []MemberMapping
	)
	switch authType {
	case "organization":
		org, err := GetOrganizationByNameContext(ctx, iq, authName)
		if err == nil {
			authID = org.ID
			baseEndpoint = restRoleMembersOrgDeprecated
			mapping, err = OrganizationAuthorizationsContext(ctx, iq, authName)
		}
	case "application":
		app, err := GetApplicationByPublicIDContext(ctx, iq, authName)
		if err == nil {
			authID = app.ID
			baseEndpoint = restRoleMembersAppDeprecated
			mapping, err = ApplicationAuthorizationsContext(ctx, iq, authName)
		}
	}
	if err != nil && mapping != nil {
		return fmt.Errorf("could not get current authorizations for %s: %v", authName, err)
	}

	for i, auth := range mapping {
		if auth.RoleID == role.ID {
			for j, member := range auth.Members {
				if member.Type == memberType && member.UserOrGroupName == memberName {
					copy(mapping[i].Members[j:], mapping[i].Members[j+1:])
					mapping[i].Members[len(mapping[i].Members)-1] = Member{}
					mapping[i].Members = mapping[i].Members[:len(mapping[i].Members)-1]
				}
			}
		}
	}

	buf, err := json.Marshal(memberMappings{MemberMappings: mapping})
	if err != nil {
		return fmt.Errorf("could not create mapping: %v", err)
	}

	endpoint := fmt.Sprintf(baseEndpoint, authID)
	_, _, err = iq.Put(ctx, endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("could not remove role mapping: %v", err)
	}

	return nil
}

func revoke(ctx context.Context, iq IQ, authType, authName, roleName, memberType, memberName string) error {
	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var (
		authID, baseEndpoint string
	)
	switch authType {
	case "organization":
		org, err := GetOrganizationByNameContext(ctx, iq, authName)
		if err == nil {
			authID = org.ID
			switch memberType {
			case MemberTypeUser:
				baseEndpoint = restRoleMembersOrgUser
			case MemberTypeGroup:
				baseEndpoint = restRoleMembersOrgGroup
			}
		}
	case "application":
		app, err := GetApplicationByPublicIDContext(ctx, iq, authName)
		if err == nil {
			authID = app.ID
			switch memberType {
			case MemberTypeUser:
				baseEndpoint = restRoleMembersAppUser
			case MemberTypeGroup:
				baseEndpoint = restRoleMembersAppGroup
			}
		}
	}

	endpoint := fmt.Sprintf(baseEndpoint, authID, role.ID, memberName)
	_, err = iq.Del(ctx, endpoint)
	return err
}

func RevokeOrganizationUserContext(ctx context.Context, iq IQ, name, roleName, user string) error {
	if !hasRev70API(ctx, iq) {
		return revokeLT70(ctx, iq, "organization", name, roleName, MemberTypeUser, user)
	}
	return revoke(ctx, iq, "organization", name, roleName, MemberTypeUser, user)
}

func RevokeOrganizationUser(iq IQ, name, roleName, user string) error {
	return RevokeOrganizationUserContext(context.Background(), iq, name, roleName, user)
}

func RevokeOrganizationGroupContext(ctx context.Context, iq IQ, name, roleName, group string) error {
	if !hasRev70API(ctx, iq) {
		return revokeLT70(ctx, iq, "organization", name, roleName, MemberTypeGroup, group)
	}
	return revoke(ctx, iq, "organization", name, roleName, MemberTypeGroup, group)
}

// RevokeOrganizationGroup removes a group and role from the named organization
func RevokeOrganizationGroup(iq IQ, name, roleName, group string) error {
	return RevokeOrganizationGroupContext(context.Background(), iq, name, roleName, group)
}

func RevokeApplicationUserContext(ctx context.Context, iq IQ, name, roleName, user string) error {
	if !hasRev70API(ctx, iq) {
		return revokeLT70(ctx, iq, "application", name, roleName, MemberTypeUser, user)
	}
	return revoke(ctx, iq, "application", name, roleName, MemberTypeUser, user)
}

// RevokeApplicationUser removes a user and role from the named application
func RevokeApplicationUser(iq IQ, name, roleName, user string) error {
	return RevokeApplicationUserContext(context.Background(), iq, name, roleName, user)

}

func RevokeApplicationGroupContext(ctx context.Context, iq IQ, name, roleName, group string) error {
	if !hasRev70API(ctx, iq) {
		return revokeLT70(ctx, iq, "application", name, roleName, MemberTypeGroup, group)
	}
	return revoke(ctx, iq, "application", name, roleName, MemberTypeGroup, group)
}

// RevokeApplicationGroup removes a group and role from the named application
func RevokeApplicationGroup(iq IQ, name, roleName, group string) error {
	return RevokeApplicationGroupContext(context.Background(), iq, name, roleName, group)
}

func repositoriesAuth(ctx context.Context, iq IQ, method, roleName, memberType, member string) error {
	if !hasRev70API(ctx, iq) {
		return fmt.Errorf("did not find revision 70 API")
	}

	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	switch memberType {
	case MemberTypeUser:
		endpoint = fmt.Sprintf(restRoleMembersRepositoryUser, role.ID, member)
	case MemberTypeGroup:
		endpoint = fmt.Sprintf(restRoleMembersRepositoryGroup, role.ID, member)
	}

	switch method {
	case http.MethodPut:
		_, _, err = iq.Put(ctx, endpoint, nil)
	case http.MethodDelete:
		_, err = iq.Del(ctx, endpoint)
	}
	if err != nil {
		return fmt.Errorf("could not affect repositories role mapping: %v", err)
	}

	return nil
}

func repositoriesAuthorizationsByRoleID(ctx context.Context, iq IQ, roleID string) ([]MemberMapping, error) {
	auths, err := RepositoriesAuthorizationsContext(ctx, iq)
	if err != nil {
		return nil, fmt.Errorf("could not find authorization mappings for repositories: %v", err)
	}

	mappings := make([]MemberMapping, 0)
	for _, m := range auths {
		if m.RoleID == roleID {
			mappings = append(mappings, m)
		}
	}

	return mappings, nil
}

func RepositoriesAuthorizationsContext(ctx context.Context, iq IQ) ([]MemberMapping, error) {
	body, _, err := iq.Get(ctx, restRoleMembersReposGet)
	if err != nil {
		return nil, fmt.Errorf("could not get repositories mappings: %v", err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal mapping: %v", err)
	}

	return mappings.MemberMappings, nil
}

// RepositoriesAuthorizations returns the member mappings of all repositories
func RepositoriesAuthorizations(iq IQ) ([]MemberMapping, error) {
	return RepositoriesAuthorizationsContext(context.Background(), iq)
}

func RepositoriesAuthorizationsByRoleContext(ctx context.Context, iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	return repositoriesAuthorizationsByRoleID(ctx, iq, role.ID)
}

// RepositoriesAuthorizationsByRole returns the member mappings of all repositories which match the given role
func RepositoriesAuthorizationsByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	return RepositoriesAuthorizationsByRoleContext(context.Background(), iq, roleName)
}

func SetRepositoriesUserContext(ctx context.Context, iq IQ, roleName, user string) error {
	return repositoriesAuth(ctx, iq, http.MethodPut, roleName, MemberTypeUser, user)
}

// SetRepositoriesUser sets the role and user that can have access to the repositories
func SetRepositoriesUser(iq IQ, roleName, user string) error {
	return SetRepositoriesUserContext(context.Background(), iq, roleName, user)
}

func SetRepositoriesGroupContext(ctx context.Context, iq IQ, roleName, group string) error {
	return repositoriesAuth(ctx, iq, http.MethodPut, roleName, MemberTypeGroup, group)
}

// SetRepositoriesGroup sets the role and group that can have access to the repositories
func SetRepositoriesGroup(iq IQ, roleName, group string) error {
	return SetRepositoriesGroupContext(context.Background(), iq, roleName, group)
}

func RevokeRepositoriesUserContext(ctx context.Context, iq IQ, roleName, user string) error {
	return repositoriesAuth(ctx, iq, http.MethodDelete, roleName, MemberTypeUser, user)
}

// RevokeRepositoriesUser revoke the role and user that can have access to the repositories
func RevokeRepositoriesUser(iq IQ, roleName, user string) error {
	return RevokeRepositoriesUserContext(context.Background(), iq, roleName, user)
}

func RevokeRepositoriesGroupContext(ctx context.Context, iq IQ, roleName, group string) error {
	return repositoriesAuth(ctx, iq, http.MethodDelete, roleName, MemberTypeGroup, group)
}

// RevokeRepositoriesGroup revoke the role and group that can have access to the repositories
func RevokeRepositoriesGroup(iq IQ, roleName, group string) error {
	return RevokeRepositoriesGroupContext(context.Background(), iq, roleName, group)
}

func membersByRoleID(ctx context.Context, iq IQ, roleID string) ([]MemberMapping, error) {
	members := make([]MemberMapping, 0)

	if m, err := organizationAuthorizationsByRoleID(ctx, iq, roleID); err == nil && len(m) > 0 {
		members = append(members, m...)
	}

	if m, err := applicationAuthorizationsByRoleID(ctx, iq, roleID); err == nil && len(m) > 0 {
		members = append(members, m...)
	}

	if hasRev70API(ctx, iq) {
		if m, err := repositoriesAuthorizationsByRoleID(ctx, iq, roleID); err == nil && len(m) > 0 {
			members = append(members, m...)
		}
	}

	return members, nil
}

func MembersByRoleContext(ctx context.Context, iq IQ, roleName string) ([]MemberMapping, error) {
	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return nil, fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}
	return membersByRoleID(ctx, iq, role.ID)
}

// MembersByRole returns all users and groups by role name
func MembersByRole(iq IQ, roleName string) ([]MemberMapping, error) {
	return MembersByRoleContext(context.Background(), iq, roleName)
}

func GlobalAuthorizationsContext(ctx context.Context, iq IQ) ([]MemberMapping, error) {
	body, _, err := iq.Get(ctx, restRoleMembersGlobalGet)
	if err != nil {
		return nil, fmt.Errorf("could not get global members: %v", err)
	}

	var mappings memberMappings
	err = json.Unmarshal(body, &mappings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal mapping: %v", err)
	}

	return mappings.MemberMappings, nil
}

// GlobalAuthorizations returns all of the users and roles who have the administrator role across all of IQ
func GlobalAuthorizations(iq IQ) ([]MemberMapping, error) {
	return GlobalAuthorizationsContext(context.Background(), iq)
}

func globalAuth(ctx context.Context, iq IQ, method, roleName, memberType, member string) error {
	if !hasRev70API(ctx, iq) {
		return fmt.Errorf("did not find revision 70 API")
	}

	role, err := RoleByNameContext(ctx, iq, roleName)
	if err != nil {
		return fmt.Errorf("could not find role with name %s: %v", roleName, err)
	}

	var endpoint string
	switch memberType {
	case MemberTypeUser:
		endpoint = fmt.Sprintf(restRoleMembersGlobalUser, role.ID, member)
	case MemberTypeGroup:
		endpoint = fmt.Sprintf(restRoleMembersGlobalGroup, role.ID, member)
	}

	switch method {
	case http.MethodPut:
		_, _, err = iq.Put(ctx, endpoint, nil)
	case http.MethodDelete:
		_, err = iq.Del(ctx, endpoint)
	}
	if err != nil {
		return fmt.Errorf("could not affect global role mapping: %v", err)
	}

	return nil
}

func SetGlobalUserContext(ctx context.Context, iq IQ, roleName, user string) error {
	return globalAuth(ctx, iq, http.MethodPut, roleName, MemberTypeUser, user)
}

// SetGlobalUser sets the role and user that can have access to the repositories
func SetGlobalUser(iq IQ, roleName, user string) error {
	return SetGlobalUserContext(context.Background(), iq, roleName, user)
}

func SetGlobalGroupContext(ctx context.Context, iq IQ, roleName, group string) error {
	return globalAuth(ctx, iq, http.MethodPut, roleName, MemberTypeGroup, group)
}

// SetGlobalGroup sets the role and group that can have access to the global
func SetGlobalGroup(iq IQ, roleName, group string) error {
	return SetGlobalGroupContext(context.Background(), iq, roleName, group)
}

func RevokeGlobalUserContext(ctx context.Context, iq IQ, roleName, user string) error {
	return globalAuth(ctx, iq, http.MethodDelete, roleName, MemberTypeUser, user)
}

// RevokeGlobalUser revoke the role and user that can have access to the global
func RevokeGlobalUser(iq IQ, roleName, user string) error {
	return RevokeGlobalUserContext(context.Background(), iq, roleName, user)
}

func RevokeGlobalGroupContext(ctx context.Context, iq IQ, roleName, group string) error {
	return globalAuth(ctx, iq, http.MethodDelete, roleName, MemberTypeGroup, group)
}

// RevokeGlobalGroup revoke the role and group that can have access to the global
func RevokeGlobalGroup(iq IQ, roleName, group string) error {
	return RevokeGlobalGroupContext(context.Background(), iq, roleName, group)
}
