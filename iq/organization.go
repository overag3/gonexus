package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const (
	restOrganization = "api/v2/organizations"

	// RootOrganization is the ID of the ... Root ... Organization
	RootOrganization = "ROOT_ORGANIZATION_ID"
)

type iqNewOrgRequest struct {
	Name string `json:"name"`
}

type allOrgsResponse struct {
	Organizations []Organization `json:"organizations"`
}

// IQCategory encapsulates the category that can be created in IQ
type IQCategory struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Organization describes the data in IQ about a given organization
type Organization struct {
	ID   string       `json:"id"`
	Name string       `json:"name"`
	Tags []IQCategory `json:"tags,omitempty"`
}

func GetOrganizationByNameContext(ctx context.Context, iq IQ, organizationName string) (*Organization, error) {
	orgs, err := GetAllOrganizationsContext(ctx, iq)
	if err != nil {
		return nil, fmt.Errorf("organization '%s' not found: %v", organizationName, err)
	}
	for _, org := range orgs {
		if org.Name == organizationName {
			return &org, nil
		}
	}

	return nil, fmt.Errorf("organization '%s' not found", organizationName)
}

// GetOrganizationByName returns details on the named IQ organization
func GetOrganizationByName(iq IQ, organizationName string) (*Organization, error) {
	return GetOrganizationByNameContext(context.Background(), iq, organizationName)
}

func CreateOrganizationContext(ctx context.Context, iq IQ, name string) (string, error) {
	doError := func(err error) error {
		return fmt.Errorf("organization '%s' not created: %v", name, err)
	}

	request, err := json.Marshal(iqNewOrgRequest{Name: name})
	if err != nil {
		return "", doError(err)
	}

	body, _, err := iq.Post(ctx, restOrganization, bytes.NewBuffer(request))
	if err != nil {
		return "", doError(err)
	}

	var org Organization
	if err = json.Unmarshal(body, &org); err != nil {
		return "", doError(err)
	}

	return org.ID, nil
}

// CreateOrganization creates an organization in IQ with the given name
func CreateOrganization(iq IQ, name string) (string, error) {
	return CreateOrganizationContext(context.Background(), iq, name)
}

func GetAllOrganizationsContext(ctx context.Context, iq IQ) ([]Organization, error) {
	doError := func(err error) error {
		return fmt.Errorf("organizations not found: %v", err)
	}

	body, _, err := iq.Get(ctx, restOrganization)
	if err != nil {
		return nil, doError(err)
	}

	var resp allOrgsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, doError(err)
	}

	return resp.Organizations, nil
}

// GetAllOrganizations returns a slice of all of the organizations in an IQ instance
func GetAllOrganizations(iq IQ) ([]Organization, error) {
	return GetAllOrganizationsContext(context.Background(), iq)
}
