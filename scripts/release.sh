#!/usr/bin/env bash

# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.



# Build a peitho release.  This will build the binaries, create the Docker
# images and other build artifacts.

set -o errexit
set -o nounset
set -o pipefail

PEITHO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${PEITHO_ROOT}/scripts/common.sh"
source "${PEITHO_ROOT}/scripts/lib/release.sh"

PEITHO_RELEASE_RUN_TESTS=${PEITHO_RELEASE_RUN_TESTS-y}

peitho::golang::setup_env
peitho::build::verify_prereqs
peitho::release::verify_prereqs
#peitho::build::build_image
peitho::build::build_command
peitho::release::package_tarballs
#peitho::release::updload_tarballs
peitho::release::github_release
peitho::release::generate_changelog
