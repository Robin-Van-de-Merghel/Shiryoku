package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logicnmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

// SearchNmapScans returns a gin handler injected with the given NmapDBIface.
// Call with no argument in tests to use the DummyNmapDB automatically.
func SearchNmapScans(nmapDB osdb.NmapDBIface) func(c *gin.Context) {
	return func(c *gin.Context) {
		var params models.SearchParams

		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		if !utils.ValidateAndRespond(c, params, utils.SearchSchema) {
			return
		}

		params.SetDefaults()

		nmapResult, err := logicnmap.SearchNmapScans(c.Request.Context(), &params, nmapDB)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, nmapResult)
	}
}
