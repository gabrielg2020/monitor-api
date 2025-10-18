package entities

type Host struct {
	ID        int64  `json:"id" db:"id"`
	Hostname  string `json:"hostname" db:"hostname"`
	IPAddress string `json:"ip_address" db:"ip_address"`
	Role      string `json:"role" db:"role"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	LastSeen  int64  `json:"last_seen" db:"last_seen"`
}
