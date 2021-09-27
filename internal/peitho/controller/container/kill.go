// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (cc *ContainerController) Kill(c *gin.Context) {
	log.L(c).Info("kill container function called.")

	id := c.Param("id")
	signal := c.Query("signal")

	err := cc.srv.Containers().Kill(context.Background(), id, signal)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(204, nil)
}
