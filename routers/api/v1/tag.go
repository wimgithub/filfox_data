package v1

import (
	"filfox_data/pkg/app"
	"filfox_data/pkg/e"
	"filfox_data/service/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	begin := c.Query("begin")
	end := c.Query("end")

	b1, _ := strconv.ParseInt(begin, 0, 64)
	e2, _ := strconv.ParseInt(end, 0, 64)
	path, err := excel.GetExcel(c, b1, e2, 0, "", "", "")
	if err != nil {
		panic(err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"file_name": path,
	})
}

// 下载文件
func Download(c *gin.Context) {
	appG := app.Gin{C: c}
	fileName := c.Param("file")
	fmt.Println(fileName)
	if fileName == "" {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	f, err := ioutil.ReadFile("./excel/" + fileName)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(f)))
	c.Writer.Write(f)
}
