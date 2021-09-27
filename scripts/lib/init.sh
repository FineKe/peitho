#!/usr/bin/env bash

# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.



set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
PEITHO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${PEITHO_ROOT}/scripts/lib/util.sh"
source "${PEITHO_ROOT}/scripts/lib/logging.sh"
source "${PEITHO_ROOT}/scripts/lib/color.sh"

peitho::log::install_errexit

source "${PEITHO_ROOT}/scripts/lib/version.sh"
source "${PEITHO_ROOT}/scripts/lib/golang.sh"
