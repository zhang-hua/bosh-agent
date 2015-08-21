package monit

import (
	"encoding/xml"
	"strconv"
)

type status struct {
	XMLName     xml.Name `xml:"monit"`
	ID          string   `xml:"id,attr"`
	Incarnation string   `xml:"incarnation,attr"`
	Version     string   `xml:"version,attr"`

	Services      servicesTag
	ServiceGroups serviceGroupsTag
}

type servicesTag struct {
	XMLName  xml.Name     `xml:"services"`
	Services []serviceTag `xml:"service"`
}

type serviceTag struct {
	XMLName  xml.Name  `xml:"service"`
	Name     string    `xml:"name,attr"`
	CollectedSeconds uint64 `xml:"collected_sec"`
	Status   int       `xml:"status"`
	Monitor  int       `xml:"monitor"`
	Pid      int       `xml:"pid"`
	Ppid     int       `xml:"ppid"`
	Uptime   int       `xml:"uptime"`
	Children int       `xml:"children"`
	Memory   memoryTag `xml:"memory"`
	CPU      cpuTag    `xml:"cpu"`
	Port     portTag    `xml:"port"`
}

type memoryTag struct {
	XMLName       xml.Name `xml:"memory"`
	KiloByte      uint64      `xml:"kilobyte"`
	KiloByteTotal uint64      `xml:"kilobytetotal"`
	Percent       float64      `xml:"percent"`
	PercentTotal  float64      `xml:"percenttotal"`
}

type cpuTag struct {
	XMLName      xml.Name `xml:"cpu"`
	Percent      float64      `xml:"percent"`
	PercentTotal float64      `xml:"percenttotal"`
}

type portTag struct {
	XMLName       xml.Name `xml:"port"`
	ResponseTime float64   `xml:"responsetime"`
}

type serviceGroupsTag struct {
	XMLName       xml.Name          `xml:"servicegroups"`
	ServiceGroups []serviceGroupTag `xml:"servicegroup"`
}

type serviceGroupTag struct {
	XMLName xml.Name `xml:"servicegroup"`
	Name    string   `xml:"name,attr"`

	Services []string `xml:"service"`
}

func (s serviceTag) StatusString() (status string) {
	switch {
	case s.Monitor == 0:
		status = "unknown"
	case s.Monitor == 2:
		status = "starting"
	case s.Status == 0:
		status = "running"
	default:
		status = "failing"
	}
	return
}

func (t serviceGroupsTag) Get(name string) (group serviceGroupTag, found bool) {
	for _, g := range t.ServiceGroups {
		if g.Name == name {
			group = g
			found = true
			return
		}
	}
	return
}

func (t serviceGroupTag) Contains(name string) bool {
	for _, serviceName := range t.Services {
		if serviceName == name {
			return true
		}
	}
	return false
}

func (status status) ServicesInGroup(name string) (services []Service) {
	services = []Service{}

	serviceGroupTag, found := status.ServiceGroups.Get(name)
	if !found {
		return
	}

	for _, serviceTag := range status.Services.Services {
		if serviceGroupTag.Contains(serviceTag.Name) {
			service := Service{
				Monitored: serviceTag.Monitor > 0,
				Status:    serviceTag.StatusString(),
			}

			services = append(services, service)
		}
	}

	return
}

func (status status) ServiceStatsInGroup(name string) (serviceStats map[string]interface{}) {
	serviceStats = make(map[string]interface{})

	serviceGroupTag, found := status.ServiceGroups.Get(name)
	if !found {
		return
	}

	for _, serviceTag := range status.Services.Services {
		if serviceGroupTag.Contains(serviceTag.Name) {
			stat :=  map[string]interface{}{
				"CollectedSeconds": strconv.FormatUint(serviceTag.CollectedSeconds, 10),
				"Name": serviceTag.Name,
				"State": serviceTag.StatusString,
				"Monitor": strconv.Itoa(serviceTag.Monitor),
				"Pid": strconv.Itoa(serviceTag.Pid),
				"Ppid": strconv.Itoa(serviceTag.Ppid),
				"Uptime": strconv.Itoa(serviceTag.Uptime),
				"Children": strconv.Itoa(serviceTag.Children),
				"MemoryKiloByte": strconv.FormatUint(serviceTag.Memory.KiloByte, 10),
				"MemoryKiloByteTotal": strconv.FormatUint(serviceTag.Memory.KiloByteTotal, 10),
				"MemoryPercent": strconv.FormatFloat(serviceTag.Memory.Percent, 'f', -1, 64),
				"MemoryPercentTotal": strconv.FormatFloat(serviceTag.Memory.PercentTotal, 'f', -1, 64),
				"CPUPercent": strconv.FormatFloat(serviceTag.CPU.Percent, 'f', -1, 64),
				"CPUPercentTotal": strconv.FormatFloat(serviceTag.CPU.PercentTotal, 'f', -1, 64),
				"PortResponseTime": strconv.FormatFloat(serviceTag.Port.ResponseTime, 'f', -1, 64),
			}
			serviceStats[serviceTag.Name] = stat
		}
	}

	return
}

func (status status) GetIncarnation() (int, error) {
	return strconv.Atoi(status.Incarnation)
}
