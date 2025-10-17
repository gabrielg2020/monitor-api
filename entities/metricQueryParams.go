package entities

type MetricQueryParams struct {
	HostID    *int64 `form:"host_id"`
	StartTime *int64 `form:"start_time"`
	EndTime   *int64 `form:"end_time"`
	Limit     int    `form:"limit"`
	Order     string `form:"order"` // "asc" or "desc"
}
