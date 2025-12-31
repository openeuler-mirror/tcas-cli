package teetool

type TrustDeviceReport struct {
	DeviceType   string      `json:"type"`
	DeviceReport string      `json:"device_report"`
	Parameter    interface{} `json:"parameter"`
}

type NodeEvidence struct {
	Tee         string               `json:"tee"`
	TeeReport   string               `json:"tee_report"`
	Parameter   interface{}          `json:"parameter"`
	TrustDevice []*TrustDeviceReport `json:"trust_devices"`
	ReportData  []byte               `json:"-"`
	UserData    string               `json:"user_data"`
	InitData    string               `json:"init_data"`
	EventLog    string               `json:"event_log"`
}
