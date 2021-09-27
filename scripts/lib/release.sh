#!/usr/bin/env bash

# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# This file creates release artifacts (tar files, container images) that are
# ready to distribute to install or distribute to end users.

###############################################################################
# Most of the ::release:: namespace functions have been moved to
# github.com/tianrandailove/peitho/release.  Have a look in that repo and specifically in
# lib/releaselib.sh for ::release::-related functionality.
###############################################################################

# This is where the final release artifacts are created locally
readonly RELEASE_STAGE="${LOCAL_OUTPUT_ROOT}/release-stage"
readonly RELEASE_TARS="${LOCAL_OUTPUT_ROOT}/release-tars"
readonly RELEASE_IMAGES="${LOCAL_OUTPUT_ROOT}/release-images"

# peitho github account info
readonly PEITHO_GITHUB_ORG=tianrandailove
readonly PEITHO_GITHUB_REPO=Peitho

readonly ARTIFACT=peitho.tar.gz
readonly CHECKSUM=${ARTIFACT}.sha1sum

PEITHO_BUILD_CONFORMANCE=${PEITHO_BUILD_CONFORMANCE:-y}
PEITHO_BUILD_PULL_LATEST_IMAGES=${PEITHO_BUILD_PULL_LATEST_IMAGES:-y}

