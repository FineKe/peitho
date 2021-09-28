// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package image

import (
	"context"
	"errors"
	"github.com/tianrandailove/peitho/internal/peitho/service"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (ic *ImageController) Inspect(c *gin.Context) {
	log.L(c).Info("inspect image function called.")

	imageID := c.Param("name")

	value, err := ic.srv.Images().Inspect(context.Background(), imageID)
	if err != nil {
		if errors.As(err, service.ErrNoSuchImage) {
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

	c.JSON(200, value)
}
