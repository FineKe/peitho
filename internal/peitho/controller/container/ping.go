// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package container

import (
	"github.com/gin-gonic/gin"
)

func (cc *ContainerController) Ping(c *gin.Context) {
	//log.L(c).Info("ping function called.")

	c.JSON(200, gin.H{"message": "OK"})
}
