// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"

	"github.com/tianrandailove/peitho/internal/peitho/service"
	"github.com/tianrandailove/peitho/pkg/log"
)

func (cc *ContainerController) Create(c *gin.Context) {
	log.L(c).Info("create container function called.")

	container := service.Container{}

	if err := c.ShouldBindJSON(&container); err != nil {
		c.JSON(400, err.Error())

		return
	}

	ID := c.Query("name")

	result, err := cc.srv.Containers().Create(context.Background(), ID, container)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchImage) {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		}

		return
	}

	c.JSON(200, result)
}
