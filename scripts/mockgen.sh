#!/usr/bin/env bash

mockgen_cmd="mockgen"

$mockgen_cmd -source=x/distribution/types/expected_keepers.go -package testutil -destination x/distribution/testutil/expected_keepers_mocks.go
