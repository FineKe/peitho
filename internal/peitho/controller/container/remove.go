// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (cc *ContainerController) Remove(c *gin.Context) {
	log.L(c).Info("delete container function called.")

	id := c.Param("id")

	_ = cc.srv.Containers().Remove(context.Background(), id)

	c.JSON(204, nil)
}
