package nexusiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restUsers     = "api/v2/users/%s"
	restUsersPost = "api/v2/users"
)

// User encapsulates the information of a user in IQ
type User struct {
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

func GetUserContext(ctx context.Context, iq IQ, username string) (user User, err error) {
	endpoint := fmt.Sprintf(restUsers, username)
	body, _, err := iq.Get(ctx, endpoint)
	if err != nil {
		return user, fmt.Errorf("could not retrieve details on username %s: %v", username, err)
	}

	err = json.Unmarshal(body, &user)

	return user, err
}

// GetUser returns user details for the given name
func GetUser(iq IQ, username string) (user User, err error) {
	return GetUserContext(context.Background(), iq, username)
}

func SetUserContext(ctx context.Context, iq IQ, user User) (err error) {
	buf, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("could not read user details: %v", err)
	}
	str := bytes.NewBuffer(buf)

	if _, er := GetUserContext(ctx, iq, user.Username); er != nil {
		_, resp, er := iq.Post(ctx, restUsersPost, str)
		if er != nil && resp.StatusCode != http.StatusNoContent {
			return er
		}
	} else {
		endpoint := fmt.Sprintf(restUsers, user.Username)
		_, _, err = iq.Put(ctx, endpoint, str)
	}

	return err
}

// SetUser creates a new user
func SetUser(iq IQ, user User) (err error) {
	return SetUserContext(context.Background(), iq, user)
}

func DeleteUserContext(ctx context.Context, iq IQ, username string) error {
	endpoint := fmt.Sprintf(restUsers, username)
	if resp, err := iq.Del(ctx, endpoint); err != nil && resp.StatusCode != http.StatusNoContent {
		return err
	}
	return nil
}

// DeleteUser removes the named user
func DeleteUser(iq IQ, username string) error {
	return DeleteUserContext(context.Background(), iq, username)
}
