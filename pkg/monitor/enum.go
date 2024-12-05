package monitor

type MonitorStatus string

const (
	StatusUp     MonitorStatus = "up"
	StatusDown   MonitorStatus = "down"
	StatusPaused MonitorStatus = "paused"
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

var PLAN_VS_MONITOR_COUNT = map[string]int64{
	"free":     FREE_MONITOR_LIMIT,
	"hobbyist": HOBBYIST_MONITOR_LIMIT,
	"pro":      PRO_MONITOR_LIMIT,
}
