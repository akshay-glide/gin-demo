package userhdlr

import (
	"gin-demo/handlers"
	"gin-demo/services/usersvc"

	"github.com/gin-gonic/gin"
)

func (o *UserHdlr) CreateUser(c *gin.Context) {
	var user usersvc.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handlers.APIResponseBadRequest(c, "INVALID_REQUEST", err, "Invalid request payload")
		return
	}

	err := o.usersvc.Create(&user)
	if err != nil {
		handlers.APIResponseInternalServerError(c, "USER_CREATION_FAILED", err, "Failed to create user")
		return
	}

	handlers.APIResponseOK(c, user, "User created successfully")
}
