package shiryoku_routers

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/status"
	"github.com/gin-gonic/gin"
)

func GetFilledRouter() *gin.Engine {
	// Main router
	router := gin.Default()

	// For docker-compose status
	router.GET("/ping", status.Ping)

	return router
}
