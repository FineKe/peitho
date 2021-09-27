#!/usr/bin/env bash

# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


function peitho::util::sourced_variable {
  # Call this function to tell shellcheck that a variable is supposed to
  # be used from other calling context. This helps quiet an "unused
  # variable" warning from shellcheck and also document your code.
  true
}

peitho::util::sortable_date() {
  date "+%Y%m%d-%H%M%S"
}

# arguments: target, item1, item2, item3, ...
# returns 0 if target is in the given items, 1 otherwise.
peitho::util::array_contains() {
  local search="$1"
  local element
  shift
  for element; do
    if [[ "${element}" == "${search}" ]]; then
      return 0
     fi
  done
  return 1
}

peitho::util::wait_for_url() {
  local url=$1
  local prefix=${2:-}
  local wait=${3:-1}
  local times=${4:-30}
  local maxtime=${5:-1}

  command -v curl >/dev/null || {
    peitho::log::usage "curl must be installed"
    exit 1
  }

  local i
  for i in $(seq 1 "${times}"); do
    local out
    if out=$(curl --max-time "${maxtime}" -gkfs "${url}" 2>/dev/null); then
      peitho::log::status "On try ${i}, ${prefix}: ${out}"
      return 0
    fi
    sleep "${wait}"
  done
  peitho::log::error "Timed out waiting for ${prefix} to answer at ${url}; tried ${times} waiting ${wait} between each"
  return 1
}

# Example:  peitho::util::wait_for_success 120 5
# arguments: wait time, sleep time, shell command
# returns 0 if the shell command get output, 1 otherwise.
peitho::util::wait_for_success(){
  local wait_time="$1"
  local sleep_time="$2"
  local cmd="$3"
  while [ "$wait_time" -gt 0 ]; do
    if eval "$cmd"; then
      return 0
    else
      sleep "$sleep_time"
      wait_time=$((wait_time-sleep_time))
    fi
  done
  return 1
}

# Example:  peitho::util::trap_add 'echo "in trap DEBUG"' DEBUG
# See: http://stackoverflow.com/questions/3338030/multiple-bash-traps-for-the-same-signal
peitho::util::trap_add() {
  local trap_add_cmd
  trap_add_cmd=$1
  shift

  for trap_add_name in "$@"; do
    local existing_cmd
    local new_cmd

    # Grab the currently defined trap commands for this trap
    existing_cmd=$(trap -p "${trap_add_name}" |  awk -F"'" '{print $2}')

    if [[ -z "${existing_cmd}" ]]; then
      new_cmd="${trap_add_cmd}"
    else
      new_cmd="${trap_add_cmd};${existing_cmd}"
    fi

    # Assign the test. Disable the shellcheck warning telling that trap
    # commands should be single quoted to avoid evaluating them at this
    # point instead evaluating them at run time. The logic of adding new
    # commands to a single trap requires them to be evaluated right away.
    # shellcheck disable=SC2064
    trap "${new_cmd}" "${trap_add_name}"
  done
}

# Opposite of peitho::util::ensure-temp-dir()
peitho::util::cleanup-temp-dir() {
  rm -rf "${PEITHO_TEMP}"
}

# Create a temp dir that'll be deleted at the end of this bash session.
#
# Vars set:
#   PEITHO_TEMP
peitho::util::ensure-temp-dir() {
  if [[ -z ${PEITHO_TEMP-} ]]; then
    PEITHO_TEMP=$(mktemp -d 2>/dev/null || mktemp -d -t peithornetes.XXXXXX)
    peitho::util::trap_add peitho::util::cleanup-temp-dir EXIT
  fi
}

peitho::util::host_os() {
  local host_os
  case "$(uname -s)" in
    Darwin)
      host_os=darwin
      ;;
    Linux)
      host_os=linux
      ;;
    *)
      peitho::log::error "Unsupported host OS.  Must be Linux or Mac OS X."
      exit 1
      ;;
  esac
  echo "${host_os}"
}

peitho::util::host_arch() {
  local host_arch
  case "$(uname -m)" in
    x86_64*)
      host_arch=amd64
      ;;
    i?86_64*)
      host_arch=amd64
      ;;
    amd64*)
      host_arch=amd64
      ;;
    aarch64*)
      host_arch=arm64
      ;;
    arm64*)
      host_arch=arm64
      ;;
    arm*)
      host_arch=arm
      ;;
    i?86*)
      host_arch=x86
      ;;
    s390x*)
      host_arch=s390x
      ;;
    ppc64le*)
      host_arch=ppc64le
      ;;
    *)
      peitho::log::error "Unsupported host arch. Must be x86_64, 386, arm, arm64, s390x or ppc64le."
      exit 1
      ;;
  esac
  echo "${host_arch}"
}

