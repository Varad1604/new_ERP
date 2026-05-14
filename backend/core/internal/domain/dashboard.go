package domain

type DashboardMetrics struct {
	TotalRevenue     float64 `json:"total_revenue"`
	ActiveEmployees  int     `json:"active_employees"`
	PendingApprovals int     `json:"pending_approvals"`
}
