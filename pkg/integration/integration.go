package integration

type Notification interface {
	SendAlert(monitorId, monitorName string) error
	Notify() error
}
