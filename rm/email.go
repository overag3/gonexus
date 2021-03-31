package nexusrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restEmail = "service/rest/v1/email"
)

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

func SetEmailConfig(rm RM, config EmailConfig) error {

	json, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if _, resp, err := rm.Put(restEmail, bytes.NewBuffer(json)); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("email config not set: %v", err)
	}

	return nil
}
