package entities

type HostQueryParams struct {
	ID        int64  `form:"id"`
	Hostname  string `form:"hostname"`
	IPAddress string `form:"ip_address"`
}
