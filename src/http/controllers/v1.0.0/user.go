package v1_0_0

import (
	"log"
	"main/src/http/controllers"
	"main/src/http/controllers/v1.0.0/requests"
	"main/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*controllers.Controller

	userService *services.UserService
}

func (controller *UserController) GetUser(c *gin.Context) {

	request := requests.UserGetByIdRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		controller.SendError(c, err, 422)
		return
	}
	log.Println("User get", request.Id)
	user, err := controller.userService.GetUserById(c, request.Id)
	if err != nil {
		c.Status(503)
		return
	}
	c.JSON(http.StatusOK, user)
}

func NewUserController(
	userService *services.UserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}
