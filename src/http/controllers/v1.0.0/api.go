package v1_0_0

import (
	"main/src/http/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	*controllers.Controller
}

func (controller *ApiController) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "1.0.0"})
}

func NewApiController() *ApiController {
	return &ApiController{}
}
