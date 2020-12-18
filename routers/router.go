package routers

import (
	"filfox_data/middleware/core"
	"github.com/gin-gonic/gin"

	_ "filfox_data/docs"
	"filfox_data/routers/api/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(core.Cors())
	r.GET("/download/:file", v1.Download) //下载文件
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/get_fil", v1.GetTags)
	}
	return r
}
