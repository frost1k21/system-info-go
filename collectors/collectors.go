package collectors

import (
	. "example/generics/core"
	"github.com/StackExchange/wmi"
	"log"
	"reflect"
	"strings"
	"sync"
)

type wmiMonitorID struct {
	UserFriendlyName []int32
	ManufacturerName []int32
	InstanceName     string
}

type wmiMonitorConnectionParams struct {
	VideoOutputTechnology int64
	InstanceName          string
}

type physicalMemory struct {
	Capacity int64
	Speed    int64
}

type msSmBios_rawSMBiosTables struct {
	SMBiosData []byte
}

type computerInfoConstraint interface {
	OperatingSystem | Cpu | Motherboard | SystemUser | physicalMemory | msSmBios_rawSMBiosTables | DiskDrive | VideoAdapter | wmiMonitorID | wmiMonitorConnectionParams
}

type sysInfoCollect struct {
	value     interface{}
	wmiClass  string
	namespace string
}

type monitorConnectionPorts int64

const (
	D3DKMDT_VOT_UNINITIALIZED        monitorConnectionPorts = -2
	D3DKMDT_VOT_OTHER                                       = -1
	D3DKMDT_VOT_HD15                                        = 0
	D3DKMDT_VOT_SVIDEO                                      = 1
	D3DKMDT_VOT_COMPOSITE_VIDEO                             = 2
	D3DKMDT_VOT_COMPONENT_VIDEO                             = 3
	D3DKMDT_VOT_DVI                                         = 4
	D3DKMDT_VOT_HDMI                                        = 5
	D3DKMDT_VOT_LVDS                                        = 6
	D3DKMDT_VOT_D_JPN                                       = 8
	D3DKMDT_VOT_SDI                                         = 9
	D3DKMDT_VOT_DISPLAYPORT_EXTERNAL                        = 10
	D3DKMDT_VOT_DISPLAYPORT_EMBEDDED                        = 11
	D3DKMDT_VOT_UDI_EXTERNAL                                = 12
	D3DKMDT_VOT_UDI_EMBEDDED                                = 13
	D3DKMDT_VOT_SDTVDONGLE                                  = 14
	D3DKMDT_VOT_MIRACAST                                    = 15
)

var (
	mapAbbreviationsToMonitorManufacturer = map[string]string{
		"AAC": "AcerView",
		"ACR": "Acer",
		"AOC": "AOC",
		"AIC": "AG Neovo",
		"APP": "Apple Computer",
		"AST": "AST Research",
		"AUO": "Asus",
		"BNQ": "BenQ",
		"CMO": "Acer",
		"CPL": "Compal",
		"CPQ": "Compaq",
		"CPT": "Chunghwa Pciture Tubes, Ltd.",
		"CTX": "CTX",
		"DEC": "DEC",
		"DEL": "Dell",
		"DPC": "Delta",
		"DWE": "Daewoo",
		"EIZ": "EIZO",
		"ELS": "ELSA",
		"ENC": "EIZO",
		"EPI": "Envision",
		"FCM": "Funai",
		"FUJ": "Fujitsu",
		"FUS": "Fujitsu-Siemens",
		"GSM": "LG Electronics",
		"GWY": "Gateway 2000",
		"HEI": "Hyundai",
		"HIT": "Hyundai",
		"HSL": "Hansol",
		"HTC": "Hitachi/Nissei",
		"HWP": "HP",
		"IBM": "IBM",
		"ICL": "Fujitsu ICL",
		"IVM": "Iiyama",
		"KDS": "Korea Data Systems",
		"LEN": "Lenovo",
		"LGD": "Asus",
		"LPL": "Fujitsu",
		"MAX": "Belinea",
		"MEI": "Panasonic",
		"MEL": "Mitsubishi Electronics",
		"MS_": "Panasonic",
		"NAN": "Nanao",
		"NEC": "NEC",
		"NOK": "Nokia Data",
		"NVD": "Fujitsu",
		"OPT": "Optoma",
		"PHL": "Philips",
		"REL": "Relisys",
		"SAN": "Samsung",
		"SAM": "Samsung",
		"SBI": "Smarttech",
		"SGI": "SGI",
		"SNY": "Sony",
		"SRC": "Shamrock",
		"SUN": "Sun Microsystems",
		"SEC": "Hewlett-Packard",
		"TAT": "Tatung",
		"TOS": "Toshiba",
		"TSB": "Toshiba",
		"VSC": "ViewSonic",
		"ZCM": "Zenith",
		"UNK": "Unknown",
		"_YV": "Fujitsu",
	}

	monitorPorts = map[monitorConnectionPorts]string{
		D3DKMDT_VOT_HD15: "VGA",
		D3DKMDT_VOT_HDMI: "HDMI",
		D3DKMDT_VOT_DVI:  "DVI",
	}

	sysInfoCollectSlice = []sysInfoCollect{
		{value: []OperatingSystem{}, wmiClass: "Win32_OperatingSystem", namespace: "root\\cimv2"},
		{value: []Cpu{}, wmiClass: "Win32_Processor", namespace: "root\\cimv2"},
		{value: []Motherboard{}, wmiClass: "win32_BaseBoard", namespace: "root\\cimv2"},
		{value: []SystemUser{}, wmiClass: "Win32_ComputerSystem", namespace: "root\\cimv2"},
		{value: []physicalMemory{}, wmiClass: "Win32_PhysicalMemory", namespace: "root\\cimv2"},
		{value: []msSmBios_rawSMBiosTables{}, wmiClass: "MSSmBios_RawSMBiosTables", namespace: "root\\wmi"},
		{value: []DiskDrive{}, wmiClass: "Win32_DiskDrive", namespace: "root\\cimv2"},
		{value: []VideoAdapter{}, wmiClass: "Win32_VideoController", namespace: "root\\cimv2"},
		{value: []wmiMonitorID{}, wmiClass: "WmiMonitorID", namespace: "root\\wmi"},
		{value: []wmiMonitorConnectionParams{}, wmiClass: "WmiMonitorConnectionParams", namespace: "root\\wmi"},
	}
)

