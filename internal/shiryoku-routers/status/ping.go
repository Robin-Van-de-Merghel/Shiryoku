package status

import "github.com/gin-gonic/gin"

// Check health status
func Ping(c *gin.Context) {

	// FIXME: Later on, this will be used by docker-compose for health status.
	//
	// It will be required to verify that all connections are done before saying "pong"

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
