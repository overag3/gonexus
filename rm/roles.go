package nexusrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const restRole = "service/rest/v1/security/roles"

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Privileges  []string `json:"privileges"`
	Roles       []string `json:"roles"`
}

func CreateRoleContext(ctx context.Context, rm RM, role Role) error {
	json, err := json.Marshal(role)
	if err != nil {
		return err
	}

	_, resp, err := rm.Post(ctx, restRole, bytes.NewBuffer(json))
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}

	return nil
}

func CreateRole(rm RM, role Role) error {
	return CreateRoleContext(context.Background(), rm, role)
}

func DeleteRoleByIdContext(ctx context.Context, rm RM, id string) error {
	url := fmt.Sprintf("%s/%s", restRole, id)

	if resp, err := rm.Del(ctx, url); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("role not deleted '%s': %v", id, err)
	}

	return nil
}

func DeleteRoleById(rm RM, id string) error {
	return DeleteRoleByIdContext(context.Background(), rm, id)
}
