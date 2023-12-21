package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const restComponentVersions = "api/v2/components/versions"

func ComponentVersionsContext(ctx context.Context, iq IQ, comp Component) (versions []string, err error) {
	str, err := json.Marshal(comp)
	if err != nil {
		return nil, fmt.Errorf("could not process component: %v", err)
	}

	body, _, err := iq.Post(ctx, restComponentVersions, bytes.NewBuffer(str))
	if err != nil {
		return nil, fmt.Errorf("could not request component: %v", err)
	}

	if err := json.Unmarshal(body, &versions); err != nil {
		return nil, fmt.Errorf("could not process versions list: %v", err)
	}

	return
}

// ComponentVersions returns all known versions of a given component
func ComponentVersions(iq IQ, comp Component) (versions []string, err error) {
	return ComponentVersionsContext(context.Background(), iq, comp)
}
