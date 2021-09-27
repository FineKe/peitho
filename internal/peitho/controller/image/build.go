// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package image

import (
	"bufio"
	"context"

	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/pkg/log"
)

func (ic *ImageController) Build(c *gin.Context) {
	log.L(c).Info("build image function called.")

	dockerfile := c.Query("dockerfile")
	tags := c.QueryArray("t")

	value, err := ic.srv.Images().Build(context.Background(), dockerfile, tags, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	// c.Writer.Header().Add("Transfer-Encoding","chunked")
	c.Writer.WriteHeader(200)
	reader := bufio.NewReader(value)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		_, _ = c.Writer.Write(line)
		c.Writer.Flush()
	}
}
