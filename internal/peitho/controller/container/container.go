// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import "github.com/tianrandailove/peitho/internal/peitho/service"

type ContainerController struct {
	srv service.Service
}

func NewContainerController(srv service.Service) *ContainerController {
	return &ContainerController{
		srv: srv,
	}
}
