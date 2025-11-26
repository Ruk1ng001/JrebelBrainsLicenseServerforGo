package model

// JRebelLeaseResponse JRebel租约响应
type JRebelLeaseResponse struct {
	ServerVersion         string   `json:"serverVersion"`
	ServerProtocolVersion string   `json:"serverProtocolVersion"`
	ServerGUID            string   `json:"serverGuid"`
	GroupType             string   `json:"groupType"`
	ID                    int      `json:"id,omitempty"`
	LicenseType           int      `json:"licenseType,omitempty"`
	EvaluationLicense     bool     `json:"evaluationLicense,omitempty"`
	Signature             string   `json:"signature,omitempty"`
	ServerRandomness      string   `json:"serverRandomness,omitempty"`
	SeatPoolType          string   `json:"seatPoolType,omitempty"`
	StatusCode            string   `json:"statusCode"`
	Message               *string  `json:"msg,omitempty"`
	StatusMessage         *string  `json:"statusMessage,omitempty"`
	Offline               bool     `json:"offline,omitempty"`
	ValidFrom             *int64   `json:"validFrom,omitempty"`
	ValidUntil            *int64   `json:"validUntil,omitempty"`
	Company               string   `json:"company,omitempty"`
	OrderID               string   `json:"orderId,omitempty"`
	ZeroIds               []string `json:"zeroIds,omitempty"`
	LicenseValidFrom      int64    `json:"licenseValidFrom,omitempty"`
	LicenseValidUntil     int64    `json:"licenseValidUntil,omitempty"`
	CanGetLease           bool     `json:"canGetLease,omitempty"`
}

// JRebelValidateResponse JRebel验证响应
type JRebelValidateResponse struct {
	ServerVersion         string `json:"serverVersion"`
	ServerProtocolVersion string `json:"serverProtocolVersion"`
	ServerGUID            string `json:"serverGuid"`
	GroupType             string `json:"groupType"`
	StatusCode            string `json:"statusCode"`
	Company               string `json:"company"`
	CanGetLease           bool   `json:"canGetLease"`
	LicenseType           int    `json:"licenseType"`
	EvaluationLicense     bool   `json:"evaluationLicense"`
	SeatPoolType          string `json:"seatPoolType"`
}
