// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

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

// GetDeploymentName get a valid deployment name
func GetDeploymentName(ID string) string {
	if len(ID) <= 53 {

		return strings.ReplaceAll(ID, ".", "-")
	}

	ID = strings.ReplaceAll(ID, ".", "-")
	hasher := md5.New()
	hasher.Write([]byte(ID))
	digest := hasher.Sum(nil)

	return fmt.Sprintf("chaincode-%s-%s", ID[:10], hex.EncodeToString(digest))
}

// GetDeploymentName get a valid deployment name
func GetDeploymentName(ID string) string {
	if len(ID) <= 53 {

		return strings.ReplaceAll(ID, ".", "-")
	}

	ID = strings.ReplaceAll(ID, ".", "-")
	hasher := md5.New()
	hasher.Write([]byte(ID))
	digest := hasher.Sum(nil)

	return fmt.Sprintf("chaincode-%s-%s", ID[:10], hex.EncodeToString(digest))
}

// GetDeploymentName get a valid deployment name
func GetDeploymentName(ID string) string {
	if len(ID) <= 53 {

		return strings.ReplaceAll(ID, ".", "-")
	}

	ID = strings.ReplaceAll(ID, ".", "-")
	hasher := md5.New()
	hasher.Write([]byte(ID))
	digest := hasher.Sum(nil)

	return fmt.Sprintf("chaincode-%s-%s", ID[:10], hex.EncodeToString(digest))
}

// GetDeploymentName get a valid deployment name
func GetDeploymentName(ID string) string {
	if len(ID) <= 53 {

		return strings.ReplaceAll(ID, ".", "-")
	}

	ID = strings.ReplaceAll(ID, ".", "-")
	hasher := md5.New()
	hasher.Write([]byte(ID))
	digest := hasher.Sum(nil)

	return fmt.Sprintf("chaincode-%s-%s", ID[:10], hex.EncodeToString(digest))
}
