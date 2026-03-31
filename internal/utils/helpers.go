package utils

import (
	"strconv"

	userModels "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/models"
	"github.com/gin-gonic/gin"
)

func GetTenantID(c *gin.Context) *uint {
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			return u.TenantID
		}
	}
	return nil
}

func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

func GetUserName(c *gin.Context) string {
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*userModels.User); ok {
			return u.FirstName + " " + u.LastName
		}
	}
	return ""
}

func GetPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize
}
