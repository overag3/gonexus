package nexusrm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const restMaintenanceDBCheck = "service/rest/v1/maintenance/%s/check"

// Define database types
const (
	AccessLogDB = "accesslog"
	ComponentDB = "component"
	ConfigDB    = "config"
	SecurityDB  = "security"
)

// DatabaseState contains state information about a given state
type DatabaseState struct {
	PageCorruption bool `json:"pageCorruption"`
	IndexErrors    int  `json:"indexErrors"`
}

func CheckDatabaseContext(ctx context.Context, rm RM, dbName string) (DatabaseState, error) {
	doError := func(err error) error {
		return fmt.Errorf("error checking status of database '%s': %v", dbName, err)
	}

	var state DatabaseState

	url := fmt.Sprintf(restMaintenanceDBCheck, dbName)
	body, resp, err := rm.Put(ctx, url, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return state, doError(err)
	}

	if err := json.Unmarshal(body, &state); err != nil {
		return state, doError(err)
	}

	return state, nil
}

// CheckDatabase returns the state of the named database
func CheckDatabase(rm RM, dbName string) (DatabaseState, error) {
	return CheckDatabaseContext(context.Background(), rm, dbName)
}

func CheckAllDatabasesContext(ctx context.Context, rm RM) (states map[string]DatabaseState, err error) {
	states = make(map[string]DatabaseState)

	check := func(dbName string) {
		if err != nil {
			return
		}

		if state, er := CheckDatabaseContext(ctx, rm, dbName); er != nil {
			err = fmt.Errorf("error with '%s' database when all states: %v", dbName, er)
		} else {
			states[dbName] = state
		}
	}

	check(AccessLogDB)
	check(ComponentDB)
	check(ConfigDB)
	check(SecurityDB)

	return
}

// CheckAllDatabases returns state on all of the databases
func CheckAllDatabases(rm RM) (states map[string]DatabaseState, err error) {
	return CheckAllDatabasesContext(context.Background(), rm)
}
