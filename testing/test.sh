#!/usr/bin/env bash

ENV_FILE=$1

function run_action() {
    rm -rf ./repo || true
    rm ./bin/speakeasy || true
    go run main.go
}

# Default environment variables not subject to change by different tests
export INPUT_DEBUG=true
export INPUT_OPENAPI_DOC_LOCATION="https://docs.speakeasyapi.dev/openapi.yaml"
export INPUT_GITHUB_ACCESS_TOKEN=${GITHUB_ACCESS_TOKEN}
export GITHUB_SERVER_URL="https://github.com"
export GITHUB_REPOSITORY_OWNER="speakeasy-api"
export GITHUB_REF="refs/heads/main"
export GITHUB_OUTPUT="./output.txt"
export GITHUB_WORKFLOW="test"

set -o allexport && source ${ENV_FILE} && set +o allexport

rm output.txt || true
run_action

if [ "$RUN_FINALIZE" = "true" ]; then
    BRANCH_NAME=$(go run testing/getoutput.go -output branch_name)
    PREVIOUS_GEN_VERSION=$(go run testing/getoutput.go -output previous_gen_version)
    export INPUT_BRANCH_NAME=${BRANCH_NAME}
    export INPUT_PREVIOUS_GEN_VERSION=${PREVIOUS_GEN_VERSION}
    INPUT_ACTION="finalize"
    run_action
fi