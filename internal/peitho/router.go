// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peitho

import (
	"github.com/gin-gonic/gin"

	"github.com/tianrandailove/peitho/internal/peitho/controller/container"
	"github.com/tianrandailove/peitho/internal/peitho/controller/image"
	"github.com/tianrandailove/peitho/internal/peitho/service"
)

func initRouter(g *gin.Engine) {
	installController(g)
}

func installController(g *gin.Engine) {
	containerController := container.NewContainerController(service.Srv)

	g.HEAD("/_ping", containerController.Ping)
	g.GET("/_ping", containerController.Ping)

	g.POST("/containers/create", containerController.Create)

	g.PUT("/containers/:id/archive", containerController.Upload)
	g.GET("/containers/:id/archive", containerController.Fetch)

	g.POST("/containers/:id/attach", containerController.Attach)
	g.POST("/containers/:id/start", containerController.Start)
	g.POST("/containers/:id/stop", containerController.Stop)
	g.POST("/containers/:id/kill", containerController.Kill)
	g.DELETE("/containers/:id", containerController.Remove)
	g.POST("/containers/:id/wait", containerController.Wait)

	imageController := image.NewImageController(service.Srv)

	g.POST("/images/create", imageController.Create)
	g.GET("/images/:name/*json", imageController.Inspect)
	g.POST("/build", imageController.Build)
	g.GET("/tar/:name", imageController.Download)
}
