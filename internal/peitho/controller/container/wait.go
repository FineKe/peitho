// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (cc *ContainerController) Wait(c *gin.Context) {
	log.L(c).Info("wait container function called.")

	id := c.Param("id")

	err := cc.srv.Containers().Wait(context.Background(), id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, nil)
}
