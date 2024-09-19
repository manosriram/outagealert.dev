package integration

type Notify interface {
	SendAlert(monitorId, monitorName string) error
}
