package tracer

type LogProperties struct {
	Level    string `json:"level"`
	Logger   string `json:"logger"`
	Language string `json:"language"`
	Env      string `json:"env"`
	Service  string `json:"service"`
	Hostname string `json:"hostname"`
	DC       string `json:"dc"`
	//CloudMinion   string `json:"cloudMinion"`
	Message       string `json:"message"`
	ThrownMessage string `json:"thrownMessage,omitempty"`
	RequestID     string `json:"requestId,omitempty"`
}

type UploadBean struct {
	OSVersion     string        `json:"osVersion"`
	Vendor        string        `json:"vendor"`
	VersionName   string        `json:"versionName"`
	VersionCode   int           `json:"versionCode"`
	DeviceID      string        `json:"deviceId"`
	Module        string        `json:"module"`
	Properties    LogProperties `json:"properties"`
	Tags          []string      `json:"tags"`
	CrashIDSource string        `json:"crashIdSourceField"`
}

type LogMessage struct {
	UploadBean UploadBean `json:"uploadBean"`
	StackTrace string     `json:"stackTrace"`
	Type       string     `json:"type"`
	Timestamp  string     `json:"timestamp"` // ISO 8601 timestamp
}
