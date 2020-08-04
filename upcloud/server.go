package upcloud

import (
	"encoding/json"
)

// Constants
const (
	ServerStateStarted     = "started"
	ServerStateStopped     = "stopped"
	ServerStateMaintenance = "maintenance"
	ServerStateError       = "error"

	VideoModelVGA    = "vga"
	VideoModelCirrus = "cirrus"

	StopTypeSoft = "soft"
	StopTypeHard = "hard"
)

// ServerConfigurations represents a /server_size response
type ServerConfigurations struct {
	ServerConfigurations []ServerConfiguration `xml:"server_size"`
}

// ServerConfiguration represents a server configuration
type ServerConfiguration struct {
	CoreNumber   int `xml:"core_number"`
	MemoryAmount int `xml:"memory_amount"`
}

// Servers represents a /server response
type Servers struct {
	Servers []Server `xml:"server" json:"servers"`
}

func (s *Servers) UnmarshalJSON(b []byte) error {
	type serverWrapper struct {
		Servers []Server `json:"server"`
	}

	v := struct {
		Servers serverWrapper `json:"servers"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.Servers = v.Servers.Servers

	return nil
}

type TagSlice []string

func (t *TagSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		Tags []string `json:"tag"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*t) = v.Tags

	return nil
}

// Server represents a server
type Server struct {
	CoreNumber   int      `xml:"core_number" json:"core_number,string"`
	Hostname     string   `xml:"hostname" json:"hostname"`
	License      float64  `xml:"license" json:"license"`
	MemoryAmount int      `xml:"memory_amount" json:"memory_amount,string"`
	Plan         string   `xml:"plan" json:"plan"`
	Progress     int      `xml:"progress" json:"progress,string"`
	State        string   `xml:"state" json:"state"`
	Tags         TagSlice `xml:"tags>tag" json:"tags"`
	Title        string   `xml:"title" json:"title"`
	UUID         string   `xml:"uuid" json:"uuid"`
	Zone         string   `xml:"zone" json:"zone"`
}

type IPAddressSlice []IPAddress

func (i *IPAddressSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		IPAddresses []IPAddress `json:"ip_address"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*i) = v.IPAddresses

	return nil
}

type ServerStorageDeviceSlice []ServerStorageDevice

func (s *ServerStorageDeviceSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		StorageDevices []ServerStorageDevice `json:"storage_device"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.StorageDevices

	return nil
}

// ServerDetails represents details about a server
type ServerDetails struct {
	Server

	BootOrder  string `xml:"boot_order" json:"boot_order"`
	CoreNumber int    `xml:"core_number" json:"core_number,string"`
	// TODO: Convert to boolean
	Firewall       string                   `xml:"firewall" json:"firewall"`
	Host           int                      `xml:"host" json:"host"`
	IPAddresses    IPAddressSlice           `xml:"ip_addresses>ip_address" json:"ip_addresses"`
	NICModel       string                   `xml:"nic_model" json:"nic_model"`
	StorageDevices ServerStorageDeviceSlice `xml:"storage_devices>storage_device" json:"storage_devices"`
	Timezone       string                   `xml:"timezone" json:"timezone"`
	VideoModel     string                   `xml:"video_model" json:"video_model"`
	// TODO: Convert to boolean
	VNC         string `xml:"vnc" json:"vnc"`
	VNCHost     string `xml:"vnc_host" json:"vnc_host"`
	VNCPassword string `xml:"vnc_password" json:"vnc_password"`
	VNCPort     int    `xml:"vnc_port" json:"vnc_port"`
}

func (s *ServerDetails) UnmarshalJSON(b []byte) error {
	type localServerDetails ServerDetails

	v := struct {
		ServerDetails localServerDetails `json:"server"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = ServerDetails(v.ServerDetails)

	return nil
}
