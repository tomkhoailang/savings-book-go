package presenter

import "time"

type DashboardDayRevenueStats struct {
	RegulationType string
	Income         float64
	Outcome        float64
	DifferentCount float64
}
type DashboardDayCountStats struct {
	CurrentDate    time.Time
	OpenCount      int
	ClosedCount    int
	DifferentCount int
}
type DashboardDayCountStatsQuery struct {
	Time         time.Time `form:"time"`
	RegulationId string   `form:"regulationId"`
	Term         int      `form:"term"`
	InterestRate float64  `form:"interestRate"`
}