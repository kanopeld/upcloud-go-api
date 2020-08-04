package request

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

// Constants
const (
	PasswordDeliveryNone  = "none"
	PasswordDeliveryEmail = "email"
	PasswordDeliverySMS   = "sms"

	ServerStopTypeSoft = "soft"
	ServerStopTypeHard = "hard"

	RestartTimeoutActionDestroy = "destroy"
	RestartTimeoutActionIgnore  = "ignore"
)

// GetServerDetailsRequest represents a request for retrieving details about a server
type GetServerDetailsRequest struct {
	UUID string
}

// RequestURL implements the Request interface
func (r *GetServerDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s", r.UUID)
}

// CreateServerIPAddressSlice is a slice of strings
// It exists to allow for a custom JSON marshaller.
type CreateServerIPAddressSlice []CreateServerIPAddress

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (s CreateServerIPAddressSlice) MarshalJSON() ([]byte, error) {
	v := struct {
		IPAddress []CreateServerIPAddress `json:"ip_address"`
	}{}
	v.IPAddress = s

	return json.Marshal(v)
}

// CreateServerStorageDeviceSlice is a slice of strings
// It exists to allow for a custom JSON marshaller.
type CreateServerStorageDeviceSlice []upcloud.CreateServerStorageDevice

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (s CreateServerStorageDeviceSlice) MarshalJSON() ([]byte, error) {
	v := struct {
		StorageDevice []upcloud.CreateServerStorageDevice `json:"storage_device"`
	}{}
	v.StorageDevice = s

	return json.Marshal(v)
}

// CreateServerRequest represents a request for creating a new server
type CreateServerRequest struct {
	XMLName xml.Name `xml:"server" json:"-"`

	AvoidHost  string `xml:"avoid_host,omitempty" json:"avoid_host,omitempty"`
	BootOrder  string `xml:"boot_order,omitempty" json:"boot_order,omitempty"`
	CoreNumber int    `xml:"core_number,omitempty" json:"core_number,omitempty"`
	// TODO: Convert to boolean
	Firewall         string                         `xml:"firewall,omitempty" json:"firewall,omitempty"`
	Hostname         string                         `xml:"hostname" json:"hostname"`
	IPAddresses      CreateServerIPAddressSlice     `xml:"ip_addresses>ip_address" json:"ip_addresses"`
	LoginUser        *LoginUser                     `xml:"login_user,omitempty" json:"login_user,omitempty"`
	MemoryAmount     int                            `xml:"memory_amount,omitempty" json:"memory_amount,omitempty"`
	PasswordDelivery string                         `xml:"password_delivery,omitempty" json:"password_delivery,omitempty"`
	Plan             string                         `xml:"plan,omitempty" json:"plan,omitempty"`
	StorageDevices   CreateServerStorageDeviceSlice `xml:"storage_devices>storage_device" json:"storage_devices"`
	TimeZone         string                         `xml:"timezone,omitempty" json:"timezone,omitempty"`
	Title            string                         `xml:"title" json:"title"`
	UserData         string                         `xml:"user_data,omitempty" json:"user_data,omitempty"`
	VideoModel       string                         `xml:"video_model,omitempty" json:"video_model,omitempty"`
	// TODO: Convert to boolean
	VNC         string `xml:"vnc,omitempty" json:"vnc,omitempty"`
	VNCPassword string `xml:"vnc_password,omitempty" json:"vnc_password,omitempty"`
	Zone        string `xml:"zone" json:"zone"`
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r CreateServerRequest) MarshalJSON() ([]byte, error) {
	type localCreateServerRequest CreateServerRequest
	v := struct {
		Server localCreateServerRequest `json:"server"`
	}{}
	v.Server = localCreateServerRequest(r)

	return json.Marshal(&v)
}

// RequestURL implements the Request interface
func (r *CreateServerRequest) RequestURL() string {
	return "/server"
}

