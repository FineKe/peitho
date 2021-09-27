// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package image

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (ic *ImageController) Inspect(c *gin.Context) {
	log.L(c).Info("inspect image function called.")

	imageID := c.Param("name")

	value, err := ic.srv.Images().Inspect(context.Background(), imageID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, value)
}
