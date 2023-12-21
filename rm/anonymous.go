package nexusrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const restAnonymous = "service/rest/v1/security/anonymous"

type SettingsAnonAccess struct {
	Enabled   bool   `json:"enabled"`
	UserId    string `json:"userId"`
	RealmName string `json:"realmName"`
}

func GetAnonAccessContext(ctx context.Context, rm RM) (SettingsAnonAccess, error) {
	var settings SettingsAnonAccess

	body, resp, err := rm.Get(ctx, restAnonymous)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return SettingsAnonAccess{}, fmt.Errorf("anonymous access settings can't getting: %v", err)
	}

	if err := json.Unmarshal(body, &settings); err != nil {
		return SettingsAnonAccess{}, fmt.Errorf("anonymous access settings can't getting: %v", err)
	}

	return settings, nil
}

func GetAnonAccess(rm RM) (SettingsAnonAccess, error) {
	return GetAnonAccessContext(context.Background(), rm)
}

func SetAnonAccessContext(ctx context.Context, rm RM, settings SettingsAnonAccess) error {
	json, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	if _, resp, err := rm.Put(ctx, restAnonymous, bytes.NewBuffer(json)); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("email config not set: %v", err)
	}

	return nil
}

func SetAnonAccess(rm RM, settings SettingsAnonAccess) error {
	return SetAnonAccessContext(context.Background(), rm, settings)
}
