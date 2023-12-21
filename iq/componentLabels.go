package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	restLabelComponent         = "api/v2/components/%s/labels/%s/applications/%s"
	restLabelComponentByOrg    = "api/v2/labels/organization/%s"
	restLabelComponentByOrgDel = "api/v2/labels/organization/%s/%s"
	restLabelComponentByApp    = "api/v2/labels/application/%s"
	restLabelComponentByAppDel = "api/v2/labels/application/%s/%s"
)

// IqComponentLabel describes a component label
type IqComponentLabel struct {
	ID             string `json:"id,omitempty"`
	OwnerID        string `json:"ownerId,omitempty"`
	Label          string `json:"label"`
	LabelLowercase string `json:"labelLowercase,omitempty"`
	Description    string `json:"description,omitempty"`
	Color          string `json:"color"`
}

func ComponentLabelApplyContext(ctx context.Context, iq IQ, comp Component, appID, label string) error {
	app, err := GetApplicationByPublicIDContext(ctx, iq, appID)
	if err != nil {
		return fmt.Errorf("could not retrieve application with ID %s: %v", appID, err)
	}

	endpoint := fmt.Sprintf(restLabelComponent, comp.Hash, url.PathEscape(label), app.ID)
	_, resp, err := iq.Post(ctx, endpoint, nil)
	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("could not apply label: %v", err)
		}
	}

	return nil
}

// ComponentLabelApply adds an existing label to a component for a given application
func ComponentLabelApply(iq IQ, comp Component, appID, label string) error {
	return ComponentLabelApplyContext(context.Background(), iq, comp, appID, label)
}

func ComponentLabelUnapplyContext(ctx context.Context, iq IQ, comp Component, appID, label string) error {
	app, err := GetApplicationByPublicIDContext(ctx, iq, appID)
	if err != nil {
		return fmt.Errorf("could not retrieve application with ID %s: %v", appID, err)
	}

	endpoint := fmt.Sprintf(restLabelComponent, comp.Hash, url.PathEscape(label), app.ID)
	resp, err := iq.Del(ctx, endpoint)
	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("could not unapply label: %v", err)
		}
	}

	return nil
}

// ComponentLabelUnapply removes an existing association between a label and a component
func ComponentLabelUnapply(iq IQ, comp Component, appID, label string) error {
	return ComponentLabelUnapplyContext(context.Background(), iq, comp, appID, label)
}

func getComponentLabels(ctx context.Context, iq IQ, endpoint string) ([]IqComponentLabel, error) {
	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var labels []IqComponentLabel
	err = json.Unmarshal(body, &labels)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func GetComponentLabelsByOrganizationContext(ctx context.Context, iq IQ, organization string) ([]IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByOrg, organization)
	return getComponentLabels(ctx, iq, endpoint)
}

// GetComponentLabelsByOrganization retrieves an array of an organization's component label
func GetComponentLabelsByOrganization(iq IQ, organization string) ([]IqComponentLabel, error) {
	return GetComponentLabelsByOrganizationContext(context.Background(), iq, organization)
}

func GetComponentLabelsByAppIDContext(ctx context.Context, iq IQ, appID string) ([]IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByApp, appID)
	return getComponentLabels(ctx, iq, endpoint)
}

// GetComponentLabelsByAppID retrieves an array of an organization's component label
func GetComponentLabelsByAppID(iq IQ, appID string) ([]IqComponentLabel, error) {
	return GetComponentLabelsByAppIDContext(context.Background(), iq, appID)
}

func createLabel(ctx context.Context, iq IQ, endpoint, label, description, color string) (IqComponentLabel, error) {
	var labelResponse IqComponentLabel
	request, err := json.Marshal(IqComponentLabel{Label: label, Description: description, Color: color})
	if err != nil {
		return labelResponse, fmt.Errorf("could not marshal label: %v", err)
	}

	body, resp, err := iq.Post(ctx, endpoint, bytes.NewBuffer(request))
	if resp.StatusCode != http.StatusOK {
		return labelResponse, fmt.Errorf("did not succeeed in creating label: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &labelResponse); err != nil {
		return labelResponse, fmt.Errorf("could not read json of new label: %v", err)
	}

	return labelResponse, nil
}

func CreateComponentLabelForOrganizationContext(ctx context.Context, iq IQ, organization, label, description, color string) (IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByOrg, organization)
	return createLabel(ctx, iq, endpoint, label, description, color)
}

// CreateComponentLabelForOrganization creates a label for an organization
func CreateComponentLabelForOrganization(iq IQ, organization, label, description, color string) (IqComponentLabel, error) {
	return CreateComponentLabelForOrganizationContext(context.Background(), iq, organization, label, description, color)
}

func CreateComponentLabelForApplicationContext(ctx context.Context, iq IQ, appID, label, description, color string) (IqComponentLabel, error) {
	endpoint := fmt.Sprintf(restLabelComponentByApp, appID)
	return createLabel(ctx, iq, endpoint, label, description, color)
}

// CreateComponentLabelForApplication creates a label for an application
func CreateComponentLabelForApplication(iq IQ, appID, label, description, color string) (IqComponentLabel, error) {
	return CreateComponentLabelForApplicationContext(context.Background(), iq, appID, label, description, color)
}

func DeleteComponentLabelForOrganizationContext(ctx context.Context, iq IQ, organization, label string) error {
	endpoint := fmt.Sprintf(restLabelComponentByOrgDel, organization, label)
	resp, err := iq.Del(ctx, endpoint)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("did not succeeed in deleting label: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// DeleteComponentLabelForOrganization deletes a label from an organization
func DeleteComponentLabelForOrganization(iq IQ, organization, label string) error {
	return DeleteComponentLabelForOrganizationContext(context.Background(), iq, organization, label)
}

func DeleteComponentLabelForApplicationContext(ctx context.Context, iq IQ, appID, label string) error {
	endpoint := fmt.Sprintf(restLabelComponentByAppDel, appID, label)
	resp, err := iq.Del(ctx, endpoint)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("did not succeeed in deleting label: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// DeleteComponentLabelForApplication deletes a label from an application
func DeleteComponentLabelForApplication(iq IQ, appID, label string) error {
	return DeleteComponentLabelForApplicationContext(context.Background(), iq, appID, label)
}
