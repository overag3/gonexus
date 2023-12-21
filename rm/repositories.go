package nexusrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	restRepositories               = "service/rest/v1/repositories"
	restRepositoriesHostedApt      = "service/rest/v1/repositories/apt/hosted"
	restRepositoriesHostedBower    = "service/rest/v1/repositories/bower/hosted"
	restRepositoriesHostedDocker   = "service/rest/v1/repositories/docker/hosted"
	restRepositoriesHostedGitLfs   = "service/rest/v1/repositories/gitlfs/hosted"
	restRepositoriesHostedHelm     = "service/rest/v1/repositories/helm/hosted"
	restRepositoriesHostedMaven    = "service/rest/v1/repositories/maven/hosted"
	restRepositoriesHostedNpm      = "service/rest/v1/repositories/npm/hosted"
	restRepositoriesHostedNuget    = "service/rest/v1/repositories/nuget/hosted"
	restRepositoriesHostedPypi     = "service/rest/v1/repositories/pypi/hosted"
	restRepositoriesHostedR        = "service/rest/v1/repositories/r/hosted"
	restRepositoriesHostedRaw      = "service/rest/v1/repositories/raw/hosted"
	restRepositoriesHostedRubygems = "service/rest/v1/repositories/rubygems/hosted"
	restRepositoriesHostedYum      = "service/rest/v1/repositories/yum/hosted"

	restRepositoriesProxyApt       = "service/rest/v1/repositories/apt/proxy"
	restRepositoriesProxyBower     = "service/rest/v1/repositories/bower/proxy"
	restRepositoriesProxyCocoapods = "service/rest/v1/repositories/cocoapods/proxy"
	restRepositoriesProxyConan     = "service/rest/v1/repositories/conan/proxy"
	restRepositoriesProxyConda     = "service/rest/v1/repositories/conda/proxy"
	restRepositoriesProxyDocker    = "service/rest/v1/repositories/docker/proxy"
	restRepositoriesProxyGolang    = "service/rest/v1/repositories/go/proxy"
	restRepositoriesProxyHelm      = "service/rest/v1/repositories/helm/proxy"
	restRepositoriesProxyMaven     = "service/rest/v1/repositories/maven/proxy"
	restRepositoriesProxyNpm       = "service/rest/v1/repositories/npm/proxy"
	restRepositoriesProxyNuget     = "service/rest/v1/repositories/nuget/proxy"
	restRepositoriesProxyP2        = "service/rest/v1/repositories/p2/proxy"
	restRepositoriesProxyPypi      = "service/rest/v1/repositories/pypi/proxy"
	restRepositoriesProxyR         = "service/rest/v1/repositories/r/proxy"
	restRepositoriesProxyRaw       = "service/rest/v1/repositories/raw/proxy"
	restRepositoriesProxyRubygems  = "service/rest/v1/repositories/rubygems/proxy"
	restRepositoriesProxyYum       = "service/rest/v1/repositories/yum/proxy"
)

/*
// RepositoryType enumerates the types of repositories in RM
type repositoryType int

const (
	Hosted = iota
	Proxy
	Group
)

func (r repositoryType) String() string {
	switch r {
	case proxy:
		return "proxy"
	case hosted:
		return "hosted"
	case group:
		return "group"
	default:
		return ""
	}
}
*/

type repositoryFormat int

// Enumerates the formats which can be created as Repository Manager repositories
const (
	Unknown repositoryFormat = iota
	Apt
	Bower
	Cocoapods
	Conan
	Conda
	Docker
	GitLfs
	Golang
	Helm
	Maven
	Npm
	Nuget
	P2
	Pypi
	R
	Raw
	Rubygems
	Yum
)

// Repository collects the information returned by RM about a repository
type Repository struct {
	Name       string `json:"name"`
	Format     string `json:"format"`
	Type       string `json:"type"`
	URL        string `json:"url"`
	Attributes struct {
		Proxy struct {
			RemoteURL string `json:"remoteUrl"`
		} `json:"proxy"`
	} `json:"attributes,omitempty"`
}

type AttributesStorageHosted struct {
	BlobStoreName               string `json:"blobStoreName"`
	StrictContentTypeValidation bool   `json:"strictContentTypeValidation"`
	WritePolicy                 string `json:"writePolicy"`
}
type AttributesStorage struct {
	BlobStoreName               string `json:"blobStoreName"`
	StrictContentTypeValidation bool   `json:"strictContentTypeValidation"`
}

type AttributesCleanupPolicy struct {
	PolicyNames []string `json:"policyNames"`
}

type AttributesRepositoriesAptHosted struct {
	Distribution string `json:"distribution"`
}

type AttributesRepositoriesAptSigning struct {
	Keypair    string `json:"keypair"`
	Passphrase string `json:"passphrase"`
}

type AttributesComponent struct {
	ProprietaryComponents bool `json:"proprietaryComponents"`
}

type AttributesRaw struct {
	ContentDisposition string `json:"contentDisposition"`
}

type RepositoryNugetHosted struct {
	Name      string                  `json:"name"`
	Online    bool                    `json:"online"`
	Storage   AttributesStorageHosted `json:"storage"`
	Cleanup   AttributesCleanupPolicy `json:"cleanup"`
	Component AttributesComponent     `json:"component"`
}

type RepositoryAptHosted struct {
	Name       string                           `json:"name"`
	Online     bool                             `json:"online"`
	Storage    AttributesStorageHosted          `json:"storage"`
	Cleanup    AttributesCleanupPolicy          `json:"cleanup"`
	Component  AttributesComponent              `json:"component"`
	Apt        AttributesRepositoriesAptHosted  `json:"apt"`
	AptSigning AttributesRepositoriesAptSigning `json:"aptSigning"`
}

