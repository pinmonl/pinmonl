package main

import (
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/bitbucket"
	"github.com/pinmonl/pinmonl/monler/docker"
	"github.com/pinmonl/pinmonl/monler/github"
	"github.com/pinmonl/pinmonl/monler/gitlab"
	"github.com/pinmonl/pinmonl/monler/helm"
	"github.com/pinmonl/pinmonl/monler/npm"
)

func initMonler(cfg *config.Config, ss stores) *monler.Repository {
	repo, err := monler.NewRepository()
	if err != nil {
		logx.Panic(err)
	}

	if githubPrd, err := newGithubProvider(cfg.Github.Token); err == nil {
		repo.Register(githubPrd)
	}
	if gitlabPrd, err := newGitlabProvider(); err == nil {
		repo.Register(gitlabPrd)
	}
	if bitbucketPrd, err := newBitbucketProvider(); err == nil {
		repo.Register(bitbucketPrd)
	}
	if npmPrd, err := newNpmProvider(); err == nil {
		repo.Register(npmPrd)
	}
	if dockerPrd, err := newDockerProvider(); err == nil {
		repo.Register(dockerPrd)
	}
	if helmPrd, err := newHelmProvider(); err == nil {
		repo.Register(helmPrd)
	}

	return repo
}

func newGithubProvider(token string) (*github.Provider, error) {
	opts := github.ProviderOpts{
		Token: token,
	}
	return github.NewProvider(&opts)
}

func newGitlabProvider() (*gitlab.Provider, error) {
	return gitlab.NewProvider(&gitlab.ProviderOpts{})
}

func newBitbucketProvider() (*bitbucket.Provider, error) {
	return bitbucket.NewProvider(&bitbucket.ProviderOpts{})
}

func newNpmProvider() (*npm.Provider, error) {
	return npm.NewProvider(&npm.ProviderOpts{})
}

func newDockerProvider() (*docker.Provider, error) {
	return docker.NewProvider(&docker.ProviderOpts{})
}

func newHelmProvider() (*helm.Provider, error) {
	return helm.NewProvider(&helm.ProviderOpts{})
}
