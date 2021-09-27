// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (cc *ContainerController) Upload(c *gin.Context) {
	log.L(c).Info("Upload archive function called.")

	path := c.Query("path")
	id := c.Param("id")

	log.Debugf("upload archive, id: %s, path: %s", id, path)

	if err := cc.srv.Containers().Upload(context.Background(), id, path, c.Request.Body); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, nil)
}

func (cc *ContainerController) Fetch(c *gin.Context) {
	log.L(c).Info("Fetch archive function called.")

	id := c.Param("id")
	path := c.Query("path")

	// fetch archive
	reader, err := cc.srv.Containers().Fetch(context.Background(), id, path)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.Writer.WriteHeader(200)

	defer func() {
		reader.Close()
		c.Writer.Flush()
	}()

	// copy to response
	wlen, err := io.Copy(c.Writer, reader)
	if err != nil {
		log.Errorf("copy content to response failed: %v", err)

		return
	}

	log.Debugf("fetch archive size: %d byte", wlen)
}
