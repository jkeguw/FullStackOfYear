package claims

type Claims struct {
	UserID   string `json:"uid"`
	Role     string `json:"role"`
	DeviceID string `json:"deviceId"`
	Type     string `json:"type"`
}
