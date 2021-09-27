// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package image

import "github.com/tianrandailove/peitho/internal/peitho/service"

type ImageController struct {
	srv service.Service
}

func NewImageController(srv service.Service) *ImageController {
	return &ImageController{
		srv: srv,
	}
}