type RepositoryRawHosted struct {
	Name    string                  `json:"name"`
	Online  bool                    `json:"online"`
	Storage AttributesStorageHosted `json:"storage"`
	Cleanup AttributesCleanupPolicy `json:"cleanup"`
	Raw     AttributesRaw           `json:"raw"`
}

func CreateRepositoryHostedContext(ctx context.Context, rm RM, format repositoryFormat, r interface{}) error {
	buf, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not marshal: %v", err)
	}

	var restEndpointRepository string
	switch format {
	case Apt:
		restEndpointRepository = restRepositoriesHostedApt
	case Bower:
		restEndpointRepository = restRepositoriesHostedBower
	case Docker:
		restEndpointRepository = restRepositoriesHostedDocker
	case GitLfs:
		restEndpointRepository = restRepositoriesHostedGitLfs
	case Helm:
		restEndpointRepository = restRepositoriesHostedHelm
	case Maven:
		restEndpointRepository = restRepositoriesHostedMaven
	case Npm:
		restEndpointRepository = restRepositoriesHostedNpm
	case Nuget:
		restEndpointRepository = restRepositoriesHostedNuget
	case Pypi:
		restEndpointRepository = restRepositoriesHostedPypi
	case R:
		restEndpointRepository = restRepositoriesHostedR
	case Raw:
		restEndpointRepository = restRepositoriesHostedRaw
	case Rubygems:
		restEndpointRepository = restRepositoriesHostedRubygems
	case Yum:
		restEndpointRepository = restRepositoriesHostedYum
	}
	_, resp, err := rm.Post(ctx, restEndpointRepository, bytes.NewBuffer(buf))
	if err != nil && resp == nil {
		return fmt.Errorf("could not create repository: %v", err)
	}

	return nil
}

func CreateRepositoryHosted(rm RM, format repositoryFormat, r interface{}) error {
	return CreateRepositoryHostedContext(context.Background(), rm, format, r)
}

func CreateRepositoryProxyContext(ctx context.Context, rm RM, format repositoryFormat, r interface{}) error {
	buf, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not marshal: %v", err)
	}

	var restEndpointRepository string
	switch format {
	case Apt:
		restEndpointRepository = restRepositoriesProxyApt
	case Bower:
		restEndpointRepository = restRepositoriesProxyBower
	case Cocoapods:
		restEndpointRepository = restRepositoriesProxyCocoapods
	case Conan:
		restEndpointRepository = restRepositoriesProxyConan
	case Conda:
		restEndpointRepository = restRepositoriesProxyConda
	case Docker:
		restEndpointRepository = restRepositoriesProxyDocker
	case Golang:
		restEndpointRepository = restRepositoriesProxyGolang
	case Helm:
		restEndpointRepository = restRepositoriesProxyHelm
	case Maven:
		restEndpointRepository = restRepositoriesProxyMaven
	case Npm:
		restEndpointRepository = restRepositoriesProxyNpm
	case Nuget:
		restEndpointRepository = restRepositoriesProxyNuget
	case P2:
		restEndpointRepository = restRepositoriesProxyP2
	case Pypi:
		restEndpointRepository = restRepositoriesProxyPypi
	case R:
		restEndpointRepository = restRepositoriesProxyR
	case Raw:
		restEndpointRepository = restRepositoriesProxyRaw
	case Rubygems:
		restEndpointRepository = restRepositoriesProxyRubygems
	case Yum:
		restEndpointRepository = restRepositoriesProxyYum
	}
	_, resp, err := rm.Post(ctx, restEndpointRepository, bytes.NewBuffer(buf))
	if err != nil && resp == nil {
		return fmt.Errorf("could not create repository: %v", err)
	}

	return nil
}

func CreateRepositoryProxy(rm RM, format repositoryFormat, r interface{}) error {
	return CreateRepositoryProxyContext(context.Background(), rm, format, r)
}

func DeleteRepositoryByNameContext(ctx context.Context, rm RM, name string) error {
	url := fmt.Sprintf("%s/%s", restRepositories, name)

	if resp, err := rm.Del(ctx, url); err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("repository not deleted '%s': %v", name, err)
	}

	return nil
}

func DeleteRepositoryByName(rm RM, name string) error {
	return DeleteRepositoryByNameContext(context.Background(), rm, name)
}

func GetRepositoriesContext(ctx context.Context, rm RM) ([]Repository, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not find repositories: %v", err)
	}

	body, resp, err := rm.Get(ctx, restRepositories)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, doError(err)
	}

	repos := make([]Repository, 0)
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, doError(err)
	}

	return repos, nil
}

// GetRepositories returns a list of components in the indicated repository
func GetRepositories(rm RM) ([]Repository, error) {
	return GetRepositoriesContext(context.Background(), rm)
}

func GetRepositoryByNameContext(ctx context.Context, rm RM, name string) (repo Repository, err error) {
	repos, err := GetRepositoriesContext(ctx, rm)
	if err != nil {
		return repo, fmt.Errorf("could not get list of repositories: %v", err)
	}

	for _, repo = range repos {
		if repo.Name == name {
			return
		}
	}

	return repo, fmt.Errorf("did not find repository '%s': %v", name, err)
}

// GetRepositoryByName returns information on a named repository
func GetRepositoryByName(rm RM, name string) (repo Repository, err error) {
	return GetRepositoryByNameContext(context.Background(), rm, name)
}
