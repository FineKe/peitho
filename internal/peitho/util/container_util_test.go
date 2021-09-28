// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestIsContainerID(t *testing.T) {
	containerID := "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2"
	want := true
	got := IsContainerID(containerID)

	if want != got {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestIsChincodeID(t *testing.T) {
	containerID := "dev.peer0.org1.mycc.v1.0"
	want := false
	got := IsContainerID(containerID)

	if want != got {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetDeployment(t *testing.T) {
	ID := "dev-peer0-org1-example-com-mycc_1-cb8910118e5b08ff9f8c59065afb8658b4be0fe0987a43d0556a30d0066630a4"
	deploymentName := GetDeploymentName(ID)

	if len(deploymentName) != 53 || !strings.HasPrefix(deploymentName, "chaincode") {
		t.Errorf("deploymentName is invalid: %s", deploymentName)
	}
	fmt.Println(deploymentName)
}
