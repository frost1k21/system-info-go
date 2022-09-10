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
)

func getSysInfoCollection(wsName string) (map[reflect.Type]sysInfoCollect, error) {
	sysInfoCollectSlice := map[reflect.Type]sysInfoCollect{
		reflect.TypeOf(OperatingSystem{}):            {value: &[]OperatingSystem{}, wmiClass: "Win32_OperatingSystem", namespace: "root\\cimv2"},
		reflect.TypeOf(Cpu{}):                        {value: &[]Cpu{}, wmiClass: "Win32_Processor", namespace: "root\\cimv2"},
		reflect.TypeOf(Motherboard{}):                {value: &[]Motherboard{}, wmiClass: "win32_BaseBoard", namespace: "root\\cimv2"},
		reflect.TypeOf(SystemUser{}):                 {value: &[]SystemUser{}, wmiClass: "Win32_ComputerSystem", namespace: "root\\cimv2"},
		reflect.TypeOf(physicalMemory{}):             {value: &[]physicalMemory{}, wmiClass: "Win32_PhysicalMemory", namespace: "root\\cimv2"},
		reflect.TypeOf(msSmBios_rawSMBiosTables{}):   {value: &[]msSmBios_rawSMBiosTables{}, wmiClass: "MSSmBios_RawSMBiosTables", namespace: "root\\wmi"},
		reflect.TypeOf(DiskDrive{}):                  {value: &[]DiskDrive{}, wmiClass: "Win32_DiskDrive", namespace: "root\\cimv2"},
		reflect.TypeOf(VideoAdapter{}):               {value: &[]VideoAdapter{}, wmiClass: "Win32_VideoController", namespace: "root\\cimv2"},
		reflect.TypeOf(wmiMonitorID{}):               {value: &[]wmiMonitorID{}, wmiClass: "WmiMonitorID", namespace: "root\\wmi"},
		reflect.TypeOf(wmiMonitorConnectionParams{}): {value: &[]wmiMonitorConnectionParams{}, wmiClass: "WmiMonitorConnectionParams", namespace: "root\\wmi"},
	}
	for _, sysInfo := range sysInfoCollectSlice {
		q := wmi.CreateQuery(sysInfo.value, "", sysInfo.wmiClass)
		err := wmi.Query(q, sysInfo.value, wsName, sysInfo.namespace)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	return sysInfoCollectSlice, nil
}

func mapSystemInfoCollectToWsInfo(wsName string, sysInfoMap map[reflect.Type]sysInfoCollect) *WsInfo {
	result := WsInfo{}
	result.WsName = wsName
	result.Cpu = mapToCpu(sysInfoMap[reflect.TypeOf(Cpu{})])
	result.Motherboard = mapToMotherBoard(sysInfoMap[reflect.TypeOf(Motherboard{})])
	result.SystemUser = mapToSystemUser(sysInfoMap[reflect.TypeOf(SystemUser{})])
	result.OperatingSystem = mapToOperatingSystem(sysInfoMap[reflect.TypeOf(OperatingSystem{})])
	result.Rams = mapToRams(sysInfoMap[reflect.TypeOf(physicalMemory{})], sysInfoMap[reflect.TypeOf(msSmBios_rawSMBiosTables{})])
	result.DiskDrives = mapToDiskDrives(sysInfoMap[reflect.TypeOf(DiskDrive{})])
	result.VideoAdapters = mapToVideoAdapters(sysInfoMap[reflect.TypeOf(VideoAdapter{})])
	result.Monitors = mapToMonitors(sysInfoMap[reflect.TypeOf(wmiMonitorID{})], sysInfoMap[reflect.TypeOf(wmiMonitorConnectionParams{})])
	return &result
}

func mapToMonitors(wmiMonitorIDsCollection sysInfoCollect, wmiMonitorConnectionParamsCollection sysInfoCollect) []Monitor {
	var monitors []Monitor
	for _, monitorIDInformation := range *wmiMonitorIDsCollection.value.(*[]wmiMonitorID) {
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

	for _, monitorConnectionParams := range *wmiMonitorConnectionParamsCollection.value.(*[]wmiMonitorConnectionParams) {
		instanceName := monitorConnectionParams.InstanceName
		videoOutputTechnology := monitorConnectionParams.VideoOutputTechnology
		for i, monitor := range monitors {
			if monitor.InstanceName == instanceName {
				monitors[i].MonitorConnectionPort = videoOutputTechnology
				monitors[i].MonitorConnectionPortName = monitorPorts[monitorConnectionPorts(videoOutputTechnology)]
			}
		}
	}
	return monitors
}

func mapToVideoAdapters(collect sysInfoCollect) []VideoAdapter {
	var dst []VideoAdapter
	for _, adapter := range *collect.value.(*[]VideoAdapter) {
		if adapter.AdapterRAM > 0 {
			dst = append(dst, adapter)
		}
	}
	return dst
}

func mapToRams(physMemCollection sysInfoCollect, msSmBiosCollection sysInfoCollect) []Ram {
	var rams []Ram
	for _, memory := range *physMemCollection.value.(*[]physicalMemory) {
		currentMemory := Ram{
			Capacity: memory.Capacity,
			Speed:    memory.Speed,
		}
		rams = append(rams, currentMemory)
	}

	var memoryTypeStrings []string
	for _, bio := range *msSmBiosCollection.value.(*[]msSmBios_rawSMBiosTables) {
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
	return rams
}

func mapToDiskDrives(collect sysInfoCollect) []DiskDrive {
	return *collect.value.(*[]DiskDrive)
}

func mapToSystemUser(collect sysInfoCollect) SystemUser {
	return (*collect.value.(*[]SystemUser))[0]
}

func mapToMotherBoard(collect sysInfoCollect) Motherboard {
	return (*collect.value.(*[]Motherboard))[0]
}

func mapToCpu(collect sysInfoCollect) Cpu {
	return (*collect.value.(*[]Cpu))[0]
}

func mapToOperatingSystem(collect sysInfoCollect) OperatingSystem {
	value := (*collect.value.(*[]OperatingSystem))[0].Name
	value = strings.TrimSpace(strings.Split(value, "|")[0])
	(*collect.value.(*[]OperatingSystem))[0].Name = value
	return (*collect.value.(*[]OperatingSystem))[0]
}

func asyncInfoCollectFunction(wsName string, ch chan *WsInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	collection, err := getSysInfoCollection(wsName)
	if err != nil {
		log.Println(err.Error())
	}
	result := mapSystemInfoCollectToWsInfo(wsName, collection)

	ch <- result
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
