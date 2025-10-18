package entities

type Host struct {
	ID        int64  `json:"id" db:"id"`
	Hostname  string `json:"hostname" db:"hostname"`
	IPAddress string `json:"ip_address" db:"ip_address"`
	Role      string `json:"role" db:"role"`
}