# This figures out the host platform without relying on golang.  We need this as
# we don't want a golang install to be a prerequisite to building yet we need
# this info to figure out where the final binaries are placed.
peitho::util::host_platform() {
  echo "$(peitho::util::host_os)/$(peitho::util::host_arch)"
}

# looks for $1 in well-known output locations for the platform ($2)
# $PEITHO_ROOT must be set
peitho::util::find-binary-for-platform() {
  local -r lookfor="$1"
  local -r platform="$2"
  local locations=(
    "${PEITHO_ROOT}/_output/bin/${lookfor}"
    "${PEITHO_ROOT}/_output/${platform}/${lookfor}"
    "${PEITHO_ROOT}/_output/local/bin/${platform}/${lookfor}"
    "${PEITHO_ROOT}/_output/platforms/${platform}/${lookfor}"
  )

  # List most recently-updated location.
  local -r bin=$( (ls -t "${locations[@]}" 2>/dev/null || true) | head -1 )
  echo -n "${bin}"
}

# looks for $1 in well-known output locations for the host platform
# $PEITHO_ROOT must be set
peitho::util::find-binary() {
  peitho::util::find-binary-for-platform "$1" "$(peitho::util::host_platform)"
}


# Returns the name of the upstream remote repository name for the local git
# repo, e.g. "upstream" or "origin".
peitho::util::git_upstream_remote_name() {
  git remote -v | grep fetch |\
    grep -E 'github.com[/:]tiranrandailove/peitho' |\
    head -n 1 | awk '{print $1}'
}

# Exits script if working directory is dirty. If it's run interactively in the terminal
# the user can commit changes in a second terminal. This script will wait.
peitho::util::ensure_clean_working_dir() {
  while ! git diff HEAD --exit-code &>/dev/null; do
    echo -e "\nUnexpected dirty working directory:\n"
    if tty -s; then
        git status -s
    else
        git diff -a # be more verbose in log files without tty
        exit 1
    fi | sed 's/^/  /'
    echo -e "\nCommit your changes in another terminal and then continue here by pressing enter."
    read -r
  done 1>&2
}

# Find the base commit using:
# $PULL_BASE_SHA if set (from Prow)
# current ref from the remote upstream branch
peitho::util::base_ref() {
  local -r git_branch=$1

  if [[ -n ${PULL_BASE_SHA:-} ]]; then
    echo "${PULL_BASE_SHA}"
    return
  fi

  full_branch="$(peitho::util::git_upstream_remote_name)/${git_branch}"

  # make sure the branch is valid, otherwise the check will pass erroneously.
  if ! git describe "${full_branch}" >/dev/null; then
    # abort!
    exit 1
  fi

  echo "${full_branch}"
}

# Checks whether there are any files matching pattern $2 changed between the
# current branch and upstream branch named by $1.
# Returns 1 (false) if there are no changes
#         0 (true) if there are changes detected.
peitho::util::has_changes() {
  local -r git_branch=$1
  local -r pattern=$2
  local -r not_pattern=${3:-totallyimpossiblepattern}

  local base_ref
  base_ref=$(peitho::util::base_ref "${git_branch}")
  echo "Checking for '${pattern}' changes against '${base_ref}'"

  # notice this uses ... to find the first shared ancestor
  if git diff --name-only "${base_ref}...HEAD" | grep -v -E "${not_pattern}" | grep "${pattern}" > /dev/null; then
    return 0
  fi
  # also check for pending changes
  if git status --porcelain | grep -v -E "${not_pattern}" | grep "${pattern}" > /dev/null; then
    echo "Detected '${pattern}' uncommitted changes."
    return 0
  fi
  echo "No '${pattern}' changes detected."
  return 1
}

peitho::util::download_file() {
  local -r url=$1
  local -r destination_file=$2

  rm "${destination_file}" 2&> /dev/null || true

  for i in $(seq 5)
  do
    if ! curl -fsSL --retry 3 --keepalive-time 2 "${url}" -o "${destination_file}"; then
      echo "Downloading ${url} failed. $((5-i)) retries left."
      sleep 1
    else
      echo "Downloading ${url} succeed"
      return 0
    fi
  done
  return 1
}



