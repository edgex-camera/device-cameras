package jdevice

type JDevice struct {
	Name    string
	Type    string
	Id      string
	Camera  Camera
	Control Control
}

type JDeviceConfig struct {
	Enabled bool   `json:enabled`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Id      string `json:"id"`
	Control string `json:"control"` // onvif/sip/none
}
