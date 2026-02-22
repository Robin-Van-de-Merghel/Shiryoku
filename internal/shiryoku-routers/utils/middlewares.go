package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		// Log the panic
		fmt.Printf("Panic occurred: %v\n", err)

		// Get timestamp to help debugging
		timestamp := time.Now().UTC()
		
		// Return custom 500 response
		c.JSON(500, gin.H{
			"code":    500,
			"timestamp": timestamp, 
			"message": "Internal Server Error. Please contact an administrator",
			"detail":   fmt.Sprintf("%v", err),
		})
		c.Abort()
	})
}