func getOperatingSystemInformation(wsName string) (OperatingSystem, error) {
	dst, err := getGenericInfo[OperatingSystem](wsName, "Win32_OperatingSystem")
	if err != nil {
		return OperatingSystem{}, err
	}
	dst[0].Name = strings.TrimSpace(strings.Split(dst[0].Name, "|")[0])
	return dst[0], nil
}

func getCpuInformation(wsName string) (Cpu, error) {
	dst, err := getGenericInfo[Cpu](wsName, "Win32_Processor")
	if err != nil {
		return Cpu{}, err
	}
	return dst[0], nil
}

func getMotherboardInformation(wsName string) (Motherboard, error) {
	dst, err := getGenericInfo[Motherboard](wsName, "win32_BaseBoard")
	if err != nil {
		return Motherboard{}, err
	}
	return dst[0], nil
}

func getSystemUserInformation(wsName string) (SystemUser, error) {
	dst, err := getGenericInfo[SystemUser](wsName, "Win32_ComputerSystem")
	if err != nil {
		return SystemUser{}, err
	}
	return dst[0], nil
}

func getPhysicalMemoriesInformation(wsName string) ([]physicalMemory, error) {
	dst, err := getGenericInfo[physicalMemory](wsName, "Win32_PhysicalMemory")
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func getMsSmBios_rawSMBiosTables(wsName string) ([]msSmBios_rawSMBiosTables, error) {
	var dst []msSmBios_rawSMBiosTables
	q := wmi.CreateQuery(&dst, "", "MSSmBios_RawSMBiosTables")
	err := wmi.Query(q, &dst, wsName, "root\\wmi")

	if err != nil {
		return nil, err
	}
	return dst, nil
}

func getDiskDrivesInformation(wsName string) ([]DiskDrive, error) {
	dst, err := getGenericInfo[DiskDrive](wsName, "Win32_DiskDrive")
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func getVideoAdaptersInformation(wsName string) ([]VideoAdapter, error) {
	dst, err := getGenericInfo[VideoAdapter](wsName, "Win32_VideoController")
	if err != nil {
		return nil, err
	}
	var filteredDst []VideoAdapter
	for _, adapter := range dst {
		if adapter.AdapterRAM > 0 {
			filteredDst = append(filteredDst, adapter)
		}
	}
	return filteredDst, err
}

func getWmiMonitorIDsInformation(wsName string) ([]wmiMonitorID, error) {
	var dst []wmiMonitorID
	q := wmi.CreateQuery(&dst, "", "WmiMonitorID")
	err := wmi.Query(q, &dst, wsName, "root\\wmi")

	if err != nil {
		return nil, err
	}
	return dst, nil
}

func getWmiMonitorConnectionParams(wsName string) ([]wmiMonitorConnectionParams, error) {
	var dst []wmiMonitorConnectionParams
	q := wmi.CreateQuery(&dst, "", "WmiMonitorConnectionParams")
	err := wmi.Query(q, &dst, wsName, "root\\wmi")

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return dst, nil
}

func getGenericInfo[T computerInfoConstraint](wsName string, wmiClass string) ([]T, error) {
	var dst []T
	q := wmi.CreateQuery(&dst, "", wmiClass)
	err := wmi.Query(q, &dst, wsName)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func getSysInfoCollection() {
	for _, sysInfo := range sysInfoCollectSlice {
		log.Println(reflect.TypeOf(sysInfo.value))
	}
}

func asyncInfoCollectFunction(wsName string, ch chan *WsInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	getSysInfoCollection()

	operatingSystem, _ := getOperatingSystemInformation(wsName)
	cpu, _ := getCpuInformation(wsName)
	motherBoard, _ := getMotherboardInformation(wsName)
	user, _ := getSystemUserInformation(wsName)
	disks, _ := getDiskDrivesInformation(wsName)
	videoAdapters, _ := getVideoAdaptersInformation(wsName)

	physicalMemories, _ := getPhysicalMemoriesInformation(wsName)
	msRawBios, _ := getMsSmBios_rawSMBiosTables(wsName)

	var rams []Ram
	for _, memory := range physicalMemories {
		currentMemory := Ram{
			Capacity: memory.Capacity,
			Speed:    memory.Speed,
		}
		rams = append(rams, currentMemory)
	}

	var memoryTypeStrings []string
	for _, bio := range msRawBios {
		dataByte := bio.SMBiosData
		for i := 0; i < len(dataByte); i++ {
			if dataByte[i] == 17 && dataByte[i-1] == 00 && dataByte[i-2] == 00 {
				memType := dataByte[i+0x12]
				switch memType {
				case 26:
					memoryTypeStrings = append(memoryTypeStrings, "DDR4")
					break
				case 20:
					memoryTypeStrings = append(memoryTypeStrings, "DDR")
					break
				case 21:
					memoryTypeStrings = append(memoryTypeStrings, "DDR2")
					break
				case 24:
					memoryTypeStrings = append(memoryTypeStrings, "DDR3")
					break
				default:
					break
				}
			}
		}
	}
	if len(memoryTypeStrings) > 0 {
		for i, _ := range rams {
			rams[i].MemType = memoryTypeStrings[0]
		}
	}

	monitorIDsInformation, _ := getWmiMonitorIDsInformation(wsName)
	monitorWmiConnectionInformation, _ := getWmiMonitorConnectionParams(wsName)

	var monitors []Monitor
	for _, monitorIDInformation := range monitorIDsInformation {
		abbMan := strings.Trim(string(monitorIDInformation.ManufacturerName), "\x00")
		manName := mapAbbreviationsToMonitorManufacturer[abbMan]
		userFriendlyName := strings.Trim(string(monitorIDInformation.UserFriendlyName), "\x00")
		currentMonitor := Monitor{
			Model:        userFriendlyName,
			Manufacturer: manName,
			InstanceName: monitorIDInformation.InstanceName,
		}

		monitors = append(monitors, currentMonitor)
	}

	for _, monitorConnectionParams := range monitorWmiConnectionInformation {
		instanceName := monitorConnectionParams.InstanceName
		videoOutputTechnology := monitorConnectionParams.VideoOutputTechnology
		for i, monitor := range monitors {
			if monitor.InstanceName == instanceName {
				monitors[i].MonitorConnectionPort = videoOutputTechnology
				monitors[i].MonitorConnectionPortName = monitorPorts[monitorConnectionPorts(videoOutputTechnology)]
			}
		}
	}

	wsInfo := WsInfo{
		WsName:          wsName,
		OperatingSystem: operatingSystem,
		Cpu:             cpu,
		Motherboard:     motherBoard,
		SystemUser:      user,
		Rams:            rams,
		DiskDrives:      disks,
		VideoAdapters:   videoAdapters,
		Monitors:        monitors,
	}
	ch <- &wsInfo
}

func GetComputersInfo(wsNames []string) []WsInfo {
	var wsInfos []WsInfo
	var wg sync.WaitGroup
	responseChannel := make(chan *WsInfo, len(wsNames))

	wg.Add(len(wsNames))
	for _, name := range wsNames {
		go asyncInfoCollectFunction(name, responseChannel, &wg)
	}
	wg.Wait()
	close(responseChannel)
	for system := range responseChannel {
		wsInfos = append(wsInfos, *system)
	}
	return wsInfos
}
