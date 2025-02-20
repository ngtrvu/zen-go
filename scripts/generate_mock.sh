#!/usr/bin/env bash

function generate_mock() {
    SOURCE=$1
    DESTINATION=$2
    echo 'Generating '$SOURCE
    mockgen -source=$SOURCE -destination=$DESTINATION
}

# repositories
generate_mock 'zen/repo.go' 'zen/mocks/repo_mock.go'
generate_mock 'zen/service.go' 'zen/mocks/service_mock.go'
generate_mock 'storage/cloud_storage.go' 'storage/mocks/cloud_storage_mock.go'
