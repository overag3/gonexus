package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restApplication         = "api/v2/applications"
	restApplicationByPublic = "api/v2/applications?publicId=%s"
)

type iqNewAppRequest struct {
	PublicID        string `json:"publicId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName,omitempty"`
	ApplicationTags []struct {
		TagID string `json:"tagId"`
	} `json:"applicationTags,omitempty"`
}

type iqAppDetailsResponse struct {
	Applications []Application `json:"applications"`
}

type allAppsResponse struct {
	Applications []Application `json:"applications"`
}

// Application captures information of an IQ application
type Application struct {
	ID              string `json:"id"`
	PublicID        string `json:"publicId,omitempty"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	ContactUserName string `json:"contactUserName,omitempty"`
	ApplicationTags []struct {
		ID            string `json:"id,omitempty"`
		TagID         string `json:"tagId,omitempty"`
		ApplicationID string `json:"applicationId,omitempty"`
	} `json:"applicationTags,omitempty"`
}

func GetApplicationByPublicIDContext(ctx context.Context, iq IQ, applicationPublicID string) (*Application, error) {
	doError := func(err error) error {
		return fmt.Errorf("application '%s' not found: %v", applicationPublicID, err)
	}
	endpoint := fmt.Sprintf(restApplicationByPublic, applicationPublicID)
	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return nil, doError(err)
	}

	var resp iqAppDetailsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, doError(err)
	}

	if len(resp.Applications) == 0 {
		return nil, fmt.Errorf("application %s not found", applicationPublicID)
	}

	return &resp.Applications[0], nil
}

// GetApplicationByPublicID returns details on the named IQ application
func GetApplicationByPublicID(iq IQ, applicationPublicID string) (*Application, error) {
	return GetApplicationByPublicIDContext(context.Background(), iq, applicationPublicID)
}

func CreateApplicationContext(ctx context.Context, iq IQ, name, id, organizationID string) (string, error) {
	if name == "" || id == "" || organizationID == "" {
		return "", fmt.Errorf("cannot create application with empty values")
	}

	doError := func(err error) (string, error) {
		return "", fmt.Errorf("application '%s' not created: %v", name, err)
	}

	request, err := json.Marshal(iqNewAppRequest{Name: name, PublicID: id, OrganizationID: organizationID})
	if err != nil {
		return doError(err)
	}

	body, _, err := iq.Post(ctx, restApplication, bytes.NewBuffer(request))
	if err != nil {
		return doError(err)
	}

	var resp Application
	if err = json.Unmarshal(body, &resp); err != nil {
		return doError(err)
	}

	return resp.ID, nil
}

// CreateApplication creates an application in IQ with the given name and identifier
func CreateApplication(iq IQ, name, id, organizationID string) (string, error) {
	return CreateApplicationContext(context.Background(), iq, name, id, organizationID)
}

func DeleteApplicationContext(ctx context.Context, iq IQ, applicationID string) error {
	if resp, err := iq.Del(ctx, fmt.Sprintf("%s/%s", restApplication, applicationID)); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("application '%s' not deleted: %v", applicationID, err)
	}
	return nil
}

// DeleteApplication deletes an application in IQ with the given id
func DeleteApplication(iq IQ, applicationID string) error {
	return DeleteApplicationContext(context.Background(), iq, applicationID)
}

func GetAllApplicationsContext(ctx context.Context, iq IQ) ([]Application, error) {
	body, _, err := iq.Get(ctx, restApplication)
	if err != nil {
		return nil, fmt.Errorf("applications not found: %v", err)
	}

	var resp allAppsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("applications not found: %v", err)
	}

	return resp.Applications, nil
}

// GetAllApplications returns a slice of all of the applications in an IQ instance
func GetAllApplications(iq IQ) ([]Application, error) {
	return GetAllApplicationsContext(context.Background(), iq)
}

func GetApplicationsByOrganizationContext(ctx context.Context, iq IQ, organizationName string) ([]Application, error) {
	org, err := GetOrganizationByNameContext(ctx, iq, organizationName)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %v", err)
	}

	apps, err := GetAllApplicationsContext(ctx, iq)
	if err != nil {
		return nil, fmt.Errorf("could not get applications list: %v", err)
	}

	orgApps := make([]Application, 0)
	for _, app := range apps {
		if app.OrganizationID == org.ID {
			orgApps = append(orgApps, app)
		}
	}

	return orgApps, nil
}

// GetApplicationsByOrganization returns all applications under a given organization
func GetApplicationsByOrganization(iq IQ, organizationName string) ([]Application, error) {
	return GetApplicationsByOrganizationContext(context.Background(), iq, organizationName)
}
