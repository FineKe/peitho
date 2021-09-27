#!/bin/bash

# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


PEITHO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${PEITHO_ROOT}/scripts/lib/init.sh"

if [ $# -ne 1 ];then
  iam::log::error "Usage: force_release.sh v1.0.0"
  exit 1  
fi

version="$1"

set +o errexit
# 1. delete old version
git tag -d ${version}
git push origin --delete ${version}

# 2. create a new tag
git tag -a ${version} -m "release ${version}"
git push origin master
git push origin ${version}

# 3. release the new release
pushd ${PEITHO_ROOT}
# try to delete target github release if exist to avoid create error    
iam::log::info "delete github release with tag ${IAM_GIT_VERSION} if exist"    
github-release delete  \
  --user tianrandailove\
  --repo Peitho  \
  --tag ${version}

make release
