package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restSourceControl       = "api/v2/sourceControl/%s"
	restSourceControlDelete = "api/v2/sourceControl/%s/%s"
)

// SourceControlEntry describes a Source Control entry in IQ
type SourceControlEntry struct {
	ID            string `json:"id,omitempty"`
	ApplicationID string `json:"applicationId"`
	RepositoryURL string `json:"repositoryUrl"`
	Token         string `json:"token"`
}

func getSourceControlEntryByInternalID(ctx context.Context, iq IQ, applicationID string) (entry SourceControlEntry, err error) {
	endpoint := fmt.Sprintf(restSourceControl, applicationID)

	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &entry)

	return
}

func GetSourceControlEntryContext(ctx context.Context, iq IQ, applicationID string) (SourceControlEntry, error) {
	appInfo, err := GetApplicationByPublicIDContext(ctx, iq, applicationID)
	if err != nil {
		return SourceControlEntry{}, fmt.Errorf("no source control entry for '%s': %v", applicationID, err)
	}

	return getSourceControlEntryByInternalID(ctx, iq, appInfo.ID)
}

// GetSourceControlEntry lists of all of the Source Control entries for the given application
func GetSourceControlEntry(iq IQ, applicationID string) (SourceControlEntry, error) {
	return GetSourceControlEntryContext(context.Background(), iq, applicationID)
}

func GetAllSourceControlEntriesContext(ctx context.Context, iq IQ) ([]SourceControlEntry, error) {
	apps, err := GetAllApplicationsContext(ctx, iq)
	if err != nil {
		return nil, fmt.Errorf("no source control entries: %v", err)
	}

	entries := make([]SourceControlEntry, 0)
	for _, app := range apps {
		if entry, err := getSourceControlEntryByInternalID(ctx, iq, app.ID); err == nil {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// GetAllSourceControlEntries lists of all of the Source Control entries in the IQ instance
func GetAllSourceControlEntries(iq IQ) ([]SourceControlEntry, error) {
	return GetAllSourceControlEntriesContext(context.Background(), iq)
}

func CreateSourceControlEntryContext(ctx context.Context, iq IQ, applicationID, repositoryURL, token string) error {
	doError := func(err error) error {
		return fmt.Errorf("source control entry not created for '%s': %v", applicationID, err)
	}

	appInfo, err := GetApplicationByPublicIDContext(ctx, iq, applicationID)
	if err != nil {
		return doError(err)
	}

	request, err := json.Marshal(SourceControlEntry{"", appInfo.ID, repositoryURL, token})
	if err != nil {
		return doError(err)
	}

	endpoint := fmt.Sprintf(restSourceControl, appInfo.ID)
	if _, _, err = iq.Post(ctx, endpoint, bytes.NewBuffer(request)); err != nil {
		return doError(err)
	}

	return nil
}

// CreateSourceControlEntry creates a source control entry in IQ
func CreateSourceControlEntry(iq IQ, applicationID, repositoryURL, token string) error {
	return CreateSourceControlEntryContext(context.Background(), iq, applicationID, repositoryURL, token)
}

func UpdateSourceControlEntryContext(ctx context.Context, iq IQ, applicationID, repositoryURL, token string) error {
	doError := func(err error) error {
		return fmt.Errorf("source control entry not updated for '%s': %v", applicationID, err)
	}

	appInfo, err := GetApplicationByPublicIDContext(ctx, iq, applicationID)
	if err != nil {
		return doError(err)
	}

	request, err := json.Marshal(SourceControlEntry{"", appInfo.ID, repositoryURL, token})
	if err != nil {
		return doError(err)
	}

	endpoint := fmt.Sprintf(restSourceControl, appInfo.ID)
	if _, _, err = iq.Put(ctx, endpoint, bytes.NewBuffer(request)); err != nil {
		return doError(err)
	}

	return nil
}

// UpdateSourceControlEntry updates a source control entry in IQ
func UpdateSourceControlEntry(iq IQ, applicationID, repositoryURL, token string) error {
	return UpdateSourceControlEntryContext(context.Background(), iq, applicationID, repositoryURL, token)
}

func deleteSourceControlEntry(ctx context.Context, iq IQ, appInternalID, sourceControlID string) error {
	endpoint := fmt.Sprintf(restSourceControlDelete, appInternalID, sourceControlID)

	resp, err := iq.Del(ctx, endpoint)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

func DeleteSourceControlEntryContext(ctx context.Context, iq IQ, applicationID, sourceControlID string) error {
	appInfo, err := GetApplicationByPublicIDContext(ctx, iq, applicationID)
	if err != nil {
		return fmt.Errorf("source control entry not deleted from '%s': %v", applicationID, err)
	}

	return deleteSourceControlEntry(ctx, iq, appInfo.ID, sourceControlID)
}

// DeleteSourceControlEntry deletes a source control entry in IQ
func DeleteSourceControlEntry(iq IQ, applicationID, sourceControlID string) error {
	return DeleteSourceControlEntryContext(context.Background(), iq, applicationID, sourceControlID)
}

func DeleteSourceControlEntryByAppContext(ctx context.Context, iq IQ, applicationID string) error {
	doError := func(err error) error {
		return fmt.Errorf("source control entry not deleted from '%s': %v", applicationID, err)
	}

	appInfo, err := GetApplicationByPublicIDContext(ctx, iq, applicationID)
	if err != nil {
		return doError(err)
	}

	entry, err := getSourceControlEntryByInternalID(ctx, iq, appInfo.ID)
	if err != nil {
		return doError(err)
	}

	return deleteSourceControlEntry(ctx, iq, appInfo.ID, entry.ID)
}

// DeleteSourceControlEntryByApp deletes a source control entry in IQ for the given application
func DeleteSourceControlEntryByApp(iq IQ, applicationID string) error {
	return DeleteSourceControlEntryByAppContext(context.Background(), iq, applicationID)
}

// DeleteSourceControlEntryByEntry deletes a source control entry in IQ for the given entry ID
/*
func DeleteSourceControlEntryByEntry(iq IQ, sourceControlID string) error {
	entry, err := getSourceControlEntryByInternalID(iq, appInfo.ID)
	if err != nil {
		return err
	}

	return deleteSourceControlEntry(iq, entry.ApplicationID, entry.ID)
}
*/
