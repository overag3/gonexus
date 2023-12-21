package nexusrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restScript    = "service/rest/v1/script"
	restScriptRun = "service/rest/v1/script/%s/run"
)

// Script encapsulates a Repository Manager script
type Script struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type runResponse struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

func ScriptListContext(ctx context.Context, rm RM) ([]Script, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not list scripts: %v", err)
	}

	body, _, err := rm.Get(ctx, restScript)
	if err != nil {
		return nil, doError(err)
	}

	scripts := make([]Script, 0)
	if err = json.Unmarshal(body, &scripts); err != nil {
		return nil, doError(err)
	}

	return scripts, nil
}

// ScriptList lists all of the uploaded scripts in Repository Manager
func ScriptList(rm RM) ([]Script, error) {
	return ScriptListContext(context.Background(), rm)
}

func ScriptGetContext(ctx context.Context, rm RM, name string) (Script, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not find script '%s': %v", name, err)
	}

	var script Script

	endpoint := fmt.Sprintf("%s/%s", restScript, name)
	body, _, err := rm.Get(ctx, endpoint)
	if err != nil {
		return script, doError(err)
	}

	if err = json.Unmarshal(body, &script); err != nil {
		return script, doError(err)
	}

	return script, nil
}

// ScriptGet returns the named script
func ScriptGet(rm RM, name string) (Script, error) {
	return ScriptGetContext(context.Background(), rm, name)
}

func ScriptUploadContext(ctx context.Context, rm RM, script Script) error {
	doError := func(err error) error {
		return fmt.Errorf("could not upload script '%s': %v", script.Name, err)
	}

	json, err := json.Marshal(script)
	if err != nil {
		return doError(err)
	}

	_, resp, err := rm.Post(ctx, restScript, bytes.NewBuffer(json))
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return doError(err)
	}

	return nil
}

// ScriptUpload uploads the given Script to Repository Manager
func ScriptUpload(rm RM, script Script) error {
	return ScriptUploadContext(context.Background(), rm, script)
}

func ScriptUpdateContext(ctx context.Context, rm RM, script Script) error {
	doError := func(err error) error {
		return fmt.Errorf("could not update script '%s': %v", script.Name, err)
	}

	json, err := json.Marshal(script)
	if err != nil {
		return doError(err)
	}

	endpoint := fmt.Sprintf("%s/%s", restScript, script.Name)
	_, resp, err := rm.Put(ctx, endpoint, bytes.NewBuffer(json))
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return doError(err)
	}

	return nil
}

// ScriptUpdate update the contents of the given script
func ScriptUpdate(rm RM, script Script) error {
	return ScriptUpdateContext(context.Background(), rm, script)
}

func ScriptRunContext(ctx context.Context, rm RM, name string, arguments []byte) (string, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not run script '%s': %v", name, err)
	}

	endpoint := fmt.Sprintf(restScriptRun, name)
	body, _, err := rm.Post(ctx, endpoint, bytes.NewBuffer(arguments)) // TODO: Better response handling
	if err != nil {
		return "", doError(err)
	}

	var resp runResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", doError(err)
	}

	return resp.Result, nil
}

// ScriptRun executes the named Script
func ScriptRun(rm RM, name string, arguments []byte) (string, error) {
	return ScriptRunContext(context.Background(), rm, name, arguments)
}

func ScriptRunOnceContext(ctx context.Context, rm RM, script Script, arguments []byte) (string, error) {
	if err := ScriptUploadContext(ctx, rm, script); err != nil {
		return "", err
	}
	defer ScriptDeleteContext(ctx, rm, script.Name)

	return ScriptRunContext(ctx, rm, script.Name, arguments)
}

// ScriptRunOnce takes the given Script, uploads it, executes it, and deletes it
func ScriptRunOnce(rm RM, script Script, arguments []byte) (string, error) {
	return ScriptRunOnceContext(context.Background(), rm, script, arguments)
}

func ScriptDeleteContext(ctx context.Context, rm RM, name string) error {
	endpoint := fmt.Sprintf("%s/%s", restScript, name)
	resp, err := rm.Del(ctx, endpoint)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete '%s': %v", name, err)
	}
	return nil
}

// ScriptDelete removes the name, uploaded script
func ScriptDelete(rm RM, name string) error {
	return ScriptDeleteContext(context.Background(), rm, name)
}