# Determines if docker can be run, failures may simply require that the user be added to the docker group.
function peitho::util::ensure_docker_daemon_connectivity {
  IFS=" " read -ra DOCKER <<< "${DOCKER_OPTS}"
  # Expand ${DOCKER[@]} only if it's not unset. This is to work around
  # Bash 3 issue with unbound variable.
  DOCKER=(docker ${DOCKER[@]:+"${DOCKER[@]}"})
  if ! "${DOCKER[@]}" info > /dev/null 2>&1 ; then
    cat <<'EOF' >&2
Can't connect to 'docker' daemon.  please fix and retry.

Possible causes:
  - Docker Daemon not started
    - Linux: confirm via your init system
    - macOS w/ docker-machine: run `docker-machine ls` and `docker-machine start <name>`
    - macOS w/ Docker for Mac: Check the menu bar and start the Docker application
  - DOCKER_HOST hasn't been set or is set incorrectly
    - Linux: domain socket is used, DOCKER_* should be unset. In Bash run `unset ${!DOCKER_*}`
    - macOS w/ docker-machine: run `eval "$(docker-machine env <name>)"`
    - macOS w/ Docker for Mac: domain socket is used, DOCKER_* should be unset. In Bash run `unset ${!DOCKER_*}`
  - Other things to check:
    - Linux: User isn't in 'docker' group.  Add and relogin.
      - Something like 'sudo usermod -a -G docker ${USER}'
      - RHEL7 bug and workaround: https://bugzilla.redhat.com/show_bug.cgi?id=1119282#c8
EOF
    return 1
  fi
}

# Wait for background jobs to finish. Return with
# an error status if any of the jobs failed.
peitho::util::wait-for-jobs() {
  local fail=0
  local job
  for job in $(jobs -p); do
    wait "${job}" || fail=$((fail + 1))
  done
  return ${fail}
}

# peitho::util::join <delim> <list...>
# Concatenates the list elements with the delimiter passed as first parameter
#
# Ex: peitho::util::join , a b c
#  -> a,b,c
function peitho::util::join {
  local IFS="$1"
  shift
  echo "$*"
}



# peitho::util::check-file-in-alphabetical-order <file>
# Check that the file is in alphabetical order
#
function peitho::util::check-file-in-alphabetical-order {
  local failure_file="$1"
  if ! diff -u "${failure_file}" <(LC_ALL=C sort "${failure_file}"); then
    {
      echo
      echo "${failure_file} is not in alphabetical order. Please sort it:"
      echo
      echo "  LC_ALL=C sort -o ${failure_file} ${failure_file}"
      echo
    } >&2
  false
  fi
}

# peitho::util::require-jq
# Checks whether jq is installed.
function peitho::util::require-jq {
  if ! command -v jq &>/dev/null; then
    echo "jq not found. Please install." 1>&2
    return 1
  fi
}

# outputs md5 hash of $1, works on macOS and Linux
function peitho::util::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

# peitho::util::read-array
# Reads in stdin and adds it line by line to the array provided. This can be
# used instead of "mapfile -t", and is bash 3 compatible.
#
# Assumed vars:
#   $1 (name of array to create/modify)
#
# Example usage:
# peitho::util::read-array files < <(ls -1)
#
function peitho::util::read-array {
  local i=0
  unset -v "$1"
  while IFS= read -r "$1[i++]"; do :; done
  eval "[[ \${$1[--i]} ]]" || unset "$1[i]" # ensures last element isn't empty
}

# Some useful colors.
if [[ -z "${color_start-}" ]]; then
  declare -r color_start="\033["
  declare -r color_red="${color_start}0;31m"
  declare -r color_yellow="${color_start}0;33m"
  declare -r color_green="${color_start}0;32m"
  declare -r color_blue="${color_start}1;34m"
  declare -r color_cyan="${color_start}1;36m"
  declare -r color_norm="${color_start}0m"

  peitho::util::sourced_variable "${color_start}"
  peitho::util::sourced_variable "${color_red}"
  peitho::util::sourced_variable "${color_yellow}"
  peitho::util::sourced_variable "${color_green}"
  peitho::util::sourced_variable "${color_blue}"
  peitho::util::sourced_variable "${color_cyan}"
  peitho::util::sourced_variable "${color_norm}"
fi

# ex: ts=2 sw=2 et filetype=sh
