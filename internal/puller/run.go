// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package puller

import (
	"github.com/tianrandailove/peitho/internal/puller/config"
	"github.com/tianrandailove/peitho/internal/puller/service"
	"github.com/tianrandailove/peitho/pkg/docker"
	"github.com/tianrandailove/peitho/pkg/log"
	"github.com/tianrandailove/peitho/pkg/options"
)

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	// new docker client
	dockerService, err := docker.NewDockerService(&options.DockerOption{
		Endpoint: cfg.PullerOption.DockerEndpoint,
	}, options.NewPeithoOption())
	if err != nil {
		panic(err)
	}

	// new service
	srv := service.NewService(dockerService)
	err = srv.Pullers().PullImage(cfg.PullerOption.Image, cfg.PullerOption.PullAddress)
	if err != nil {
		log.Errorf("pull and load image failed: %v", err)

		return err
	}

	return nil
}
