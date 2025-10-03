package healthcheckhdlr

import (
	"gin-demo/handlers"
	"gin-demo/services/healthchecksvc"

	"github.com/gin-gonic/gin"
)

type HealthCheckHdlr struct {
	healthCheckSvc healthchecksvc.HealthCheckSvc
}

func (o *HealthCheckHdlr) HealthCheck(c *gin.Context) {
	healthStatus := o.healthCheckSvc.GetHealth()

	if healthStatus.IsOk {
		handlers.APIResponseOK(c, healthStatus.Status, healthStatus.Msg)
	} else {
		handlers.APIResponseInternalServerError(c, "HEALTH_BAD", healthStatus.Status, healthStatus.Msg)
	}
}

func (o *HealthCheckHdlr) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", o.HealthCheck)
}

func NewHealthCheckHdlr(healthCheckSvc healthchecksvc.HealthCheckSvc) *HealthCheckHdlr {
	return &HealthCheckHdlr{
		healthCheckSvc: healthCheckSvc,
	}
}
