// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"regexp"

	"github.com/tianrandailove/peitho/pkg/log"
)

// IsContainerID check the containerID is universal containerID
// if it is chaincodeID, like "dev.peer0.org" will return false
func IsContainerID(id string) bool {
	result, err := regexp.MatchString("^[0-9a-f]{64}", id)
	if err != nil {
		log.Errorf("%v", err)

		return false
	}

	return result
}
