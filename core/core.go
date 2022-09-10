package core

type Result struct {
	Success WsInfo `json:"success"`
	Error   string `json:"error"`
}

type WsInfo struct {
	WsName          string          `json:"wsName"`
	Cpu             Cpu             `json:"cpu"`
	Motherboard     Motherboard     `json:"motherboard"`
	SystemUser      SystemUser      `json:"systemUser"`
	OperatingSystem OperatingSystem `json:"operatingSystem"`
	Rams            []Ram           `json:"rams"`
	DiskDrives      []DiskDrive     `json:"diskDrives"`
	VideoAdapters   []VideoAdapter  `json:"videoAdapters"`
	Monitors        []Monitor       `json:"monitors"`
}

type Cpu struct {
	Name          string `json:"name"`
	MaxClockSpeed int64  `json:"frequency"`
}

type Motherboard struct {
	Manufacturer string `json:"manufacturer"`
	Product      string `json:"model"`
}

type SystemUser struct {
	UserName string `json:"login"`
}

type OperatingSystem struct {
	Name                    string `json:"name"`
	OSArchitecture          string `json:"osArchitecture"`
	ServicePackMajorVersion byte   `json:"servicePack"`
}

type Ram struct {
	Capacity int64  `json:"capacity"`
	Speed    int64  `json:"speed"`
	MemType  string `json:"memType"`
}

type DiskDrive struct {
	Model string `json:"name"`
	Size  int64  `json:"size"`
}

type VideoAdapter struct {
	Name       string `json:"name"`
	AdapterRAM int64  `json:"memory"`
}

type Monitor struct {
	Model                     string `json:"model"`
	Manufacturer              string `json:"manufacturer"`
	InstanceName              string `json:"instanceName"`
	MonitorConnectionPort     int64  `json:"monitorConnectionPort"`
	MonitorConnectionPortName string `json:"monitorConnectionPortName"`
}
