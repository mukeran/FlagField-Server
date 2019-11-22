package resource

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func loadResource(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	resourceUUID := c.Param("resourceUUID")
	resource, err := cpt.GetResourceByUUID(tx, resourceUUID)
	if err != nil {
		return false
	}
	c.Set("resource", resource)
	return true
}

func LoadResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadResource(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