# Validate a ci version
#
# Globals:
#   None
# Arguments:
#   version
# Returns:
#   If version is a valid ci version
# Sets:                    (e.g. for '1.2.3-alpha.4.56+abcdef12345678')
#   VERSION_MAJOR          (e.g. '1')
#   VERSION_MINOR          (e.g. '2')
#   VERSION_PATCH          (e.g. '3')
#   VERSION_PRERELEASE     (e.g. 'alpha')
#   VERSION_PRERELEASE_REV (e.g. '4')
#   VERSION_BUILD_INFO     (e.g. '.56+abcdef12345678')
#   VERSION_COMMITS        (e.g. '56')
function peitho::release::parse_and_validate_ci_version() {
  # Accept things like "v1.2.3-alpha.4.56+abcdef12345678" or "v1.2.3-beta.4"
  local -r version_regex="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)-([a-zA-Z0-9]+)\\.(0|[1-9][0-9]*)(\\.(0|[1-9][0-9]*)\\+[0-9a-f]{7,40})?$"
  local -r version="${1-}"
  [[ "${version}" =~ ${version_regex} ]] || {
    peitho::log::error "Invalid ci version: '${version}', must match regex ${version_regex}"
    return 1
  }

  # The VERSION variables are used when this file is sourced, hence
  # the shellcheck SC2034 'appears unused' warning is to be ignored.

  # shellcheck disable=SC2034
  VERSION_MAJOR="${BASH_REMATCH[1]}"
  # shellcheck disable=SC2034
  VERSION_MINOR="${BASH_REMATCH[2]}"
  # shellcheck disable=SC2034
  VERSION_PATCH="${BASH_REMATCH[3]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE="${BASH_REMATCH[4]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE_REV="${BASH_REMATCH[5]}"
  # shellcheck disable=SC2034
  VERSION_BUILD_INFO="${BASH_REMATCH[6]}"
  # shellcheck disable=SC2034
  VERSION_COMMITS="${BASH_REMATCH[7]}"
}

# ---------------------------------------------------------------------------
# Build final release artifacts
function peitho::release::clean_cruft() {
  # Clean out cruft
  find "${RELEASE_STAGE}" -name '*~' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '#*#' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '.DS*' -exec rm {} \;
}

function peitho::release::package_tarballs() {
  # Clean out any old releases
  rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"
  mkdir -p "${RELEASE_TARS}"
  peitho::release::package_src_tarball &
  peitho::util::wait-for-jobs || { peitho::log::error "previous tarball phase failed"; return 1; }
}

# Package the source code we built, compliance/licensing/audit/yaddafor .
function peitho::release::package_src_tarball() {
  local -r src_tarball="${RELEASE_TARS}/peitho-src.tar.gz"
  peitho::log::status "Building tarball: src"
  if [[ "${PEITHO_GIT_TREE_STATE-}" = 'clean' ]]; then
    git archive -o "${src_tarball}" HEAD
  else
    find "${PEITHO_ROOT}" -mindepth 1 -maxdepth 1 \
      ! \( \
      \( -path "${PEITHO_ROOT}"/_\* -o \
      -path "${PEITHO_ROOT}"/.git\* -o \
      -path "${PEITHO_ROOT}"/.gitignore\* -o \
      -path "${PEITHO_ROOT}"/.gsemver.yaml\* -o \
      -path "${PEITHO_ROOT}"/.config\* -o \
      -path "${PEITHO_ROOT}"/.chglog\* -o \
      -path "${PEITHO_ROOT}"/.gitlint -o \
      -path "${PEITHO_ROOT}"/.golangci.yaml -o \
      -path "${PEITHO_ROOT}"/.goreleaser.yml -o \
      -path "${PEITHO_ROOT}"/.note.md -o \
      -path "${PEITHO_ROOT}"/.todo.md \
      \) -prune \
      \) -print0 \
      | "${TAR}" czf "${src_tarball}" --transform "s|${PEITHO_ROOT#/*}|peitho|" --null -T -
  fi
}




# Package up all of the server binaries in docker images
function peitho::release::build_server_images() {
  # Clean out any old images
  rm -rf "${RELEASE_IMAGES}"
  local platform
  for platform in "${PEITHO_SERVER_PLATFORMS[@]}"; do
    local platform_tag
    local arch
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    arch=$(basename "${platform}")
    peitho::log::status "Building images: $platform_tag"

    local release_stage
    release_stage="${RELEASE_STAGE}/server/${platform_tag}/peitho"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    # This fancy expression will expand to prepend a path
    # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
    # PEITHO_SERVER_IMAGE_BINARIES array.
    cp "${PEITHO_SERVER_IMAGE_BINARIES[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
      "${release_stage}/server/bin/"

    peitho::release::create_docker_images_for_server "${release_stage}/server/bin" "${arch}"
  done
}

function peitho::release::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

function peitho::release::sha1() {
  if which sha1sum >/dev/null 2>&1; then
    sha1sum "$1" | awk '{ print $1 }'
  else
    shasum -a1 "$1" | awk '{ print $1 }'
  fi
}

function peitho::release::build_conformance_image() {
  local -r arch="$1"
  local -r registry="$2"
  local -r version="$3"
  local -r save_dir="${4-}"
  peitho::log::status "Building conformance image for arch: ${arch}"
  ARCH="${arch}" REGISTRY="${registry}" VERSION="${version}" \
    make -C cluster/images/conformance/ build >/dev/null

  local conformance_tag
  conformance_tag="${registry}/conformance-${arch}:${version}"
  if [[ -n "${save_dir}" ]]; then
    "${DOCKER[@]}" save "${conformance_tag}" > "${save_dir}/conformance-${arch}.tar"
  fi
  peitho::log::status "Deleting conformance image ${conformance_tag}"
  "${DOCKER[@]}" rmi "${conformance_tag}" &>/dev/null || true
}

# This builds all the release docker images (One docker image per binary)
# Args:
#  $1 - binary_dir, the directory to save the tared images to.
#  $2 - arch, architecture for which we are building docker images.
function peitho::release::create_docker_images_for_server() {
  # Create a sub-shell so that we don't pollute the outer environment
  (
    local binary_dir
    local arch
    local binaries
    local images_dir
    binary_dir="$1"
    arch="$2"
    binaries=$(peitho::build::get_docker_wrapped_binaries "${arch}")
    images_dir="${RELEASE_IMAGES}/${arch}"
    mkdir -p "${images_dir}"

    # k8s.gcr.io is the constant tag in the docker archives, this is also the default for config scripts in GKE.
    # We can use PEITHO_DOCKER_REGISTRY to include and extra registry in the docker archive.
    # If we use PEITHO_DOCKER_REGISTRY="k8s.gcr.io", then the extra tag (same) is ignored, see release_docker_image_tag below.
    local -r docker_registry="k8s.gcr.io"
    # Docker tags cannot contain '+'
    local docker_tag="${PEITHO_GIT_VERSION/+/_}"
    if [[ -z "${docker_tag}" ]]; then
      peitho::log::error "git version information missing; cannot create Docker tag"
      return 1
    fi

    # provide `--pull` argument to `docker build` if `PEITHO_BUILD_PULL_LATEST_IMAGES`
    # is set to y or Y; otherwise try to build the image without forcefully
    # pulling the latest base image.
    local docker_build_opts
    docker_build_opts=
    if [[ "${PEITHO_BUILD_PULL_LATEST_IMAGES}" =~ [yY] ]]; then
        docker_build_opts='--pull'
    fi

    for wrappable in $binaries; do

      local binary_name=${wrappable%%,*}
      local base_image=${wrappable##*,}
      local binary_file_path="${binary_dir}/${binary_name}"
      local docker_build_path="${binary_file_path}.dockerbuild"
      local docker_file_path="${docker_build_path}/Dockerfile"
      local docker_image_tag="${docker_registry}/${binary_name}-${arch}:${docker_tag}"

      peitho::log::status "Starting docker build for image: ${binary_name}-${arch}"
      (
        rm -rf "${docker_build_path}"
        mkdir -p "${docker_build_path}"
        ln "${binary_file_path}" "${docker_build_path}/${binary_name}"
        ln "${PEITHO_ROOT}/build/nsswitch.conf" "${docker_build_path}/nsswitch.conf"
        chmod 0644 "${docker_build_path}/nsswitch.conf"
        cat <<EOF > "${docker_file_path}"
FROM ${base_image}
COPY ${binary_name} /usr/local/bin/${binary_name}
EOF
        # ensure /etc/nsswitch.conf exists so go's resolver respects /etc/hosts
        if [[ "${base_image}" =~ busybox ]]; then
          echo "COPY nsswitch.conf /etc/" >> "${docker_file_path}"
        fi

        "${DOCKER[@]}" build ${docker_build_opts:+"${docker_build_opts}"} -q -t "${docker_image_tag}" "${docker_build_path}" >/dev/null
        # If we are building an official/alpha/beta release we want to keep
        # docker images and tag them appropriately.
        local -r release_docker_image_tag="${PEITHO_DOCKER_REGISTRY-$docker_registry}/${binary_name}-${arch}:${PEITHO_DOCKER_IMAGE_TAG-$docker_tag}"
        if [[ "${release_docker_image_tag}" != "${docker_image_tag}" ]]; then
          peitho::log::status "Tagging docker image ${docker_image_tag} as ${release_docker_image_tag}"
          "${DOCKER[@]}" rmi "${release_docker_image_tag}" 2>/dev/null || true
          "${DOCKER[@]}" tag "${docker_image_tag}" "${release_docker_image_tag}" 2>/dev/null
        fi
        "${DOCKER[@]}" save -o "${binary_file_path}.tar" "${docker_image_tag}" "${release_docker_image_tag}"
        echo "${docker_tag}" > "${binary_file_path}.docker_tag"
        rm -rf "${docker_build_path}"
        ln "${binary_file_path}.tar" "${images_dir}/"

        peitho::log::status "Deleting docker image ${docker_image_tag}"
        "${DOCKER[@]}" rmi "${docker_image_tag}" &>/dev/null || true
      ) &
    done

    if [[ "${PEITHO_BUILD_CONFORMANCE}" =~ [yY] ]]; then
      peitho::release::build_conformance_image "${arch}" "${docker_registry}" \
        "${docker_tag}" "${images_dir}" &
    fi

    peitho::util::wait-for-jobs || { peitho::log::error "previous Docker build failed"; return 1; }
    peitho::log::status "Docker builds done"
  )

}

# Build a release tarball.  $1 is the output tar name.  $2 is the base directory
# of the files to be packaged.  This assumes that ${2}/PEITHOis what is
# being packaged.
function peitho::release::create_tarball() {
  peitho::build::ensure_tar

  local tarfile=$1
  local stagingdir=$2

  "${TAR}" czf "${tarfile}" -C "${stagingdir}" peitho --owner=0 --group=0
}

function peitho::release::install_github_release(){
  GO111MODULE=off go get -u github.com/github-release/github-release
}

# Require the following tools:
# - github-release
# - gsemver
# - git-chglog
# - coscmd
function peitho::release::verify_prereqs(){
  if [ -z "$(which github-release 2>/dev/null)" ]; then
    peitho::log::info "'github-release' tool not installed, try to install it."

    if ! peitho::release::install_github_release; then
      peitho::log::error "failed to install 'github-release'"
      return 1
    fi
  fi

  if [ -z "$(which git-chglog 2>/dev/null)" ]; then
    peitho::log::info "'git-chglog' tool not installed, try to install it."

    if ! go get github.com/git-chglog/git-chglog/cmd/git-chglog &>/dev/null; then
      peitho::log::error "failed to install 'git-chglog'"
      return 1
    fi
  fi

  if [ -z "$(which gsemver 2>/dev/null)" ]; then
    peitho::log::info "'gsemver' tool not installed, try to install it."

    if ! go get github.com/arnaud-deprez/gsemver &>/dev/null; then
      peitho::log::error "failed to install 'gsemver'"
      return 1
    fi
  fi
}

# Create a github release with specified tarballs.
# NOTICE: Must export 'GITHUB_TOKEN' env in the shell, details:
# https://github.com/github-release/github-release
function peitho::release::github_release() {
  # create a github release
  peitho::log::info "create a new github release with tag ${PEITHO_GIT_VERSION}"
  github-release release \
    --user ${PEITHO_GITHUB_ORG} \
    --repo ${PEITHO_GITHUB_REPO} \
    --tag ${PEITHO_GIT_VERSION} \
    --description "" \
    --pre-release

  peitho::log::info "upload peitho-src.tar.gz to release ${PEITHO_GIT_VERSION}"
  github-release upload \
    --user ${PEITHO_GITHUB_ORG} \
    --repo ${PEITHO_GITHUB_REPO} \
    --tag ${PEITHO_GIT_VERSION} \
    --name "peitho-src.tar.gz" \
    --file ${RELEASE_TARS}/peitho-src.tar.gz
}

function peitho::release::generate_changelog() {
  peitho::log::info "generate CHANGELOG-${PEITHO_GIT_VERSION#v}.md and commit it"

  git-chglog ${PEITHO_GIT_VERSION} > ${PEITHO_ROOT}/CHANGELOG/CHANGELOG-${PEITHO_GIT_VERSION#v}.md

  (set +o errexit git add ${PEITHO_ROOT}/CHANGELOG/CHANGELOG-${PEITHO_GIT_VERSION#v}.md)
  git commit -a -m "docs(changelog): add CHANGELOG-${PEITHO_GIT_VERSION#v}.md"
  git push -f origin master # push CHANGELOG
}

