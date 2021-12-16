// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

//go:generate mockgen -self_package=github.com/tianrandailove/peitho/internal/peitho/service -destination mock_service.go -package service github.com/tianrandailove/peitho/internal/peitho/service Service,ImageSrv,ContainerSrv
import (
	"github.com/tianrandailove/peitho/pkg/docker"
)

var Srv Service

// Service defines functions used to return resource interface.
type Service interface {
	Pullers() PullerSrv
}

type service struct {
	docker docker.DockerService
}

func (s *service) Pullers() PullerSrv {
	return newPuller(s)
}

// NewService returns Service interface.
func NewService(docker docker.DockerService) Service {
	return &service{
		docker: docker,
	}
}
