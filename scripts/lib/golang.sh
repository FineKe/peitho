#!/usr/bin/env bash

# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# shellcheck disable=SC2034 # Variables sourced in other scripts.

# The server platform we are building on.
readonly PEITHO_SUPPORTED_SERVER_PLATFORMS=(
  linux/amd64
  linux/arm64
)

# If we update this we should also update the set of platforms whose standard
# library is precompiled for in build/build-image/cross/Dockerfile
readonly PEITHO_SUPPORTED_CLIENT_PLATFORMS=(
  linux/amd64
  linux/arm64
)

# The set of server targets that we are only building for Linux
# If you update this list, please also update build/BUILD.
peitho::golang::server_targets() {
  local targets=(
    peitho
  )
  echo "${targets[@]}"
}

IFS=" " read -ra PEITHO_SERVER_TARGETS <<< "$(peitho::golang::server_targets)"
readonly PEITHO_SERVER_TARGETS
readonly PEITHO_SERVER_BINARIES=("${PEITHO_SERVER_TARGETS[@]##*/}")

# The set of server targets we build docker images for
peitho::golang::server_image_targets() {
  # NOTE: this contains cmd targets for peitho::build::get_docker_wrapped_binaries
  local targets=(
    cmd/peitho
  )
  echo "${targets[@]}"
}

IFS=" " read -ra PEITHO_SERVER_IMAGE_TARGETS <<< "$(peitho::golang::server_image_targets)"
readonly PEITHO_SERVER_IMAGE_TARGETS
readonly PEITHO_SERVER_IMAGE_BINARIES=("${PEITHO_SERVER_IMAGE_TARGETS[@]##*/}")

# ------------
# NOTE: All functions that return lists should use newlines.
# bash functions can't return arrays, and spaces are tricky, so newline
# separators are the preferred pattern.
# To transform a string of newline-separated items to an array, use peitho::util::read-array:
# peitho::util::read-array FOO < <(peitho::golang::dups a b c a)
#
# ALWAYS remember to quote your subshells. Not doing so will break in
# bash 4.3, and potentially cause other issues.
# ------------

# Returns a sorted newline-separated list containing only duplicated items.
peitho::golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
peitho::golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

# Depends on values of user-facing PEITHO_BUILD_PLATFORMS, PEITHO_FASTBUILD,
# and PEITHO_BUILDER_OS.
# Configures PEITHO_SERVER_PLATFORMS, then sets them
# to readonly.
# The configured vars will only contain platforms allowed by the
# PEITHO_SUPPORTED* vars at the top of this file.
declare -a PEITHO_SERVER_PLATFORMS
peitho::golang::setup_platforms() {
  if [[ -n "${PEITHO_BUILD_PLATFORMS:-}" ]]; then
    # PEITHO_BUILD_PLATFORMS needs to be read into an array before the next
    # step, or quoting treats it all as one element.
    local -a platforms
    IFS=" " read -ra platforms <<< "${PEITHO_BUILD_PLATFORMS}"

    # Deduplicate to ensure the intersection trick with peitho::golang::dups
    # is not defeated by duplicates in user input.
    peitho::util::read-array platforms < <(peitho::golang::dedup "${platforms[@]}")

    # Use peitho::golang::dups to restrict the builds to the platforms in
    # PEITHO_SUPPORTED_*_PLATFORMS. Items should only appear at most once in each
    # set, so if they appear twice after the merge they are in the intersection.
    peitho::util::read-array PEITHO_SERVER_PLATFORMS < <(peitho::golang::dups \
        "${platforms[@]}" \
        "${PEITHO_SUPPORTED_SERVER_PLATFORMS[@]}" \
      )
    readonly PEITHO_SERVER_PLATFORMS

    peitho::util::read-array PEITHO_CLIENT_PLATFORMS < <(peitho::golang::dups \
        "${platforms[@]}" \
        "${PEITHO_SUPPORTED_CLIENT_PLATFORMS[@]}" \
      )
    readonly PEITHO_CLIENT_PLATFORMS

  elif [[ "${PEITHO_FASTBUILD:-}" == "true" ]]; then
    PEITHO_SERVER_PLATFORMS=(linux/amd64)
    readonly PEITHO_SERVER_PLATFORMS
  else
    PEITHO_SERVER_PLATFORMS=("${PEITHO_SUPPORTED_SERVER_PLATFORMS[@]}")
    readonly PEITHO_SERVER_PLATFORMS
  fi
}

peitho::golang::setup_platforms

readonly PEITHO_ALL_TARGETS=(
  "${PEITHO_SERVER_TARGETS[@]}"
)
readonly PEITHO_ALL_BINARIES=("${PEITHO_ALL_TARGETS[@]##*/}")

# Asks golang what it thinks the host platform is. The go tool chain does some
# slightly different things when the target platform matches the host platform.
peitho::golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}

# Ensure the go tool exists and is a viable version.
peitho::golang::verify_go_version() {
  if [[ -z "$(command -v go)" ]]; then
    peitho::log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<< "$(go version)"
  local minimum_go_version
  minimum_go_version=go1.13.4
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    peitho::log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
PEITHO requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

# peitho::golang::setup_env will check that the `go` commands is available in
# ${PATH}. It will also check that the Go version is good enough for the
# PEITHO build.
#
# Outputs:
#   env-var GOBIN is unset (we want binaries in a predictable place)
#   env-var GO15VENDOREXPERIMENT=1
#   env-var GO111MODULE=on
peitho::golang::setup_env() {
  peitho::golang::verify_go_version

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This seems to matter to some tools
  export GO15VENDOREXPERIMENT=1

  # Open go module feature
  export GO111MODULE=on

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  umask 0022
}
