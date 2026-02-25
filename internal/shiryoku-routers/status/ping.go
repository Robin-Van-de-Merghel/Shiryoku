package status

import (
	"net/http"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config/common"
	"github.com/gin-gonic/gin"
)

// Ping takes a list of checkers and returns a gin.HandlerFunc
func Ping(checkers ...common.Checker) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		for _, check := range checkers {
			ok, err := check(ctx)
			if err != nil || !ok {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"message": "not ready",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
