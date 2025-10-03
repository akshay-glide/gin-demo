package healthchecksvc

type HealthStatus struct {
	IsOk   bool
	Status any
	Msg    string
}

type HealthCheckSvc interface {
	GetHealth() *HealthStatus
}

type HealthCheckSvcImpl struct {
}

func (o *HealthCheckSvcImpl) GetHealth() *HealthStatus {
	return &HealthStatus{
		IsOk:   true,
		Status: "All functions are normal",
		Msg:    "All functions are normal",
	}
}

func NewHealthCheckSvc() HealthCheckSvc {
	return &HealthCheckSvcImpl{}
}
