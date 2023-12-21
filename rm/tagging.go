package nexusrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const restTagging = "service/rest/v1/tags"

// Tag contains the information about a component tag
type Tag struct {
	Name         string   `json:"name"`
	Attributes   struct{} `json:"attributes,omitempty"`
	FirstCreated string   `json:"firstCreated,omitempty"`
	LastUpdated  string   `json:"lastUpdated,omitempty"`
}

type tagsResponse struct {
	Items             []Tag  `json:"items"`
	ContinuationToken string `json:"continuationToken"`
}

type associateResponse struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ComponentsAssociated []componentsAssociated `json:"components associated"`
	} `json:"data"`
}

type componentsAssociated struct {
	Name    string `json:"name"`
	Group   string `json:"group"`
	Version string `json:"version"`
}

func TagsListContext(ctx context.Context, rm RM) ([]Tag, error) {
	continuation := ""
	tags := make([]Tag, 0)

	get := func() error {
		url := restTagging

		if continuation != "" {
			url += "&continuationToken=" + continuation
		}

		body, _, err := rm.Get(ctx, url)
		if err != nil {
			return fmt.Errorf("could not get list of tags: %v", err)
		}

		var resp tagsResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return fmt.Errorf("could not read tag list response: %v", err)
		}

		continuation = resp.ContinuationToken

		return nil
	}

	for {
		if err := get(); err != nil {
			return nil, err
		}

		if continuation == "" {
			break
		}
	}

	return tags, nil
}

// TagsList returns a list of tags in the given RM instance
func TagsList(rm RM) ([]Tag, error) {
	return TagsListContext(context.Background(), rm)
}

func AddTagContext(ctx context.Context, rm RM, tagName string, attributes map[string]string) (Tag, error) {
	tag := Tag{Name: tagName}
	//TODO: attributes

	buf, err := json.Marshal(tag)
	if err != nil {
		return Tag{}, fmt.Errorf("could not marshal tag: %v", err)
	}

	body, _, err := rm.Post(ctx, restTagging, bytes.NewBuffer(buf))
	if err != nil {
		return Tag{}, fmt.Errorf("could not create tag %s: %v", tagName, err)
	}

	var createdTag Tag
	if err = json.Unmarshal(body, &createdTag); err != nil {
		return Tag{}, fmt.Errorf("could not read response: %v", err)
	}

	return createdTag, nil
}

// AddTag adds a tag to the given instance
func AddTag(rm RM, tagName string, attributes map[string]string) (Tag, error) {
	return AddTagContext(context.Background(), rm, tagName, attributes)
}

func GetTagContext(ctx context.Context, rm RM, tagName string) (Tag, error) {
	endpoint := fmt.Sprintf("%s/%s", restTagging, tagName)

	body, _, err := rm.Get(ctx, endpoint)
	if err != nil {
		return Tag{}, fmt.Errorf("could not find tag %s: %v", tagName, err)
	}

	var tag Tag
	if err = json.Unmarshal(body, &tag); err != nil {
		return Tag{}, fmt.Errorf("could not read response: %v", err)
	}

	return tag, nil
}

// GetTag retrieve the named tag
func GetTag(rm RM, tagName string) (Tag, error) {
	return GetTagContext(context.Background(), rm, tagName)
}

func AssociateTagContext(ctx context.Context, rm RM, query QueryBuilder) error {
	endpoint := fmt.Sprintf("%s?%s", restTagging, query.Build())

	// TODO: handle response
	_, _, err := rm.Post(ctx, endpoint, nil)
	return err
}

// AssociateTag associates a tag to any component which matches the search criteria
func AssociateTag(rm RM, query QueryBuilder) error {
	return AssociateTagContext(context.Background(), rm, query)
}

func DisassociateTagContext(ctx context.Context, rm RM, query QueryBuilder) error {
	endpoint := fmt.Sprintf("%s?%s", restTagging, query.Build())

	_, err := rm.Del(ctx, endpoint)
	return err
}

// DisassociateTag associates a tag to any component which matches the search criteria
func DisassociateTag(rm RM, query QueryBuilder) error {
	return DisassociateTagContext(context.Background(), rm, query)
}
