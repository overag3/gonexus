package nexusrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const restEmail = "service/rest/v1/email"

type EmailConfig struct {
	Enabled                       bool   `json:"enabled"`
	Host                          string `json:"host"`
	Port                          int32  `json:"port"`
	Username                      string `json:"username"`
	Password                      string `json:"password"`
	FromAddress                   string `json:"fromAddress"`
	SubjectPrefix                 string `json:"subjectPrefix"`
	StartTlsEnabled               bool   `json:"startTlsEnabled"`
	StartTlsRequired              bool   `json:"startTlsRequired"`
	SslOnConnectEnabled           bool   `json:"sslOnConnectEnabled"`
	SslServerIdentityCheckEnabled bool   `json:"sslServerIdentityCheckEnabled"`
	NexusTrustStoreEnabled        bool   `json:"nexusTrustStoreEnabled"`
}

func SetEmailConfigContext(ctx context.Context, rm RM, config EmailConfig) error {
	json, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if _, resp, err := rm.Put(ctx, restEmail, bytes.NewBuffer(json)); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("email config not set: %v", err)
	}

	return nil
}

func SetEmailConfig(rm RM, config EmailConfig) error {
	return SetEmailConfigContext(context.Background(), rm, config)
}

func GetEmailConfigContext(ctx context.Context, rm RM) (EmailConfig, error) {
	var config EmailConfig

	body, resp, err := rm.Get(ctx, restEmail)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return EmailConfig{}, fmt.Errorf("email config can't getting: %v", err)
	}

	if err := json.Unmarshal(body, &config); err != nil {
		return EmailConfig{}, fmt.Errorf("email config can't getting: %v", err)
	}

	return config, nil
}

func GetEmailConfig(rm RM) (EmailConfig, error) {
	return GetEmailConfigContext(context.Background(), rm)
}

func DeleteEmailConfigContext(ctx context.Context, rm RM) error {
	if resp, err := rm.Del(ctx, restEmail); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("email config not deleted: %v", err)
	}

	return nil
}

func DeleteEmailConfig(rm RM) error {
	return DeleteEmailConfigContext(context.Background(), rm)
}
