package entities

type HostQueryParams struct {
	ID        int64  `json:"id"`
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
}
