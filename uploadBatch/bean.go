package uploadBatch

type LogProperties struct {
	Level         string `json:"level"`
	Logger        string `json:"logger"`
	Language      string `json:"language"`
	Environment   string `json:"env"`
	Service       string `json:"service"`
	Hostname      string `json:"hostname"`
	DataCenter    string `json:"dc"`
	CloudMinion   string `json:"cloudMinion"`
	Message       string `json:"message"`
	ThrownMessage string `json:"thrownMessage"`
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
	Tags          []string      `json:"tags,omitempty"`
	CrashIDSource string        `json:"crashIdSourceField"`
}

type LogMessage struct {
	UploadBean UploadBean `json:"uploadBean"`
	StackTrace string     `json:"stackTrace"`
	Type       string     `json:"type"`
	Timestamp  string     `json:"timestamp"` // ISO 8601 timestamp
}
