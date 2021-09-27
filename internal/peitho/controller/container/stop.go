// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (cc *ContainerController) Stop(c *gin.Context) {
	log.L(c).Info("stop container function called.")

	id := c.Param("id")

	timeout := c.GetDuration("t")

	_ = cc.srv.Containers().Stop(context.Background(), id, timeout)

	c.JSON(204, nil)
}
