package image

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tianrandailove/peitho/pkg/log"
	"os"
)

func (ic *ImageController) Download(c *gin.Context) {
	log.L(c).Info("Download image function called.")

	imageID := c.Param("name")
	if imageID == "" {
		c.JSON(400, gin.H{"message": "image not be empty"})

		return
	}
	fileName := fmt.Sprintf("%s.tar", imageID)
	_, err := os.Stat(fileName)
	if err != nil {
		log.Errorf("file not exists:%v", err)
		c.JSON(404, gin.H{"message": err.Error()})

		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(fileName)
	return
}
