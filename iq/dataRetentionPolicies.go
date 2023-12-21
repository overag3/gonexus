package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const restDataRetentionPolicies = "api/v2/dataRetentionPolicies/organizations/%s"

// DataRetentionPolicies encapsulates an organization's retention policies
type DataRetentionPolicies struct {
	ApplicationReports ApplicationReports  `json:"applicationReports"`
	SuccessMetrics     DataRetentionPolicy `json:"successMetrics"`
}

// ApplicationReports captures the policies related to application reports
type ApplicationReports struct {
	Stages map[Stage]DataRetentionPolicy `json:"stages"`
}

// DataRetentionPolicy describes the retention policies for a pipeline stage
type DataRetentionPolicy struct {
	InheritPolicy bool   `json:"inheritPolicy"`
	EnablePurging bool   `json:"enablePurging"`
	MaxAge        string `json:"maxAge"`
}

func GetRetentionPoliciesContext(ctx context.Context, iq IQ, orgName string) (policies DataRetentionPolicies, err error) {
	org, err := GetOrganizationByNameContext(ctx, iq, orgName)
	if err != nil {
		return policies, fmt.Errorf("could not retrieve organization named %s: %v", orgName, err)
	}

	endpoint := fmt.Sprintf(restDataRetentionPolicies, org.ID)

	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return policies, fmt.Errorf("did not retrieve retention policies for organization %s: %v", orgName, err)
	}

	err = json.Unmarshal(body, &policies)

	return
}

// GetRetentionPolicies returns the current retention policies
func GetRetentionPolicies(iq IQ, orgName string) (policies DataRetentionPolicies, err error) {
	return GetRetentionPoliciesContext(context.Background(), iq, orgName)
}

func SetRetentionPoliciesContext(ctx context.Context, iq IQ, orgName string, policies DataRetentionPolicies) error {
	org, err := GetOrganizationByNameContext(ctx, iq, orgName)
	if err != nil {
		return fmt.Errorf("could not retrieve organization named %s: %v", orgName, err)
	}

	request, err := json.Marshal(policies)
	if err != nil {
		return fmt.Errorf("could not parse policies: %v", err)
	}

	endpoint := fmt.Sprintf(restDataRetentionPolicies, org.ID)

	_, _, err = iq.Put(ctx, endpoint, bytes.NewBuffer(request))
	if err != nil {
		return fmt.Errorf("did not set retention policies for organization %s: %v", orgName, err)
	}

	return nil
}

// SetRetentionPolicies updates the retention policies
func SetRetentionPolicies(iq IQ, orgName string, policies DataRetentionPolicies) error {
	return SetRetentionPoliciesContext(context.Background(), iq, orgName, policies)
}
