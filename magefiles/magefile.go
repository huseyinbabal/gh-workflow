package main

import (
	"botkub-cloud-backend/tools/deployment"
	"errors"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/carolynvs/magex/mgx"
	"github.com/magefile/mage/mg"
	"go.szostok.io/magex/shx"
	"os"
)

// Test targets
type Test mg.Namespace

// Integration executes integration tests.
// TODO: race condition is out of context for now since we don't have a use case where we run
// http server in one app instance, but we have multiple running apps in e2e tests. Revisit
// this part once we want to parallelize tests.
func (Test) Integration() {
	mgx.Must(shx.MustCmdf("go test -tags=integration -race -p 1 -v ./e2e/...").RunV())
}

// Unit executes unit tests.
func (Test) Unit() {
	mgx.Must(shx.MustCmdf("go test -v  -race ./...").RunV())
}

// E2e executes end-to-end tests
func (Test) E2e() {
	mgx.Must(shx.MustCmdf("go test -tags=e2e -race -p 1 -v ./e2e/...").RunV())
}

// Slack_e2e executes end-to-end tests for Cloud Slack.
func (Test) Slack_e2e() {
	cmd := "go test -tags=slack_e2e -race -p 1 -v ./e2e/..."
	if os.Getenv("SHOW_BROWSER") == "true" {
		cmd += " -rod=show"
	}
	mgx.Must(shx.MustCmdf(cmd).RunV())
}

type Gen mg.Namespace

// Gql generates required Go source code for graphql schemas.
func (Gen) Gql() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	mgx.Must(err)

	err = api.Generate(cfg)
	mgx.Must(err)
}

type Deployment mg.Namespace

func (Deployment) Wait() {
	var err error
	url := os.Getenv("DEPLOYMENT_STATUS_URL")
	if url == "" {
		err = errors.New("DEPLOYMENT_STATUS_URL env variable is not set")
	}
	version := os.Getenv("EXPECTED_VERSION")
	if version == "" {
		err = errors.New("EXPECTED_VERSION env variable is not set")
	}
	mgx.Must(err)
	mgx.Must(deployment.WaitFor(url, version))
}
