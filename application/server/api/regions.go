package api

import (
	"grets_server/db"
	"grets_server/db/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProvinces 获取所有省份数据
func GetProvinces(c *gin.Context) {
	var provinces []models.Region

	result := db.GlobalMysql.Order("province_code").Find(&provinces)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取省份数据失败",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取省份数据成功",
		"data":    provinces,
	})
}
