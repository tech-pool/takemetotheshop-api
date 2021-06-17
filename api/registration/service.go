package registration

import (
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"message": "Hello Registration",
	})

}