// SSHKeySlice is a slice of strings
// It exists to allow for a custom JSON unmarshaller.
type SSHKeySlice []string

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *SSHKeySlice) UnmarshalJSON(b []byte) error {
	v := struct {
		SSHKey []string `json:"ssh_key"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.SSHKey

	return nil
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (s SSHKeySlice) MarshalJSON() ([]byte, error) {
	v := struct {
		SSHKey []string `json:"ssh_key"`
	}{}

	v.SSHKey = s

	return json.Marshal(v)
}

// LoginUser represents the login_user block when creating a new server
type LoginUser struct {
	CreatePassword string      `xml:"create_password,omitempty" json:"create_password,omitempty"`
	Username       string      `xml:"username,omitempty" json:"username,omitempty"`
	SSHKeys        SSHKeySlice `xml:"ssh_keys>ssh_key,omitempty" json:"ssh_keys"`
}

// CreateServerIPAddress represents an IP address for a CreateServerRequest
type CreateServerIPAddress struct {
	Access string `xml:"access" json:"access"`
	Family string `xml:"family" json:"family"`
}

// WaitForServerStateRequest represents a request to wait for a server to enter or exit a specific state
type WaitForServerStateRequest struct {
	UUID           string
	DesiredState   string
	UndesiredState string
	Timeout        time.Duration
}

// StartServerRequest represents a request to start a server
type StartServerRequest struct {
	UUID string

	// TODO: Start server requests have no timeout in the API
	Timeout time.Duration
}

// RequestURL implements the Request interface
func (r *StartServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/start", r.UUID)
}

// StopServerRequest represents a request to stop a server
type StopServerRequest struct {
	XMLName xml.Name `xml:"stop_server" json:"-"`

	UUID string `xml:"-" json:"-"`

	StopType string        `xml:"stop_type,omitempty" json:"stop_type,omitempty"`
	Timeout  time.Duration `xml:"timeout,omitempty" json:"timeout,omitempty,string"`
}

// RequestURL implements the Request interface
func (r *StopServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/stop", r.UUID)
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r StopServerRequest) MarshalJSON() ([]byte, error) {
	type localStopServerRequest StopServerRequest
	v := struct {
		StopServerRequest localStopServerRequest `json:"stop_server"`
	}{}
	v.StopServerRequest = localStopServerRequest(r)
	v.StopServerRequest.Timeout = v.StopServerRequest.Timeout / 1e9

	return json.Marshal(&v)
}

// MarshalXML implements a custom marshaller for StopServerRequest which converts the timeout to seconds
func (r *StopServerRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias StopServerRequest

	return e.Encode(&struct {
		Timeout int `xml:"timeout,omitempty"`
		*Alias
	}{
		Timeout: int(r.Timeout.Seconds()),
		Alias:   (*Alias)(r),
	})
}

// RestartServerRequest represents a request to restart a server
type RestartServerRequest struct {
	XMLName xml.Name `xml:"restart_server" json:"-"`

	UUID string `xml:"-" json:"-"`

	StopType      string        `xml:"stop_type,omitempty" json:"stop_type,omitempty"`
	Timeout       time.Duration `xml:"timeout,omitempty"  json:"timeout,omitempty,string"`
	TimeoutAction string        `xml:"timeout_action,omitempty" json:"timeout_action,omitempty"`
}

// RequestURL implements the Request interface
func (r *RestartServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/restart", r.UUID)
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r RestartServerRequest) MarshalJSON() ([]byte, error) {
	type localRestartServerRequest RestartServerRequest
	v := struct {
		RestartServerRequest localRestartServerRequest `json:"restart_server"`
	}{}
	v.RestartServerRequest = localRestartServerRequest(r)
	v.RestartServerRequest.Timeout = v.RestartServerRequest.Timeout / 1e9

	return json.Marshal(&v)
}

// MarshalXML implements a custom marshaller for RestartServerRequest which converts the timeout to seconds
func (r *RestartServerRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias RestartServerRequest

	return e.Encode(&struct {
		Timeout int `xml:"timeout,omitempty"`
		*Alias
	}{
		Timeout: int(r.Timeout.Seconds()),
		Alias:   (*Alias)(r),
	})
}

// ModifyServerRequest represents a request to modify a server
type ModifyServerRequest struct {
	XMLName xml.Name `xml:"server" json:"-"`

	UUID string `xml:"-" json:"-"`

	AvoidHost  string `xml:"avoid_host,omitempty" json:"avoid_host,omitempty"`
	BootOrder  string `xml:"boot_order,omitempty" json:"boot_order,omitempty"`
	CoreNumber int    `xml:"core_number,omitempty" json:"core_number,omitempty,string"`
	// TODO: Convert to boolean
	Firewall     string `xml:"firewall,omitempty" json:"firewall,omitempty"`
	Hostname     string `xml:"hostname,omitempty" json:"hostname,omitempty"`
	MemoryAmount int    `xml:"memory_amount,omitempty" json:"memory_amount,omitempty,string"`
	Plan         string `xml:"plan,omitempty" json:"plan,omitempty"`
	TimeZone     string `xml:"timezone,omitempty" json:"timezone,omitempty"`
	Title        string `xml:"title,omitempty" json:"title,omitempty"`
	VideoModel   string `xml:"video_model,omitempty" json:"video_model,omitempty"`
	// TODO: Convert to boolean
	VNC         string `xml:"vnc,omitempty" json:"vnc,omitempty"`
	VNCPassword string `xml:"vnc_password,omitempty" json:"vnc_password,omitempty"`
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r ModifyServerRequest) MarshalJSON() ([]byte, error) {
	type localModifyServerRequest ModifyServerRequest
	v := struct {
		ModifyServerRequest localModifyServerRequest `json:"server"`
	}{}
	v.ModifyServerRequest = localModifyServerRequest(r)

	return json.Marshal(&v)
}

// RequestURL implements the Request interface
func (r *ModifyServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s", r.UUID)
}

// DeleteServerRequest represents a request to delete a server
type DeleteServerRequest struct {
	UUID string
}

// RequestURL implements the Request interface
func (r *DeleteServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s", r.UUID)
}

// DeleteServerAndStoragesRequest represents a request to delete a server and all attached storages
type DeleteServerAndStoragesRequest struct {
	UUID string
}

// RequestURL implements the Request interface
func (r *DeleteServerAndStoragesRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/?storages=1", r.UUID)
}

// TagServerRequest represents a request to tag a server with one or more tags
type TagServerRequest struct {
	UUID string
	Tags []string
}

// RequestURL implements the Request interface
func (r *TagServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/tag/%s", r.UUID, strings.Join(r.Tags, ","))
}

// UntagServerRequest represents a request to remove one or more tags from a server
type UntagServerRequest struct {
	UUID string
	Tags []string
}

// RequestURL implements the Request interface
func (r *UntagServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/untag/%s", r.UUID, strings.Join(r.Tags, ","))
}
