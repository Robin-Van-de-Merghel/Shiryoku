package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logic_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
	logic_nmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

// SearchNmapScans returns a gin handler injected with the given opensearchDB.
func SearchNmapScans(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
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

		nmapResult, err := logic_common.Search[models.NmapHostDocument](
			c.Request.Context(), 
			nmapDB, &params, 
			logic_nmap.NMAP_HOSTS_INDEX,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, nmapResult)
	}
}
