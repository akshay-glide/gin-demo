package userhdlr

import (
	"gin-demo/services/usersvc"
	"log"

	"github.com/gin-gonic/gin"
)

type UserHdlr struct {
	usersvc usersvc.UserService
	logger  *log.Logger
}

func (o *UserHdlr) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/create", o.CreateUser)
}

func NewUserHdlr(usersvcI usersvc.UserService, logger *log.Logger) *UserHdlr {
	return &UserHdlr{
		usersvc: usersvcI,
		logger:  logger,
	}
}
