package monitor

type MonitorStatus string

const (
	StatusUp   MonitorStatus = "up"
	StatusDown MonitorStatus = "down"
)

const (
	FREE_PLAN     = "free"
	HOBBYIST_PLAN = "hobbyist"
	PRO_PLAN      = "pro"
)

const (
	FREE_MONITOR_LIMIT     = 20
	HOBBYIST_MONITOR_LIMIT = 50
	PRO_MONITOR_LIMIT      = 150
)
