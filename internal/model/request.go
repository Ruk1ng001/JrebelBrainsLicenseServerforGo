package model

// JRebelLeaseRequest JRebel租约请求
type JRebelLeaseRequest struct {
	Randomness  string `json:"randomness"`
	Username    string `json:"username"`
	GUID        string `json:"guid"`
	Offline     bool   `json:"offline"`
	ClientTime  int64  `json:"clientTime,omitempty"`
	OfflineDays int    `json:"offlineDays,omitempty"`
}
