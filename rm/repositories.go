package nexusrm

import (
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
	Maven
	Npm
	Nuget
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
	Bower
	Pypi
	Yum
	GitLfs
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

// GetRepositories returns a list of components in the indicated repository
func GetRepositories(rm RM) ([]Repository, error) {
	doError := func(err error) error {
		return fmt.Errorf("could not find repositories: %v", err)
	}

	body, resp, err := rm.Get(restRepositories)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, doError(err)
	}

	repos := make([]Repository, 0)
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, doError(err)
	}

	return repos, nil
}

// GetRepositoryByName returns information on a named repository
func GetRepositoryByName(rm RM, name string) (repo Repository, err error) {
	repos, err := GetRepositories(rm)
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
