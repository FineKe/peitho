// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peitho

import (
	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/internal/peitho/config"
	"github.com/tianrandailove/peitho/internal/peitho/service"
	"github.com/tianrandailove/peitho/pkg/docker"
	"github.com/tianrandailove/peitho/pkg/k8s"
)

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	// new k8s client
	k8sService, err := k8s.NewK8sService(cfg.K8sOption)
	if err != nil {
		panic(err)
	}

	// new docker client
	dockerService, err := docker.NewDockerService(cfg.DockerOption)
	if err != nil {
		panic(err)
	}

	// new service
	service.Srv = service.NewService(dockerService, k8sService)

	engine := gin.Default()
	initRouter(engine)

	return engine.Run()
}
